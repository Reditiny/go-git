package common

import (
	"fmt"
	"os"
	"path"
)

func ListAllTag() {
	tagDir := path.Join(CURRENT, REPO, REFS, TAGS)
	dirEntries, err := os.ReadDir(tagDir)
	if err != nil {
		fmt.Printf("error reading directory: %v\n", err)
		return
	}
	for _, entry := range dirEntries {
		fmt.Println(entry.Name())
	}
}

func NewTag(tagName, sha1 string) {
	if sha1 == "" {
		head, err := os.ReadFile(path.Join(CURRENT, REPO, HEAD))
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return
		}
		sha1 = string(head)
	}

	tagPath := path.Join(".", REPO, REFS, TAGS, tagName)
	err := os.WriteFile(tagPath, []byte(sha1), 0644)
	if err != nil {
		fmt.Printf("error writing file: %v\n", err)
		return
	}
}
