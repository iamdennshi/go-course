package hw3

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//easyjson:json
type User struct {
	Browsers []string `json:",intern"`
	Email    string
	Name     string
}

var TOTAL_USERS = 1000

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	users := make([]User, TOTAL_USERS)
	in := bufio.NewScanner(file)

	for i := 0; in.Scan(); i++ {
		user := User{}
		user.UnmarshalJSON(in.Bytes())
		if i >= TOTAL_USERS {
			users = append(users, user)
		} else {
			users[i] = user
		}
	}

	seenBrowsers := []string{}
	foundUsers := ""

	for i, user := range users {
		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true

				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
