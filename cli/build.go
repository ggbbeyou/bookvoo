package main

import (
	"os"

	"github.com/yzimhao/utilgo/pack"
)

func main() {
	pack.Build("../", "../dist/bookVoo_{{.OS}}_{{.Arch}}", os.Args[1])
}
