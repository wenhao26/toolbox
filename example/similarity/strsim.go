package main

import (
	"fmt"

	"github.com/antlabs/strsim"
)

func main() {
	text1 := `This was a five star "Good morning sir." She curtsied and smiled sheepishly.`
	text2 := `She was a beautiful girl with an unusual personality. Coming late to work and lying her way out was her thing.`

	similarity := strsim.Compare(text1, text2)
	fmt.Printf("相似度: %.4f\n", similarity)
}