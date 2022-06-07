package app

import (
	"flag"
	"fmt"
	"github.com/CarsonSlovoka/go-controller/v2/app/script"
)

const Version = "0.0.0"

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
		panic(fmt.Errorf("script.NetTemplate error %w", err))
	}
	Script = script.NewScript(t)
	return err
}
