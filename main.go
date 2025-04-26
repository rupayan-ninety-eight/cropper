package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args...]\nCommands:\n  crop <folder> <split_width>\n  combine <folder> <direction>", os.Args[0])
	}

	command := os.Args[1]

	switch command {
	case "crop":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: %s crop <folder_path> <split_width>", os.Args[0])
		}
		folder := os.Args[2]
		width := os.Args[3]
		RunCropper(folder, width)

	case "combine":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: %s combine <folder_path> <left|right>", os.Args[0])
		}
		folder := os.Args[2]
		direction := os.Args[3]
		RunCombiner(folder, direction)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
