package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
)

//const CueTemplate = `cmd: [...string]
//
//rules?: [...{
//	path:          string
//	rewriteTarget: *"" | string
//}]`

func main() {
	path, _ := filepath.Abs("./test.cue")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	r := cue.Runtime{}
	template, _ := r.Compile("", string(data))
	stc, _ := template.Value().Struct()
	for i := 0; i < stc.Len(); i++ {
		val := stc.Field(i).Value
		switch val.IncompleteKind() {
		case cue.ListKind:
			iter, _ := val.List()
			for iter.Next() {
				fmt.Println(iter.Value())
			}

		}
	}
}
