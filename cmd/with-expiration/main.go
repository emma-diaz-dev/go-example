package main

import (
	"time"

	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
	"github.com/emma-diaz-dev/jampp-test/internal/file"
	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

const (
	expiration = 60
)

var (
	conf  = cfg.GetConfig()
	chImp = make(chan file.Imp)
)

func main() {
	go func() {
		for {
			expirator()
		}
	}()
	go generator.FakeInfo()
	go file.Read(conf.ImpPath, file.ImpStrategy(chImp))
	file.Read(conf.ClickPath, file.ClickStrategy(chImp))
}

func expirator() {
	imp := <-chImp
	if time.Now().Unix()-imp.Expiration >= expiration {
		// fmt.Println("[cmd: chan][impression_id: ", imp.ID, "] [campaign: ", imp.Campaign, "] [action: removed]")
		return
	}
	go func(imp file.Imp) {
		chImp <- imp
	}(imp)
}
