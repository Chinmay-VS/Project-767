package main

import "fmt"

func main() {

	for i := 0; i < 10; i++ {
		fmt.Println("Count ", i)
	}

	attempts := 0
	for attempts < 3 {
		fmt.Println("Attempts:- ", attempts)
		attempts++

	}

	urls := []string{

		"https://google.com",
		"https://github.com",
		"https://reddit.com",
	}

	for index, url := range urls {

		fmt.Printf("Site %d: %s\n", index+1, url)

	}

}





