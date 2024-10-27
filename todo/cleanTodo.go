package todo

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// clearDatabase 清空 todos 表並重置自增 ID
func clearDatabase(dbName string) error {
	// 使用 gorm 連接到 SQLite 資料庫
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 清空 todos 表
	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		return fmt.Errorf("failed to delete records from todos: %v", err)
	}
	fmt.Println("All records from 'todos' table deleted.")

	// 重置自增 ID
	if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='todos'").Error; err != nil {
		return fmt.Errorf("failed to reset auto-increment ID for todos: %v", err)
	}
	fmt.Println("Auto-increment ID reset for 'todos' table.")

	return nil
}

// CleanTodo 載入配置並清空資料庫
func CleanTodo() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	if err := clearDatabase(config.DBName); err != nil {
		fmt.Printf("Failed to clear database: %v\n", err)
	} else {
		fmt.Println("Database cleared successfully.")
	}
}
