package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
	"github.com/emma-diaz-dev/jampp-test/internal/file"
	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

type Imp struct {
	ID       string
	Campaign string
}

var (
	conf  = cfg.GetConfig()
	chImp = make(chan Imp)
)

func main() {
	go generator.FakeInfo()
	go file.Read(conf.ImpPath, impStrategy(chImp))
	file.Read(conf.ClickPath, clickStrategy(chImp))
}

func impStrategy(ch chan Imp) func(string, *bufio.Scanner) {
	return func(path string, scanner *bufio.Scanner) {
		for scanner.Scan() {
			m := strings.Split(scanner.Text(), " ")
			if len(m) > 1 {
				go func(ch chan Imp) { ch <- Imp{m[0], m[1]} }(ch)
			}

		}
	}
}
func clickStrategy(ch chan Imp) func(string, *bufio.Scanner) {
	return func(path string, scanner *bufio.Scanner) {
		for scanner.Scan() {
			m := scanner.Text()
			go func(m string) {
				for {
					a := <-ch
					if a.ID == m {
						fmt.Println("[cmd: chan][impression_id: ", a.ID, "] [campaign: ", a.Campaign, "]")
						return
					}
					go func(a Imp) { ch <- a }(a)
				}
			}(m)
		}
	}
}
