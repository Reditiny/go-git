/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// checkIgnoreCmd represents the checkIgnore command
// 检查路径是否被 .gitignore 忽略
var checkIgnoreCmd = &cobra.Command{
	Use:   "check-ignore PATH...",
	Short: "check path(s) against ignore rules.",
	Long:  ``,
	Run:   CheckIgnoreFunc,
}

func init() {
	rootCmd.AddCommand(checkIgnoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkIgnoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkIgnoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CheckIgnoreFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go-git check-ignore PATH...")
		return
	}
	// todo: 当前仅支持 .gitignore 文件里配置的全名匹配的规则
	ignoredMap, err := common.ReadIgnore()
	if err != nil {
		fmt.Println("Error reading ignore file:", err)
		return
	}
	for _, path := range args {
		if _, ok := ignoredMap[path]; ok {
			fmt.Println(path)
		}
	}
}
