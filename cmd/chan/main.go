package main

import (
	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
	"github.com/emma-diaz-dev/jampp-test/internal/file"
	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

var (
	conf  = cfg.GetConfig()
	chImp = make(chan file.Imp)
)

func main() {
	go generator.FakeInfo()
	go file.Read(conf.ImpPath, file.ImpStrategy(chImp))
	file.Read(conf.ClickPath, file.ClickStrategy(chImp))
}
