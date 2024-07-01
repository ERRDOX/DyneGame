package main

import "fmt"

func main() {
	test := "test"
	switch test {
	case "test":
		fmt.Println("this is test")
		fallthrough
	case "other test", "another test", "with test":
		fmt.Println("this is other test or another test or with test")
		fallthrough
	default:
		fmt.Println("this is default")
	}
}
