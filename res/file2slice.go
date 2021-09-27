package res

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// FileToSlice convert file to slice and restore it in a map[int][]byte, if fail, an error returned
func FileToSlice(inPath, outPath string) error {
	dir, err := ioutil.ReadDir(inPath)
	if err != nil {
		return err
	}

	fOut, err := os.Create(outPath + "/resource.go")
	if err != nil {
		return err
	}
	defer fOut.Close()

	// 写入包名, 并且初始化一下map
	// map[int][]byte {
	// 	1: []byte("xxx"),
	// }
	if _, err := fmt.Fprintf(fOut, "package ChineseChess\n\nvar resMap = map[int][]byte {\n"); err != nil {
		return err
	}

	// 遍历目录下所有文件
	for _, fIn := range dir {
		// 生成变量名
		varName := ""
		if ok := strings.HasSuffix(fIn.Name(), ".png"); ok {
			varName = "Img" + strings.TrimSuffix(fIn.Name(), ".png")
		} else if ok := strings.HasSuffix(fIn.Name(), ".wav"); ok {
			varName = "Music" + strings.TrimSuffix(fIn.Name(), ".wav")
		} else {
			continue
		}

		// 写入map
		// 打开输入文件
		f, err := os.Open(inPath + "/" + fIn.Name())
		if err != nil {
			return err
		}
		defer f.Close()
		// 转化文件
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		// 写入文件
		if _, err := fmt.Fprintf(fOut, " %s : []byte(%q),\n", varName, bs); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(fOut, "}"); err != nil {
		return err
	}
	return nil

}
