package main

// @title API
// @version 1.0
// @description API библиотеки

// @host localhost:8080
// @BasePath /

// main runs the server on the given address.
func main() {
	NewServer(":8080").Serve()
}
