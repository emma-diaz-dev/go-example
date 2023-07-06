package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

const (
	clickOutput = "/tmp/clicks.pipe"
	ImpName     = "/tmp/impressions.pipe"
)

var (
	imp   = make(map[string]string, 0)
	mutex sync.Mutex
)

func main() {
	go generator.FakeInfo()
	go readFile(ImpName)
	readFile(clickOutput)
}

func readFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[func: readFile] | [error: ", err.Error(), "]")
		t := time.NewTimer(time.Second)
		<-t.C
		readFile(path)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		switch path {
		case ImpName:
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
				fmt.Println("[impression_id: ", m, "] [campaign: ", v, "]")
				delete(imp, m)
				mutex.Unlock()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
