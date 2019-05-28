package main

import "fmt"

func main() {
	var a, b int
	var c = &b

	switch *c {
	case a:
		fmt.Println("A!")
	case b:
		fmt.Println("B!")
	default:
		fmt.Println("C!")
	}
}
