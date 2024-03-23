/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
)

// checkoutCmd represents the checkout command
// 将 commit 的内容写入当前目录
// cat-file COMMIT_ID PATH
var checkoutCmd = &cobra.Command{
	Use:   "checkout COMMIT_ID",
	Short: "checkout commit to a directory",
	Long:  `checkout commit to a directory, the directory must not exist.`,
	Run:   CheckoutFunc,
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CheckoutFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: go-git checkout COMMIT_ID")
		return
	}
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}

	err := common.CheckoutCommit(common.GetCommitID(args[0]))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
