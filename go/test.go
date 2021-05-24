package main

import "fmt"

func main() {

	var set = make(map[int]string, 0)
	var k, bo = set[100]
	fmt.Println( k, bo)
	set[100] = "1"
	k, bo = set[100]
	fmt.Print( k, bo)
}
