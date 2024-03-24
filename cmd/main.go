package main

import (
	"fmt"

	"github.com/solumD/go-social-media-api/cmd/server"
)

// Приветствие
func Greeting() int {
	fmt.Println("Hello!\n1 - Server\n2 - Client")
	fmt.Println("Enter the number of the program you want to run into the console.")
	var n int
	fmt.Print("Number: ")
	fmt.Scan(&n)
	fmt.Println("")
	return n
}

func main() {

	// выбор формата запуска
	for {
		num := Greeting()
		if num == 1 {
			fmt.Println("Starting server")
			server.Server()
			break
		} else if num == 2 {
			fmt.Println("Starting client ")
			break
		} else {
			fmt.Println("Write the correct number")
			continue
		}
	}
}
