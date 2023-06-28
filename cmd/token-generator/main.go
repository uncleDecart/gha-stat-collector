package main

import (
	"flag"
	"fmt"

	"github.com/uncleDecart/gha-stat-collector/pkg/token"
)

func main() {
	tokenLength := flag.Int("n", 64, "Length of generated token")
	flag.Parse()

	fmt.Println(token.GenerateToken(*tokenLength))
}
