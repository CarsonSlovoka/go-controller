package app

import (
	"flag"
	"github.com/CarsonSlovoka/go-controller/v2/app/script"
)

var (
	Script *script.Script
)

func InitApp() (err error) {
	inputFilepath := flag.String("file", "", "filepath (json format)")
	flag.Parse()
	if *inputFilepath == "" {
		panic("The input file path is empty.")
	}
	var t *script.Template
	t, err = script.NewTemplate(*inputFilepath)
	if err != nil {
		panic("script.NetTemplate error")
	}
	Script = script.NewScript(t)
	return err
}
