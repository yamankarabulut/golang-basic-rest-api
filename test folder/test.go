package test

import (
	"fmt"
)

type Params struct {
	VersionStartValue float64 `json:"versionStartValue" `
	VersionEndValue   float64 `json:"versionEndValue" `
	Lang1             string  `json:"lang1"`
	Lang2             string  `json:"lang2"`
}

func main() {
	fmt.Println("hello world")
}
