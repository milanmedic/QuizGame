package main

import (
	"bufio"
	"encoding/csv"
	"flag"
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
	fmt.Println("Quiz Game")

	filename := flag.String("file", "questions.csv", "A csv file in the format of 'question,answer'")
	timeLimit := flag.Int("time", 2, "The time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*filename)
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

			timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

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
