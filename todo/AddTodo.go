package todo

import (
	"fmt"
)

// AddTodo 添加一个新的 TODO 项目到数据库
func AddTodo(todo TODO) error {
	// 确保 db 已经初始化
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	// 添加 todo 项目到数据库
	if err := db.Create(&todo).Error; err != nil {
		return fmt.Errorf("failed to add todo: %w", err)
	}

	fmt.Println("Added TODO:", todo)
	return nil
}
