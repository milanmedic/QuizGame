package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	fmt.Println("Hello Quiz Game")

	file, err := os.Open("questions.csv")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reader := bufio.NewReader(os.Stdin)
	numOfCorrectAns := 0

	for scanner.Scan() {
		entry := scanner.Text()
		entries := strings.Split(entry, ",")
		question, correctAnswer := entries[0], entries[1]

		fmt.Println("------------------------------------------------")
		fmt.Println("Question: ")
		fmt.Println(question)
		fmt.Println("------------------------------------------------")
		fmt.Println("Answer")
		userAnswer, _ := reader.ReadString('\n')
		// convert CRLF to LF
		userAnswer = strings.Replace(userAnswer, "\n", "", -1)
		if strings.Compare(correctAnswer, userAnswer) != 0 {
			file.Close()
			fmt.Println("The number of correct answers is", numOfCorrectAns)
			os.Exit(0)
		} else {
			numOfCorrectAns++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
