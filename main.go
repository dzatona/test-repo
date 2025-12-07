package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("HELLO from Thunderwind Enclave!")
	fmt.Println("Waiting 5 seconds...")
	time.Sleep(5 * time.Second)
	fmt.Println("Goodbye!")
}
