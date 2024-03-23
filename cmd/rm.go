/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm NAME...",
	Short: "Remove files from the index",
	Long:  `Remove files from the index`,
	Run:   RmFunc,
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func RmFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go-git rm NAME...")
		return
	}
	if !common.RepositoryExist(common.CURRENT) {
		fmt.Println("not a repository")
		return
	}

	gitIndex, err := common.ReadIndex()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}
	newEntries := make([]*common.IndexEntry, 0)
	for i := range gitIndex.Entries {
		if !contain(gitIndex.Entries[i].Name, args) {
			newEntries = append(newEntries, gitIndex.Entries[i])
		}
	}
	gitIndex.Entries = newEntries
	err = gitIndex.WriteIndex()
	if err != nil {
		fmt.Println("Error writing index:", err)
		return
	}
}

func contain(name string, args []string) bool {
	for i := range args {
		if name == args[i] {
			return true
		}
	}
	return false
}
