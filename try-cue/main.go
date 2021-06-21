package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/openapi"
)

func genOpenAPI(fileName string, filePath string) ([]byte, error) {
	defaultConfig := &openapi.Config{}

	filename := filepath.FromSlash(fileName)

	inst := cue.Build(load.Instances([]string{filename}, &load.Config{
		Dir: filePath,
	}))[0]

	if inst.Err != nil {
		return nil, inst.Err
	}

	b, err := openapi.Gen(inst, defaultConfig)
	if err != nil {
		return nil, err
	}

	var out = &bytes.Buffer{}
	_ = json.Indent(out, b, "", "   ")
	return out.Bytes(), nil
}

func main() {
	a, err := genOpenAPI("webservice.cue", "try-cue")
	if err != nil {
		fmt.Sprintf(err.Error())
	}
	fmt.Printf(string(a))
}
