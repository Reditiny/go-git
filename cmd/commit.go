/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
	"os"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "record changes to the repository.",
	Long:  `record changes to the repository.`,
	Run:   CommitFunc,
}

var message string

func init() {

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")

	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CommitFunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("commit command does not accept any arguments")
	}
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	if message == "" {
		message = "default commit message"
	}
	// 根据 index 生成 tree 对象
	index, err := common.ReadIndex()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}
	content, err := index.MakeTreeContent()
	if err != nil {
		fmt.Println("Error making tree:", err)
		return
	}
	treeObject := common.GitObject{
		Type:    "tree",
		Size:    len(content),
		Content: content,
	}
	err = common.WriteObject(common.CURRENT, &treeObject)
	if err != nil {
		fmt.Println("Error writing tree:", err)
		return
	}
	treeSha1, err := treeObject.Sha1()
	if err != nil {
		fmt.Println("Error making tree:", err)
		return
	}
	// 获得当前 commit 的 commitID
	commitID := common.GetCommitID("HEAD")
	// 写入 commit 对象 todo: 未实现 author 和 committer
	content1 := fmt.Sprintf("tree %s\n", treeSha1)
	content2 := fmt.Sprintf("parent %s\n", commitID)
	content3 := fmt.Sprintf("author %s\n", "author")
	content4 := fmt.Sprintf("committer %s\n", "committer")
	content5 := fmt.Sprintf("\n%s\n", message)
	commitContent := content1 + content2 + content3 + content4 + content5
	object := common.GitObject{
		Type:    "commit",
		Size:    len(commitContent),
		Content: []byte(commitContent),
	}
	err = common.WriteObject(common.CURRENT, &object)
	if err != nil {
		fmt.Println("Error writing commit:", err)
		return
	}
	// 更新 HEAD
	commitSha1, _ := object.Sha1()
	os.WriteFile(common.CURRENT+"/.git/HEAD", []byte(commitSha1+"\n"), 0644)
}
