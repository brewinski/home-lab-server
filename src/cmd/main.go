package main

import (
	"fmt"
	"net/http"
)

func main() {
	limit := make(chan struct{}, 50)

	for {
		limit <- struct{}{}
		// httpGetRequest to `http://192.168.0.201/dashboard/#/udp/routers`
		go func() {
			resp, err := http.Get("http://192.168.0.201/dashboard/#/udp/routers")
			fmt.Printf("%v, %v", resp, err)
			<-limit
		}()
	}
}
