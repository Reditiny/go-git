package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	CURRENT   = "."
	REPO      = ".git"
	HEAD      = "HEAD"
	HEAD_INIT = "ref: refs/heads/master\n"
	HEADS     = "heads"
	OBJECTS   = "objects"
	REFS      = "refs"
	TAGS      = "tags"
	Sha1Full  = "^[0-9a-f]{40}$"
	Sha1Short = "^[0-9a-f]{10}$"
	IGNORE    = ".gitignore"
)

func RepositoryExist(folder string) bool {
	repoPath := filepath.Join(folder, REPO)
	_, err := os.Stat(repoPath)
	if err != nil {
		return false
	}
	return true
}

// InitRepo 初始化仓库
// dir 为仓库所在目录，当前不支持指定目录，dir 提供给未来扩展以及 test 使用
func InitRepo(dir string) (err error) {
	if RepositoryExist(dir) {
		fmt.Println("reinitialize an existed repository")
		if err = os.RemoveAll(filepath.Join(dir, REPO)); err != nil {
			return
		}
	}
	// .git
	if err = os.Mkdir(filepath.Join(dir, REPO), 0755); err != nil {
		return
	}
	// HEAD
	_, err = os.Create(filepath.Join(dir, REPO, HEAD))
	if err != nil {
		return
	}
	if err = os.WriteFile(filepath.Join(dir, REPO, HEAD), []byte(HEAD_INIT), 0644); err != nil {
		return
	}
	// refs
	if err = os.Mkdir(filepath.Join(dir, REPO, REFS), 0755); err != nil {
		return
	}
	if err = os.Mkdir(filepath.Join(dir, REPO, REFS, HEADS), 0755); err != nil {
		return
	}
	if err = os.Mkdir(filepath.Join(dir, REPO, REFS, TAGS), 0755); err != nil {
		return
	}
	// objects
	if err = os.Mkdir(filepath.Join(dir, REPO, OBJECTS), 0755); err != nil {
		return
	}
	// index
	data := make([]byte, 12)
	data[0] = 'D'
	data[1] = 'I'
	data[2] = 'R'
	data[3] = 'C'
	binary.BigEndian.PutUint32(data[4:8], 0)
	binary.BigEndian.PutUint32(data[8:12], 0)
	if err = os.WriteFile(filepath.Join(dir, REPO, INDEX), data, 0644); err != nil {
		return err
	}
	//// master
	//gitObject := GitObject{Type: Commit, Size: 0, Content: []byte("")}
	//sha1, err := gitObject.Sha1()
	//if err != nil {
	//	return
	//}
	//err = WriteObject(dir, &gitObject)
	//if err != nil {
	//	return
	//}
	//open, err = os.Create(filepath.Join(dir, REPO, REFS, HEADS, "master"))
	//if err != nil {
	//	return
	//}
	//_, err = open.Write([]byte(sha1 + "\n"))
	//if err != nil {
	//	return err
	//}
	return nil
}

func FolderCheck(path string) bool {
	folderInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("folder not exist:", path)
		} else {
			fmt.Println("fail to get file info:", err)
		}
		return false
	}
	if !folderInfo.IsDir() {
		fmt.Println(path, "is not a folder")
		return false
	}
	return true
}

// GetSha1ByPath 计算文件的 SHA-1 哈希值
func GetSha1ByPath(objType, filePath string, write bool) (string, error) {
	// 读取文件得到 gitObject
	file, err := os.ReadFile(filepath.Join(CURRENT, filePath))
	if err != nil {
		return "", errors.New("Error opening file:" + err.Error())
	}
	gitObj := GitObject{Type: objType, Size: len(file), Content: file}
	sha1, err := gitObj.Sha1()
	// write 决定是否写入对象数据库
	if write {
		err = WriteObject(CURRENT, &gitObj)
		if err != nil {
			return "", errors.New("Error writing object:" + err.Error())
		}
	}
	return sha1, nil
}

