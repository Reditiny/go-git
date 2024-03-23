/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-git/common"

	"github.com/spf13/cobra"
)

// catFileCmd represents the cat-file command
// 根据 sha1 hash 值显示对象的内容
// cat-file TYPE SHA1
var catFileCmd = &cobra.Command{
	Use:   "cat-file TYPE SHA1",
	Short: "read content of a git object",
	Long:  `read content of a git object, which can be commit, tree, blob, tag TYPE.`,
	Run:   catFileFunc,
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func catFileFunc(cmd *cobra.Command, args []string) {
	// 参数检验
	if len(args) != 2 {
		fmt.Println("Usage: go-git cat-file TYPE HASH")
		return
	}
	objType, objSha1 := args[0], args[1]
	if !common.ValidType(objType) {
		fmt.Printf("Invalid object type \"%s\"\n", objType)
		return
	}
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}
	// 读取对象
	obj, err := common.ReadObject(common.CURRENT, objSha1)
	if err != nil {
		fmt.Println("Error reading object:", err)
		return
	}
	if obj.GetType() != objType {
		fmt.Println("Object type does not match")
		return
	}
	fmt.Println(string(obj.Content))
}
