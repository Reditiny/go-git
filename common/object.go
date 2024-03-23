package common

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	Blob   = "blob"
	Commit = "commit"
	Tree   = "tree"
	Tag    = "tag"
)

type GitObject struct {
	Type    string
	Size    int
	Content []byte
}

func (g *GitObject) GetType() string {
	return g.Type
}

func (g *GitObject) Sha1() (string, error) {
	objectBytes := encodeObject(g)
	sha1Hash := sha1.New()
	_, err := sha1Hash.Write(objectBytes)
	if err != nil {
		return "", err
	}
	hashInBytes := sha1Hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

func decodeObject(data []byte) (*GitObject, error) {
	// 查找对象类型和内容的分隔符
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("Invalid Git object format")
	}

	// 提取对象类型和内容
	objectHeader := string(data[:nullIndex])
	split := strings.Split(objectHeader, " ")

	objectContent := data[nullIndex+1:]

	return &GitObject{
		Type:    split[0],
		Size:    len(objectContent),
		Content: objectContent,
	}, nil
}

func encodeObject(obj *GitObject) []byte {
	// 拼接对象类型和内容
	header := fmt.Sprintf("%s %d\x00", obj.Type, obj.Size)
	return append([]byte(header), obj.Content...)
}

func ValidType(objType string) bool {
	return objType == Blob || objType == Commit || objType == Tree || objType == Tag
}

// ReadObject 从文件中读取 Git 对象，文件使用 zlib 压缩存储
// git object 的存储格式为 [object type] [object size]\x00[object content]
// 例如 "blob 5\x00hello" (压缩前)
func ReadObject(folder, sha string) (*GitObject, error) {
	// sha1 hash 码的前两位作为文件夹名，第三位开始作为文件名
	filePath := filepath.Join(folder, REPO, OBJECTS, sha[:2], sha[2:])
	objectData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("Error reading object file: " + err.Error())
	}
	// 解压缩对象数据
	reader, err := zlib.NewReader(bytes.NewReader(objectData))
	if err != nil {
		return nil, errors.New("Error creating zlib reader: " + err.Error())
	}
	// 解析对象数据
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("Error reading uncompressed object data: " + err.Error())
	}
	gitObject, err := decodeObject(content)
	if err != nil {
		return nil, errors.New("Error parsing Git object: " + err.Error())
	}
	return gitObject, nil
}

// WriteObject 将 Git 对象写入文件，文件使用 zlib 压缩存储
func WriteObject(folder string, gitObject *GitObject) error {
	sha, err := gitObject.Sha1()
	if err != nil {
		return errors.New("Error writing data to SHA-1 hash: " + err.Error())
	}
	// 创建目录及文件
	err = os.MkdirAll(filepath.Join(folder, REPO, OBJECTS, sha[:2]), 0755)
	if err != nil {
		return errors.New("Error creating directory: " + err.Error())
	}
	filePath := filepath.Join(folder, REPO, OBJECTS, sha[:2], sha[2:])
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return errors.New("Error creating file: " + err.Error())
	}
	// 压缩后写入数据
	writer := zlib.NewWriter(file)
	defer writer.Close()
	_, err = writer.Write(encodeObject(gitObject))
	if err != nil {
		return errors.New("Error writing compressed object data: " + err.Error())
	}
	return nil
}

func FindObject(name, t string) (*GitObject, error) {
	sha := GetSha1ByName(name)
	if sha == "" {
		return nil, fmt.Errorf("not a valid object name: %s", name)
	}
	return ReadObject(".", sha)
}
