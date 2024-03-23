/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show the working tree status.",
	Long:  `Show the working tree status.`,
	Run:   StatusFunc,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StatusFunc(cmd *cobra.Command, args []string) {
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	// 当前分支情况
	curCommitNameToSha1, err := common.GetCommitInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 当前 index 情况
	indexNameToSha1, err := common.GetIndexInfo()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}
	// 比较当前 index 和当前 commit 的差异
	fmt.Println("difference between index and commit")
	for name, indexSha1 := range indexNameToSha1 {
		if commitSha1, ok := curCommitNameToSha1[name]; ok {
			if commitSha1 != indexSha1 {
				fmt.Println(name, "modified")
			}
			delete(curCommitNameToSha1, name)
		} else {
			fmt.Println(name, "added")
		}
	}
	for name, _ := range curCommitNameToSha1 {
		fmt.Println(name, "deleted")
	}
	// 当前工作区情况
	workSpaceNameToSha1, err := common.GetWorkSpaceInfo(common.CURRENT)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 比较当前工作区和当前 index 的差异
	fmt.Println("difference between work space and index")
	for name, workSha1 := range workSpaceNameToSha1 {
		if indexSha1, ok := indexNameToSha1[name]; ok {
			if indexSha1 != workSha1 {
				fmt.Println(name, "modified")
			}
			delete(indexNameToSha1, name)
		} else {
			fmt.Println(name, "untracked")
		}
	}
	for name, _ := range indexNameToSha1 {
		fmt.Println(name, "deleted")
	}
}
