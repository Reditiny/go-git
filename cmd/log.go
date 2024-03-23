/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
// 显示 commit 的提交历史
// log SHA
var logCmd = &cobra.Command{
	Use:   "log SHA",
	Short: "show commit history",
	Long:  `show commit history, which can be used to track the changes of the repository.`,
	Run:   LogFunc,
}

var short bool

func init() {
	logCmd.Flags().BoolVarP(&short, "short", "s", false, "short log")

	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func LogFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: go-git log SHA")
		return
	}
	common.ParseCommit(".", args[0], short)
}
