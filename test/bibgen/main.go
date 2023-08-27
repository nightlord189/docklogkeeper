package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const linesBetweenPause = 50
const pause = 30 * time.Second

func main() {
	fmt.Println("start")
	for {
		fmt.Println("-----START-----")
		printFile()
		fmt.Println("-----END OF FILE-----")
	}
}

func printFile() {
	readFile, err := os.Open("bible.txt")

	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	lines := 0

	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
		lines++
		if lines%linesBetweenPause == 0 {
			fmt.Println("---PAUSE---")
			time.Sleep(pause)
			fmt.Println("---CONTINUE---")
		}
	}
}
