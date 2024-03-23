/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// lsFilesCmd represents the lsFiles command
// 列出所有的 stage 文件
// todo: 通过 flag 显示跟多信息
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "List all the stage files",
	Long:  `List all the stage files`,
	Run:   LsFilesFunc,
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsFilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsFilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func LsFilesFunc(cmd *cobra.Command, args []string) {
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	gitIndex, err := common.ReadIndex()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}

	for i := range gitIndex.Entries {
		fmt.Println(gitIndex.Entries[i].Name)
	}
}
