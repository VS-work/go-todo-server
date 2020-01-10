package main

func main() {
	a := App{}
	a.Initialize("./todos.db")

	a.Run(":8080")
}
