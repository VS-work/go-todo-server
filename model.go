package main

import (
	"database/sql"
)

type Todo struct {
	ID        int    `json:"rowid" db:"rowid"`
	Priority  int    `json:"priority" db:"priority"`
	Content   string `json:"content" db:"content"`
	Completed int    `json:"completed" db:"completed"`
}

func (t *Todo) getTodo(db *sql.DB) error {
	return db.QueryRow("SELECT priority, content, completed FROM todos WHERE rowid=?", t.ID).
		Scan(&t.Priority, &t.Content, &t.Completed)
}

func (t *Todo) updateTodo(db *sql.DB) error {
	_, err := db.Exec("UPDATE todos SET priority=?, content=?, completed=? WHERE rowid=?",
		t.Priority, t.Content, t.Completed, t.ID)
	return err
}

func (t *Todo) deleteTodo(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todos WHERE rowid=?", t.ID)

	return err
}

func (t *Todo) createTodo(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO todos(content, completed, priority) VALUES(?, 0, 0)", t.Content)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ROWID()").Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func getTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT rowid, priority, content, completed FROM todos ORDER BY rowid DESC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Priority, &t.Content, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}
