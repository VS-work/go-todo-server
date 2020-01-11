package main

import (
	"database/sql"
)

type todo struct {
	ID        int    `json:"rowid"`
	Priority  int    `json:"priority"`
	Content   string `json:"content"`
	Completed int    `json:"completed"`
}

func (t *todo) getTodo(db *sql.DB) error {
	return db.QueryRow("SELECT priority, content, completed FROM todos WHERE rowid=?", t.ID).
		Scan(&t.Priority, &t.Content, &t.Completed)
}

func (t *todo) updateTodo(db *sql.DB) error {
	_, err := db.Exec("UPDATE todos SET priority=?, content=?, completed=? WHERE rowid=?",
		t.Priority, t.Content, t.Completed, t.ID)
	return err
}

func (t *todo) deleteTodo(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todos WHERE rowid=?", t.ID)

	return err
}

func (t *todo) createTodo(db *sql.DB) error {
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

func getTodos(db *sql.DB) ([]todo, error) {
	rows, err := db.Query("SELECT rowid, priority, content, completed FROM todos ORDER BY rowid DESC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []todo{}

	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Priority, &t.Content, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}
