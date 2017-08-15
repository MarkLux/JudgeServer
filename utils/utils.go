package utils

import (
	"io/ioutil"
	"strings"
)

func GetFileNames(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), suffix) {
			files = append(files, strings.TrimSuffix(fi.Name(), suffix))
		}
	}
	return files, nil
}
