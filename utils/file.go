package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteToFile(fileName string, content string, dir ...string) {
	var path string
	if len(dir) > 0 {
		dir = append(dir, fileName)
		path = Concat(dir...)
	} else {
		path = fileName
	}
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)

	}
}

func WriteToFileBytes(fileName string, content []byte, dir ...string) {
	var path string
	if len(dir) > 0 {
		dir = append(dir, fileName)
		path = Concat(dir...)
	} else {
		path = fileName
	}

	// 创建目录（如果需要）
	dirPath := filepath.Dir(path)
	_ = os.MkdirAll(dirPath, 0755)

	var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println(err)

	}
}
