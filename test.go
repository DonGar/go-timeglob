package main

import (
	"fmt"
	"github.com/DonGar/go-timeglob/timeglob"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("TimeGlob is required.")
		os.Exit(1)
	}

	args := os.Args[1:]
	glob := strings.Join(args, " ")

	fmt.Printf("Using TimeGlob: %s\n", glob)

	tg, err := timeglob.Parse(glob)
	if err != nil {
		fmt.Printf(" error: %s\n", err.Error())
		os.Exit(1)
	}

	ticker := tg.Ticker()

	for {
		tick := <-ticker.C
		fmt.Printf(" tick: %s\n", tick)
	}
}
