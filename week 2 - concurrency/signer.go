package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	quotaMd5Ch := make(chan struct{}, 1)

	for val := range in {
		wg.Add(1)

		go func(val interface{}) {
			defer wg.Done()
			data := strconv.Itoa(val.(int))

			resultMd5Ch := make(chan string)
			resultCrc32Ch := make(chan string)
			resultCrc32WithMd5Ch := make(chan string)

			// DataSignerMd5 по условию должен выполняться один раз в любой момент
			// времени. Поэтому нужно блокировать другие попытки выполнения
			// до тех пор, пока в канале quotaMd5Ch есть что-то
			go func(resultMd5Chan chan string, data string) {
				quotaMd5Ch <- struct{}{}
				resultMd5Chan <- DataSignerMd5(data)
				<-quotaMd5Ch
			}(resultMd5Ch, data)

			go func(resultCrc32Chan chan string) {
				resultCrc32Chan <- DataSignerCrc32(data)
			}(resultCrc32Ch)

			go func(resultCrc32WithMd5Chan chan string) {
				resultCrc32WithMd5Chan <- DataSignerCrc32(<-resultMd5Ch)
			}(resultCrc32WithMd5Ch)

			out <- <-resultCrc32Ch + "~" + <-resultCrc32WithMd5Ch

		}(val)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for val := range in {
		wg.Add(1)

		go func(val interface{}) {
			// Сколько раз нужно вычислить MultiHash
			const MAX_TH = 6

			defer wg.Done()
			data := val.(string)
			results := make([]string, MAX_TH)
			wgInner := sync.WaitGroup{}
			wgInner.Add(MAX_TH)

			for i := 0; i < MAX_TH; i++ {
				go func(idx int, data string) {
					results[idx] = DataSignerCrc32(strconv.Itoa(idx) + data)
					wgInner.Done()
				}(i, data)
			}

			wgInner.Wait()
			out <- strings.Join(results, "")
		}(val)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	result := []string{}

	for val := range in {
		result = append(result, val.(string))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	out <- strings.Join(result, "_")
}

func ExecutePipeline(jobs ...job) {
	rCh := make(chan interface{})
	done := make(chan struct{})

	if len(jobs) > 1 {
		go func(wCh chan interface{}) {
			jobs[0](nil, wCh)
			close(wCh)
		}(rCh)

		// Вначале было for i := range jobs[1 : len(jobs)-1]
		// Думал что i начинется с 1, по факту с 0
		// от этого два раза запустилась первая job
		for i := 1; i < len(jobs)-1; i++ {
			wCh := make(chan interface{})

			go func(idx int, rCh chan interface{}, wCh chan interface{}) {
				// Запускал сразу job, которую получал из головы цикла
				// (for i, job := range jobs...)
				// при выполнении горутины выполнялась последняя job
				jobs[idx](rCh, wCh)
				close(wCh)
			}(i, rCh, wCh)

			rCh = wCh
		}
		go func(rCh chan interface{}) {
			jobs[len(jobs)-1](rCh, nil)
			done <- struct{}{}
		}(rCh)

	} else {
		//  Если имеем только одну job, которая только пишет в канал
		go func(wCh chan interface{}) {
			jobs[0](nil, wCh)
			close(wCh)
		}(rCh)

		// Нужна дополнительная job, которая читает из этого канала.
		// Это не даст впасть в deadlock
		go func(rCh chan interface{}) {
			for v := range rCh {
				fmt.Println(v)
			}
			done <- struct{}{}
		}(rCh)
	}
	<-done
}
