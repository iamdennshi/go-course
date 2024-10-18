

slow = list(map(int, filter(lambda i : i.isnumeric(), input().split())))
fast = list(map(int, filter(lambda i : i.isnumeric(), input().split())))

#+49.3 %  -34.6 %ns/op   -60.8 %B/op  -61.9 %allocs/op

postfix = ["%", "%ns/op", "%B/op", "%allocs/op"]


result = [ round((fast[i]*100/slow[i])-100,1) for i in range(len(slow)) ]

for i in range(len(result)):
  string = f'+{result[i]}' if result[i] > 0 else result[i]
  print(f'{string} {postfix[i]}', end=' ')

