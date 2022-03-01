package main

func FetchTodos(page, pageSize int) ([]Todo, int) {
	var todos []Todo
	totalTodosCount := 0
	database := GetDb()
	database.Model(&Todo{}).Count(&totalTodosCount)
	database.Select("id, title, description, created_at, updated_at").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").
		Find(&todos)

	return todos, totalTodosCount
}

func FetchPendingTodos(page, pageSize int, completed bool) ([]Todo, int) {
	var todos []Todo
	var totalTodosCount int
	database := GetDb()
	database.Model(Todo{}).Where(Todo{Completed: completed}).Count(&totalTodosCount)
	database.Select("id, title, description, created_at, updated_at").
		Where(Todo{Completed: completed}).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").
		Find(&todos)

	return todos, totalTodosCount
}

func DeleteAllTodosServices() {
	database := GetDb()
	database.Model(Todo{}).Delete(Todo{})
}

func FetchById(id uint) (todo Todo, err error) {
	database := GetDb()
	err = database.Model(&Todo{}).First(&todo, id).Error
	return
}

func DeleteTodoServices(todo *Todo) error {
	database := GetDb()
	return database.Delete(todo).Error
}

func CreateTodoServices(title, description string, completed bool) (todo Todo, err error) {
	database := GetDb()
	todo = Todo{Title: title, Description: description, Completed: completed}
	err = database.Create(&todo).Error
	return todo, err
}

func UpdateTodoServices(id uint, title, description string, completed bool) (todo Todo, err error) {
	todo, err = FetchById(id)
	if err != nil {
		return
	}

	todo.Title = title

	if description != "" {
		todo.Description = description
	}

	todo.Completed = completed

	database := GetDb()

	database.Save(&todo)

	return
}
