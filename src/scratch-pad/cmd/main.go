package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Print("Starting application...")
	limit := make(chan struct{}, 5)

	for {
		limit <- struct{}{}
		// httpGetRequest to `http://192.168.0.201/dashboard/#/udp/routers`
		go func() {
			resp, err := http.Get("http://192.168.0.208")
			fmt.Printf("%v, %v", resp, err)
			<-limit
		}()
	}
}
