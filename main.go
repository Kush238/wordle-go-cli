package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ANSI color codes

const (
	green  = "\033[42m\033[30m"
	yellow = "\033[43m\033[30m"
	gray   = "\033[100m\033[30m"
	reset  = "\033[0m"
)

// Emoji representation

const (
	greenEmoji  = "🟩"
	yellowEmoji = "🟨"
	grayEmoji   = "⬛"
)

func main(){
	rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("How many letters would you like the WORDLE to be? (Default 5): ")
	var input string
	fmt.Scanln(&input)

	wordLength := 5
	if len(input) > 0 {
		if val, err := strconv.Atoi(input); err == nil && val > 0 {
			wordLength = val
		} else {
			fmt.Println("Invalid input, using default length of 5.")
		}
	}

	// Fetch a random word from Random Word API
	answer, err := getRandomWord(wordLength)
	if err != nil {
		fmt.Println("Error fetching word:", err)
		return
	} 

	fmt.Printf("Welcome to GORDLE! Guess the %d-letter word.\n", wordLength)

	var attempts []string
	var feedbacks []string

	for attempt := 1; attempt <= 6; attempt++ {
		fmt.Printf("Attempt %d: ", attempt)
		var guess string
		fmt.Scanln(&guess)
		guess = strings.ToLower(guess)

		if len(guess) != wordLength {
			fmt.Printf("Word must be of %d letters length!\n", wordLength)
			attempt --
			continue
		}

		coloredGuess, feedback := generateFeedback(guess, answer)

		attempts = append(attempts, guess)
		feedbacks = append(feedbacks, feedback)

		fmt.Printf("%s %s\n", coloredGuess, feedback)

		if guess == answer {
			displayGameSummary(attempts, feedbacks, attempt, true)
			return
		}
	
	}

	displayGameSummary(attempts, feedbacks, 6, false)
	fmt.Println("Game Over! The correct word was:", answer)

}

// Fetch a random word from Random Word API
func getRandomWord(length int) (string, error) {
	url := fmt.Sprintf("https://random-word-api.herokuapp.com/word?length=%d", length)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var words []string
	err = json.Unmarshal(body, &words)
	if err != nil {
		return "", err
	}

	if len(words) == 0 {
		return "", fmt.Errorf("No words returned from Random Word API")
	}

	return words[0], nil
}

// Generate feedback for a guess

func generateFeedback(guess, answer string) (string, string){
	coloredGuess := ""
	feedback := ""

	exactMatches := make(map[int]bool) 
	for i := 0; i < len(guess); i++{
		if guess[i] == answer[i] {
			exactMatches[i] = true
		}
	}

	remainingLetters := make(map[byte]int)
	for i := 0; i < len(answer); i++ {
		if !exactMatches[i]{
			remainingLetters[answer[i]]++
		}
	}

	for i := 0; i < len(guess); i++ {
		if guess[i] == answer[i]{
			coloredGuess += green + string(guess[i]) + reset
			feedback += greenEmoji
		} else if remainingLetters[guess[i]] > 0 {
			coloredGuess += yellow + string(guess[i]) + reset
			feedback += yellowEmoji
			remainingLetters[guess[i]]--
		} else {
			coloredGuess += gray + string(guess[i])+ reset 
			feedback += grayEmoji
		}
	}
	return coloredGuess, feedback
}

func displayGameSummary(attempts []string, feedbacks []string, rounds int, won bool) {
	result := "X/6"
	if won {
		result = fmt.Sprintf("%d/6", rounds)
	}

	gameNumber := rand.Intn(9000) + 1000

	fmt.Printf("\nWordle %d %s\n\n", gameNumber, result)

	for i, attempt := range attempts {
		fmt.Printf("%s : %s\n", attempt, feedbacks[i])
	}
	fmt.Println()
}
