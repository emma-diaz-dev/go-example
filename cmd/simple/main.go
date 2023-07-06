package main

import (
	"github.com/emma-diaz-dev/jampp-test/internal/cfg"
	"github.com/emma-diaz-dev/jampp-test/internal/file"
	"github.com/emma-diaz-dev/jampp-test/internal/generator"
)

var (
	conf = cfg.GetConfig()
	imp  = make(map[string]string, 0)
)

func main() {
	go generator.FakeInfo()
	go file.Read(conf.ImpPath, file.SimpleStrategy(&imp))
	file.Read(conf.ClickPath, file.SimpleStrategy(&imp))
}
