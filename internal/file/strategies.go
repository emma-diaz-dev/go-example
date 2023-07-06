package file

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

type Imp struct {
	ID         string
	Campaign   string
	Expiration int64
}

func ImpStrategy(ch chan Imp) func(string, *bufio.Scanner) {
	return func(path string, scanner *bufio.Scanner) {
		for scanner.Scan() {
			m := strings.Split(scanner.Text(), " ")
			if len(m) > 1 {
				go func(ch chan Imp) { ch <- Imp{ID: m[0], Campaign: m[1], Expiration: time.Now().Unix()} }(ch)
			}

		}
	}
}

func ClickStrategy(ch chan Imp) func(string, *bufio.Scanner) {
	return func(path string, scanner *bufio.Scanner) {
		for scanner.Scan() {
			m := scanner.Text()
			go func(m string) {
				for {
					a := <-ch
					if a.ID == m {
						fmt.Println("[cmd: chan][impression_id: ", a.ID, "] [campaign: ", a.Campaign, "] [action: printed]")
						return
					}
					go func(a Imp) { ch <- a }(a)
				}
			}(m)
		}
	}
}

func SimpleStrategy(imp *map[string]string) func(path string, scanner *bufio.Scanner) {
	return func(path string, scanner *bufio.Scanner) {
		for scanner.Scan() {
			switch path {
			case conf.ImpPath:
				m := strings.Split(scanner.Text(), " ")
				if len(m) > 1 {
					mutex.Lock()
					(*imp)[m[0]] = m[1]
					mutex.Unlock()
				}
			default:
				m := scanner.Text()
				mutex.Lock()
				if v, ok := (*imp)[m]; ok {
					fmt.Println("[cmd: simple] [impression_id: ", m, "] [campaign: ", v, "] [action: printed]")
					delete((*imp), m)
					mutex.Unlock()
				}
			}
		}
	}
}
