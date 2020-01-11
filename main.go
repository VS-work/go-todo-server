package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}

	if len(os.Args) < 3 {
		fmt.Println("Usage: go-todo-server db_file_path allowed_origins")
		return
	}

	dbfile := os.Args[1]
	a.Initialize(dbfile)

	a.Run()
}
