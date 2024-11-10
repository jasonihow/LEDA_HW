/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gogo/todo"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Todo Worker",
	Long:  `Pop a job`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("worker called")
		redisClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		config, err := todo.LoadConfig()
		db, err := gorm.Open(sqlite.Open(config.DBName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		for {
			job, err := redisClient.BRPop(context.Background(), 0, "queue").Result()
			if err != nil {
				fmt.Println("err => ", err)
				continue
			}

			// Job is in the second item of the result slice
			jobData := job[1]

			fmt.Println("Job Data => ", jobData)

			// Parse jobData JSON into TODO struct
			var todoItem todo.TODO
			err = json.Unmarshal([]byte(jobData), &todoItem)
			if err != nil {
				fmt.Println("Error parsing job data:", err)
				continue
			}

			db.Create(&todoItem)
			fmt.Println("Added TODO to db:", todoItem)

		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
