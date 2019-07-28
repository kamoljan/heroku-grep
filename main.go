package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bgentry/heroku-go"
	"github.com/korovkin/limiter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s argument1 \n ", os.Args[0])
		os.Exit(1)
	}

	re, err := regexp.Compile("(?i)" + strings.Replace(os.Args[1], " ", "[ \\._-]", -1))
	if err != nil {
		panic(err)
	}

	h := heroku.Client{Username: os.Getenv("HEROKU_USERNAME"), Password: os.Getenv("HEROKU_PASSWORD")}
	apps, err := h.AppList(&heroku.ListRange{Field: "name", Max: 1000})
	if err != nil {
		panic(err)
	}

	limit := limiter.NewConcurrencyLimiter(30)
	for _, app := range apps {
		a := app.Name // NOTE: necessary!
		limit.Execute(func() {
			m, err := h.ConfigVarInfo(a)
			if err != nil {
				panic(err)
			}
			for _, v := range m {
				if re.FindString(v) != "" {
					fmt.Println(a)
				}
			}
		})
	}
	limit.Wait()
}
