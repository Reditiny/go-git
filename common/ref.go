package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadRef 读取间接引用的 Sha1 值
func ReadRef(folder, ref string) (string, error) {
	refPath := filepath.Join(folder, REPO, ref)
	data, err := os.ReadFile(refPath)
	if err != nil {
		fmt.Println("Error reading ref:", err)
		return "", err
	}
	if strings.HasPrefix(string(data), "ref: ") {
		return ReadRef(folder, string(data[5:]))
	}
	return string(data), nil
}

// WalkRef 遍历所有引用
func WalkRef() {
	walkDir(REFS)
}

// relativePath 是相对于 .git 的路径, 输出时会用到
func walkDir(relativePath string) {
	dirPath := filepath.Join(CURRENT, REPO, relativePath)
	infos, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, info := range infos {
		if info.IsDir() {
			walkDir(filepath.Join(relativePath, info.Name()))
		} else {
			data, err := os.ReadFile(filepath.Join(dirPath, info.Name()))
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			fmt.Printf("%s %s/%s\n", data[0:len(data)-1], relativePath, info.Name())
		}
	}
}