// GetSha1ByName 计算 SHA-1 输入可以是 HEAD, tag, branch, short and long hashes
func GetSha1ByName(name string) string {
	// name 本身就是 SHA-1
	matchFull, err := regexp.Match(Sha1Full, []byte(name))
	if err == nil && matchFull {
		return name
	}
	// HEAD 引用的 SHA-1 值
	if name == "HEAD" {
		data, err := os.ReadFile(filepath.Join(CURRENT, REPO, HEAD))
		fmt.Println(string(data))
		if err == nil {
			return GetSha1ByName(string(data[:len(data)-1]))
		}
	}

	// name 是 SHA-1 的前 10 位
	matchShort, err := regexp.Match(Sha1Short, []byte(name))
	if err == nil && matchShort {
		sha := getShaByShort(name)
		if sha != "" {
			return sha
		}
	}
	// name 是 tag
	sha := getFromTag(name)
	if sha != "" {
		return sha
	}
	// name 是 branch
	fmt.Println(name)
	sha = getFromBranch(name)
	if sha != "" {
		return sha
	}
	return ""
}

func getShaByShort(shortSha string) string {
	prefix := shortSha[:2]
	dirPath := filepath.Join(".", REPO, OBJECTS, prefix)
	infos, err := os.ReadDir(dirPath)
	if err == nil {
		for _, info := range infos {
			if info.Name()[:len(shortSha)-2] == shortSha[2:] {
				return prefix + info.Name()
			}
		}
	}
	return ""
}

func getFromTag(name string) string {
	split := strings.Split(name, "/")
	name = split[len(split)-1]
	tagPath := filepath.Join(".", REPO, REFS, TAGS, name)
	file, err := os.ReadFile(tagPath)
	if err == nil {
		return string(file[:len(file)-1])
	}
	return ""
}

func getFromBranch(name string) string {
	split := strings.Split(name, "/")
	name = split[len(split)-1]
	branchPath := filepath.Join(".", REPO, REFS, "heads", name)
	file, err := os.ReadFile(branchPath)
	if err == nil {
		return string(file[:len(file)-1])
	}
	return ""
}

// ReadIgnore 读取 ignore 规则 当前只支持仓库根目录下的 .gitignore 文件
// todo: 还有其他地方配置的 ignore 规则
func ReadIgnore() (map[string]struct{}, error) {
	m := make(map[string]struct{})
	ignorePath := filepath.Join(CURRENT, IGNORE)
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		return nil, errors.New("Error reading ignore file:" + err.Error())
	}
	// 按行解析 ignore 规则
	lines := bytes.Split(data, []byte("\n"))
	for i, _ := range lines {
		if len(bytes.TrimSpace(lines[i])) != 0 {
			parsed, ok := ignoreParse1(lines[i])
			if ok {
				m[string(parsed)] = struct{}{}
			} else {
				delete(m, string(parsed))
			}
		}
	}
	return m, nil
}

func ignoreParse1(data []byte) ([]byte, bool) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || data[0] == '#' {
		return nil, false
	} else if data[0] == '!' {
		return data[1:], false
	} else if data[0] == '\\' {
		return data[1:], true
	} else {
		return data, true
	}
}

func CurrentBranch() (string, error) {
	headPath := filepath.Join(CURRENT, REPO, HEAD)
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(string(data), "ref: refs/heads/") {
		fmt.Println("On branch", string(data[16:]))
		refPath := filepath.Join(".", REPO, string(data[5:]))
		data, err = os.ReadFile(refPath)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("Detached HEAD", string(data))
	}
	return string(data), nil
}

func GetWorkSpaceInfo(folderPath string) (map[string]string, error) {
	nameToSha := make(map[string]string)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// 如果当前路径不是文件夹，则输出相对路径
		if !info.IsDir() && path[0] != '.' {
			nameToSha[path], err = GetSha1ByPath("blob", path, false)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return nameToSha, nil
}
