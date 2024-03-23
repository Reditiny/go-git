/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// lsTreeCmd represents the lsTree command
// 显示树对象的内容，当前不支持 -r 选项
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree sha1",
	Short: "Display the contents of a tree object",
	Long:  `Display the contents of a tree object.`,
	Run:   LsTreeFunc,
}

var recursive bool

func init() {

	lsTreeCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "recursive tree display")

	rootCmd.AddCommand(lsTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func LsTreeFunc(cmd *cobra.Command, args []string) {
	// 参数检验
	if len(args) != 1 {
		fmt.Println("Usage: go-git ls-tree SHA")
		return
	}
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	// 读取树对象
	tree, err := common.ParseTree(common.CURRENT, args[0])
	if err != nil {
		fmt.Println("Error reading tree:", err)
		return
	}
	for i := range tree.Children {
		tree.Children[i].Serialize(recursive)
	}
}
