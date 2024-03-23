/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add NAME...",
	Short: "Add files contents to the index",
	Long:  `Add files contents to the index.`,
	Run:   AddFunc,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func AddFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go-git add NAME...")
		return
	}
	if !common.RepositoryExist(".") {
		fmt.Println("not a repository")
		return
	}
	// index 信息
	gitIndex, err := common.ReadIndex()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}

	for i := range gitIndex.Entries {
		fmt.Println(gitIndex.Entries[i])
	}

	indexInfo, err := common.GetIndexInfo()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}
	// workspace 信息
	workSpaceInfo, err := common.GetWorkSpaceInfo(common.CURRENT)
	if err != nil {
		fmt.Println("Error reading workspace:", err)
		return
	}
	// 参数文件通过对比 workspace 和 index 信息
	// 确定添加还是更新
	for _, name := range args {
		var workSpaceSha1, indexSha1 string
		var ok bool
		// 文件不在工作区
		if workSpaceSha1, ok = workSpaceInfo[name]; !ok {
			fmt.Println("file not exist:", name)
			continue
		}

		if indexSha1, ok = indexInfo[name]; !ok {
			// 新增文件到 index
			fmt.Println("add new file:", name)
			entry, err := common.NewIndexEntry(name)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			gitIndex.Entries = append(gitIndex.Entries, entry)
			for i := range gitIndex.Entries {
				fmt.Println(gitIndex.Entries[i])
			}
			err = gitIndex.WriteIndex()
			if err != nil {
				fmt.Println("Error writing index:", err)
				continue
			}
		} else {
			if workSpaceSha1 != indexSha1 {
				// 更新 index 文件
				fmt.Println("update file:", name)
				entry, err := common.NewIndexEntry(name)
				if err != nil {
					fmt.Println("Error reading file:", err)
					continue
				}
				gitIndex.UpdateEntry(entry)
				err = gitIndex.WriteIndex()
				if err != nil {
					fmt.Println("Error writing index:", err)
					continue
				}
			}
		}
	}
}
