package main

import (
	"GO-X/auth"
	"fmt"
)

func main() {
	uuid := auth.GenUuid()
	fmt.Printf("UUID : %s\n", uuid)
}
