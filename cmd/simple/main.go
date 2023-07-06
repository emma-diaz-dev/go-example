package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
	"github.com/emma-diaz-dev/jampp-test/internal/file"
	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

var (
	conf  = cfg.GetConfig()
	imp   = make(map[string]string, 0)
	mutex sync.Mutex
)

func main() {
	go generator.FakeInfo()
	go file.Read(conf.ImpPath, simpleStrategy)
	file.Read(conf.ClickPath, simpleStrategy)
}

func simpleStrategy(path string, scanner *bufio.Scanner) {
	for scanner.Scan() {
		switch path {
		case conf.ImpPath:
			m := strings.Split(scanner.Text(), " ")
			if len(m) > 1 {
				mutex.Lock()
				imp[m[0]] = m[1]
				mutex.Unlock()
			}
		default:
			m := scanner.Text()
			mutex.Lock()
			if v, ok := imp[m]; ok {
				fmt.Println("[cmd: simple] [impression_id: ", m, "] [campaign: ", v, "]")
				delete(imp, m)
				mutex.Unlock()
			}
		}
	}
}
