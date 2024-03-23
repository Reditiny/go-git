/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// revParseCmd represents the revParse command
// 解析解析给定的标识符对应的 SHA1 值
// 当前会按优先级返回第一个匹配的值 todo: 有多个匹配时的处理
var revParseCmd = &cobra.Command{
	Use:   "rev-parse NAME",
	Short: "Parse revision (or other objects) identifiers",
	Long:  `Parse revision (or other objects) identifiers`,
	Run:   RevParseFunc,
}

var typeFlag string

func init() {

	revParseCmd.Flags().StringVarP(&typeFlag, "type", "t", "commit", "show the type of the object")

	rootCmd.AddCommand(revParseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// revParseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// revParseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func RevParseFunc(cmd *cobra.Command, args []string) {
	// 参数检验
	if len(args) != 1 {
		fmt.Println("Usage: go-git rev-parse NAME")
		return
	}
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	// 读取对象
	sha1 := common.GetSha1ByName(args[0])
	if sha1 == "" {
		fmt.Println("Not a valid object name")
		return
	}
	fmt.Println(sha1)
}
