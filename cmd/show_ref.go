/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
)

// showRefCmd represents the showRef command
// 展示所有引用 sha1 path
// 例如: 5b048242a25b614ff6c7b9a14cbcc6ef8dcaec9e refs/heads/master
var showRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "show all references in the repository",
	Long:  `show all references in the repository`,
	Run:   ShowRefFunc,
}

func init() {
	rootCmd.AddCommand(showRefCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showRefCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showRefCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ShowRefFunc(cmd *cobra.Command, args []string) {
	// 参数检验
	if len(args) > 0 {
		fmt.Println("show-ref cmd receive no args")
		return
	}
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}

	common.WalkRef()
}
