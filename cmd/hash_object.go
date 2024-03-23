/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-git/common"
)

// hashObjectCmd represents the hash-object command
// 计算文件的对象哈希，-w 选项写入对象数据库，-t 选项指定对象类型(默认为 blob)
// hash-object [-w] [-t TYPE] FILE
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object FILE",
	Short: "Compute the SHA-1 hash of an object",
	Long:  `Compute the SHA-1 hash of an object and optionally write it to the database`,
	Run:   hashObjectFunc,
}

var (
	objectType  string
	writeObject bool
	//toStdin     bool
)

func init() {
	// 读入选项
	hashObjectCmd.Flags().StringVarP(&objectType, "type", "t", "blob", "Specify the object type (e.g., blob)")
	hashObjectCmd.Flags().BoolVarP(&writeObject, "write", "w", false, "Write the object to the object database")
	// 暂时不支持 stdin
	//hashObjectCmd.Flags().BoolVarP(&toStdin, "stdin", "s", false, "Write the object to the stdin")

	rootCmd.AddCommand(hashObjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func hashObjectFunc(cmd *cobra.Command, args []string) {
	// 参数检验
	if len(args) != 1 {
		fmt.Println("Usage: go-git hash-object FILE")
		return
	}
	// 仓库检验
	if !common.RepositoryExist(".") {
		fmt.Println("Not a git repository")
		return
	}

	sha1, err := common.GetSha1ByPath(objectType, args[0], writeObject)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(sha1)
}
