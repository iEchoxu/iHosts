package utils

import (
	"fmt"
	"ihosts/global"
	"io"
	"os"
	"strings"
)

func FileDeduplication(line string) (isDuplication bool) {
	isDuplication = false
	for _, siteName := range global.EnvConfig.StartUrls {
		if strings.Contains(line, siteName) || strings.Contains(line, "更新") {
			isDuplication = true
		}
	}

	return isDuplication
}

func CopyFile(srcFile, destFile string) (written int64, err error) {
	srcFileData, err := os.OpenFile(srcFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	destFileData, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	defer func(srcFileData *os.File) {
		err := srcFileData.Close()
		if err != nil {

		}
	}(srcFileData)

	defer func(destFileData *os.File) {
		err := destFileData.Close()
		if err != nil {

		}
	}(destFileData)

	return io.Copy(destFileData, srcFileData)
}
