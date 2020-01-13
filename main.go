package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}

	if len(os.Args) < 4 {
		fmt.Println("Usage: go-todo-server db_file_path allowed_origins notification_email")
		return
	}

	dbfile := os.Args[1]
	a.Initialize(dbfile)

	a.Run()
}
