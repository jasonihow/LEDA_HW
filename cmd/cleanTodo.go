/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gogo/todo"

	"github.com/spf13/cobra"
)

// cleanTodoCmd represents the cleanTodo command
var cleanTodoCmd = &cobra.Command{
	Use:   "cleanTodo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cleanTodo called")
		todo.CleanTodo()
	},
}

func init() {
	rootCmd.AddCommand(cleanTodoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanTodoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanTodoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
