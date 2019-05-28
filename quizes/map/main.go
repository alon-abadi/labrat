package main

import "fmt"

func main() {
	quiz := map[interface{}]int{
		new(int):      1,
		new(int):      2,
		new(struct{}): 3,
		new(struct{}): 4,
	}

	fmt.Print(len(quiz))

	fmt.Print(fmt.Sprintf("%+v", quiz))
}
