package main

import (
	"fmt"
	"time"

	"tmelot.jsonparser/internal/repetitionTester"
)

func main() {
	rt := repetitionTester.NewRepetitionTester()
	fmt.Println(rt)

	for i := 0; i < 40; i++ {
		fmt.Print("\ri = ", i)
		time.Sleep(time.Duration(25) * time.Millisecond)
	}

	fmt.Println("\nDone!")
}
