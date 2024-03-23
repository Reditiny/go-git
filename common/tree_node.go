package common

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

type Node interface {
	Serialize(recursion bool)
	MakeFile() error
	MakeNameToSha1() map[string]string
}

type TreeNode struct {
	Mode     string
	Type     string
	Path     string
	Sha1     string
	Children []Node
}

func (t *TreeNode) MakeTreeObject() {

}

func (t *TreeNode) Serialize(recursion bool) {
	if recursion {
		for _, child := range t.Children {
			child.Serialize(recursion)
		}
	} else {
		fmt.Printf("%s %s %s %s\n", t.Mode, t.Type, t.Sha1, t.Path)
	}
}

func (t *TreeNode) MakeFile() error {
	// 创建目录
	_, err := os.Stat(t.Path[2:])
	if err != nil {
		err := os.Mkdir(t.Path[2:], 0755)
		if err != nil {
			return errors.New("Error making directory:" + err.Error())
		}
	}
	// 递归创建子节点
	for _, child := range t.Children {
		err = child.MakeFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeNode) MakeNameToSha1() map[string]string {
	pathToSha1 := make(map[string]string)
	for _, child := range t.Children {
		for k, v := range child.MakeNameToSha1() {
			pathToSha1[k] = v
		}
	}
	return pathToSha1
}

type BlobNode struct {
	Mode string
	Type string
	Path string
	Sha1 string
}

func (b *BlobNode) Serialize(recursion bool) {
	fmt.Printf("%s %s %s %s\n", b.Mode, b.Type, b.Sha1, b.Path)
}

func (b *BlobNode) MakeFile() error {
	// 读取对象
	object, err := ReadObject(CURRENT, b.Sha1)
	if err != nil {
		return errors.New("Error reading object:" + err.Error())
	}
	if object.GetType() != Blob {
		return errors.New("object is not a data")
	}
	// 写入文件
	err = os.WriteFile(b.Path[2:], object.Content, 0644)
	if err != nil {
		return errors.New("Error write file:" + err.Error())
	}
	return nil
}

func (t *BlobNode) MakeNameToSha1() map[string]string {
	pathToSha1 := make(map[string]string)
	pathToSha1[t.Path[2:]] = t.Sha1
	return pathToSha1
}

// ParseTree 解析树对象
func ParseTree(prefix, treeSha string) (*TreeNode, error) {
	// 读取对象
	object, err := ReadObject(".", treeSha)
	if err != nil {
		return nil, errors.New("Error reading object:" + err.Error())
	}
	if object.GetType() != Tree {
		return nil, errors.New("object is not a data")
	}
	// 解析对象内容
	treeNode := &TreeNode{}
	nodes := make([]Node, 0)
	data := object.Content
	length := len(data)
	start := 0
	for start < length {
		node, n := parseOne(prefix, data[start:])
		nodes = append(nodes, node)
		start += n
	}
	treeNode.Children = nodes
	return treeNode, nil
}

// parseOne 解析树对象的一个节点，返回节点和节点长度
// 节点构成 mode + ' ' + path + \x00 + sha1
func parseOne(prefix string, node []byte) (Node, int) {
	// 解析 mode 长度可能为 5 或 6，不足 6 时高位补 0
	spaceIndex := bytes.IndexByte(node, 0x20)
	if spaceIndex == -1 {
		return nil, -1
	}
	mode := string(node[:spaceIndex])
	if len(mode) == 5 {
		mode = "0" + mode
	}
	// 解析 path
	pathIndex := bytes.IndexByte(node, 0x00)
	path := string(node[spaceIndex+1 : pathIndex])
	// 解析 sha1 长度为 20，解析完后转化为十六进制字符串
	sha1 := hex.EncodeToString(node[pathIndex+1 : pathIndex+21])
	// 根据 mode 判断类型
	var nodeType string
	switch mode[0:2] {
	case "04":
		nodeType = "tree"
	case "16":
		panic("commit object unsupported")
	default:
		nodeType = "blob"
	}
	if nodeType == "blob" {
		return &BlobNode{mode, nodeType, prefix + "/" + path, sha1}, pathIndex + 21
	} else {
		treeNode, _ := ParseTree(prefix+"/"+path, sha1)
		treeNode.Mode = mode
		treeNode.Type = nodeType
		treeNode.Path = prefix + "/" + path
		treeNode.Sha1 = sha1
		return treeNode, pathIndex + 21
	}
}
