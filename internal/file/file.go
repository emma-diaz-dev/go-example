package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
)

var (
	conf = cfg.GetConfig()
)

func Read(path string, strategy func(string, *bufio.Scanner)) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[func: readFile] | [error: ", err.Error(), "]")
		t := time.NewTimer(time.Second)
		<-t.C
		Read(path, strategy)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	strategy(path, scanner)
	// for scanner.Scan() {
	// 	switch path {
	// 	case conf.ImpPath:
	// 		m := strings.Split(scanner.Text(), " ")
	// 		if len(m) > 1 {
	// 			mutex.Lock()
	// 			imp[m[0]] = m[1]
	// 			mutex.Unlock()
	// 		}
	// 	default:
	// 		m := scanner.Text()
	// 		mutex.Lock()
	// 		if v, ok := imp[m]; ok {
	// 			fmt.Println("[impression_id: ", m, "] [campaign: ", v, "]")
	// 			delete(imp, m)
	// 			mutex.Unlock()
	// 		}
	// 	}
	// }

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
