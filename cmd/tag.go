/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
// 展示所有标签或创建标签
var tagCmd = &cobra.Command{
	Use:   "tag [NAME] [SHA]",
	Short: "list all tags or create a tag",
	Long: `list all tags when NAME and SHA are not provided, 
create a tag named NAME to HEAD when SHA is not provided, 
create a tag named NAME to the specified commit when NAME and SHA are provided.`,
	Run: TagFunc,
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func TagFunc(cmd *cobra.Command, args []string) {
	if !common.RepositoryExist(common.CURRENT) {
		fmt.Println("Not a git repository")
		return
	}
	if len(args) == 0 {
		common.ListAllTag()
	} else if len(args) == 1 {
		common.NewTag(args[0], "")
	} else if len(args) == 2 {
		common.NewTag(args[0], args[1])
	} else {
		fmt.Println("tag cmd receive tag [NAME] [SHA]")
	}
}
