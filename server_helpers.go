package main

import "fmt"

func getAllowedOrigins() []string {
	reactTsPort := 3000
	localUiDev := fmt.Sprintf("http://localhost:%d", reactTsPort)
	return []string{localUiDev}
}
