/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
	"os"
	"path/filepath"
)

// initCmd represents the init command
// 初始化一个新的空仓库
// todo: 当前仅支持初始化当前目录下的仓库
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new, empty repository.",
	Long:  `Initialize a new, empty repository.`,
	Run:   initFunc,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initFunc(cmd *cobra.Command, args []string) {

	if len(args) != 0 {
		fmt.Println("init command does not accept any arguments")
		return
	}

	if err := common.InitRepo(common.CURRENT); err != nil {
		fmt.Println("fail to init:", err)
		os.RemoveAll(filepath.Join(".", common.REPO))
		return
	}
}
