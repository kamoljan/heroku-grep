package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/bgentry/heroku-go"
)

func main() {
	wanted := make(chan string, 10)
	var wg sync.WaitGroup
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

	for _, app := range apps {
		wg.Add(1)
		go func(appName string) {
			defer wg.Done()
			m, err := h.ConfigVarInfo(appName)
			if err != nil {
				panic(err)
			}
			for _, v := range m {
				if re.FindString(v) != "" {
					wanted <- appName
				}
			}
		}(app.Name)
	}

	go func() {
		for r := range wanted {
			fmt.Println(r)
		}
	}()

	wg.Wait()
}

//func findConfig(appName string, h *heroku.Client, re *regexp.Regexp) {
//	m, err := h.ConfigVarInfo(appName)
//	if err != nil {
//		panic(err)
//	}
//	for _, v := range m {
//		if re.FindString(v) != "" {
//			fmt.Println(appName)
//		}
//	}
//}
