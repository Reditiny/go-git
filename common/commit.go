package common

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParseCommit Key-Value List with Message
// 从 commit 字节流中解析键值对列表
func ParseCommit(folder, sha string, short bool) {
	obj, err := ReadObject(folder, sha)
	if err != nil {
		fmt.Println("Error reading object:", err)
		return
	}
	if obj.GetType() != "commit" {
		return
	}
	parents := make([]string, 0)
	// 从字节流中解析键值对列表
	data := obj.Content
	split := strings.Split(string(data), "\n")
	fmt.Println("commit ", sha)
	for i := 0; i < len(split); i++ {
		if !short {
			fmt.Println(split[i])
		}
		if len(split[i]) > 6 && split[i][:6] == "parent" {
			parents = append(parents, split[i][7:])
		}
	}
	if len(parents) == 0 {
		return
	}
	fmt.Println("|")
	fmt.Println("|")
	fmt.Println("V")
	for i := range parents {
		ParseCommit(folder, parents[i], short)
	}
}

// GetCommitParent 获取 commit 的父提交
func GetCommitParent(folder, sha string) ([]string, error) {
	obj, err := ReadObject(folder, sha)
	if err != nil {
		fmt.Println("Error reading object:", err)
		return nil, err
	}
	if obj.GetType() != "commit" {
		return nil, err
	}
	parents := make([]string, 0)
	data := obj.Content
	split := strings.Split(string(data), "\n")
	for i := 0; i < len(split); i++ {
		if len(split[i]) > 6 && split[i][:6] == "parent" {
			parents = append(parents, split[i][7:])
		}
	}
	return parents, nil
}

func GetCommitID(name string) string {
	matched, err := regexp.Match(Sha1Full, []byte(name))
	if err == nil && matched {
		return name
	}
	return GetSha1ByName(name)
}

// CheckoutCommit 将 commit 的内容写入当前目录
func CheckoutCommit(commitSha1 string) error {
	gitObject, err := ReadObject(CURRENT, commitSha1)
	if err != nil {
		return errors.New("Error reading object:" + err.Error())
	}
	if gitObject.GetType() != "commit" {
		return errors.New("not a commit object")
	}

	treeSha := gitObject.Content[5:45]
	tree, err := ParseTree(CURRENT, string(treeSha))
	if err != nil {
		return errors.New("Error reading tree:" + err.Error())
	}

	for i := range tree.Children {
		err := tree.Children[i].MakeFile()
		if err != nil {
			return errors.New("Error making file:" + err.Error())
		}
	}

	return nil
}

func GetCommitInfo() (map[string]string, error) {
	commitSha := GetSha1ByName(HEAD)
	if commitSha == "" {
		return make(map[string]string), nil
	}
	commitObject, err := ReadObject(CURRENT, commitSha)
	if err != nil {
		return nil, errors.New("Error reading object:" + err.Error())
	}
	if commitObject.GetType() != Commit {
		return nil, errors.New("object is not a commit")
	}
	// 得到当前 commit 的文件名和 workSha1
	treeSha := string(commitObject.Content[5:45])
	tree, err := ParseTree(CURRENT, treeSha)
	if err != nil {
		return nil, errors.New("Error reading tree:" + err.Error())
	}

	return tree.MakeNameToSha1(), nil
}
