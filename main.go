package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
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

	csvReader := csv.NewReader(file)
	reader := bufio.NewReader(os.Stdin)
	answerChn := make(chan string)
	numOfCorrectAns := 0

	for {
		if record, err := csvReader.Read(); err != io.EOF {
			check(err)
			question := record[0]
			correctAnswer := record[1]

			timer := time.NewTimer(2 * time.Second)

			go func() {
				PromptQuestion(question)
				answerChn <- GetUserAns(reader)
			}()

			select {
			case <-timer.C:
				fmt.Println("Time's Up!")
				fmt.Println("The number of correct answers is", numOfCorrectAns)
				return
			case userAnswer := <-answerChn:
				if !CheckCorrectness(userAnswer, correctAnswer) {
					fmt.Println("The number of correct answers is", numOfCorrectAns)
					return
				} else {
					numOfCorrectAns++
				}
			}
		} else {
			fmt.Println("The number of correct answers is", numOfCorrectAns)
			return
		}
	}
}

func PromptQuestion(question string) {
	fmt.Println("------------------------------------------------")
	fmt.Println("Question: ")
	fmt.Println(question)
	fmt.Println("------------------------------------------------")
	fmt.Println("Answer:")
}

func GetUserAns(reader *bufio.Reader) string {
	userAnswer, _ := reader.ReadString('\n')
	// convert CRLF to LF
	userAnswer = strings.Replace(userAnswer, "\n", "", -1)
	return userAnswer
}

func CheckCorrectness(userAnswer, correctAnswer string) bool {
	return strings.Compare(correctAnswer, userAnswer) == 0
}
