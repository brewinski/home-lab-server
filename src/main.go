package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("This is an example go project...")
	time.Sleep(10 * time.Second)
	os.Exit(2)
}
