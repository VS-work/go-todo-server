package main

/*import (
	"os"
	"fmt"
)*/

func main() {
	/*if len(os.Args) > 1 {
		dbname := os.Args[1]
		fmt.Println(dbname)
    return	
	}*/

	a := App{}
	a.Initialize("./db/todos_test.db")

	a.Run(":8080")
}
