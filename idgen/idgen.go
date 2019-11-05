package main

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func main() {

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := node.Generate()
	fmt.Printf("Snow Flake ID hex: \n%x\n", id.Int64())
	fmt.Printf("Snow Flake ID: \n%d\n", id.Int64())
}
