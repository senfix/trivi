package gpc

import "os"

type Generator struct {
}

func (t *Generator) writeFile(path string, data []byte) {

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(data)
}
