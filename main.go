package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ANSI color codes
const (
	green  = "\033[42m\033[30m"
	yellow = "\033[43m\033[30m"
	gray   = "\033[100m\033[30m"
	reset  = "\033[0m"
)

// Emoji representations
const (
	greenEmoji  = "🟩"
	yellowEmoji = "🟨"
	grayEmoji   = "⬛"
)

// File to store game count
const gameCountFile = "gordle_count.txt"

func main() {
	// Get the game number from file
	gameNumber := getGameNumber()
	
	// Ask for word length
	fmt.Print("How many letters would you like the WORDLE word to be? (Default 5): ")
	var input string
	fmt.Scanln(&input)
	
	wordLength := 5 // Default value
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
	
	// Game loop
	for attempt := 1; attempt <= 6; attempt++ {
		fmt.Printf("Attempt %d: ", attempt)
		var guess string
		fmt.Scanln(&guess)
		guess = strings.ToLower(guess)
		
		if len(guess) != wordLength {
			fmt.Printf("Word must be of %d letters length!\n", wordLength)
			attempt--
			continue
		}
		
		// Compare guess with answer
		coloredGuess, feedback := generateFeedback(guess, answer)
		
		attempts = append(attempts, guess)
		feedbacks = append(feedbacks, feedback)
		
		// Display colorized guess
		fmt.Printf("%s %s\n", coloredGuess, feedback)
		
		if guess == answer {
			displayGameSummary(attempts, feedbacks, attempt, true, gameNumber)
			return
		}
	}
	
	displayGameSummary(attempts, feedbacks, 6, false, gameNumber)
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
	
	body, err := ioutil.ReadAll(resp.Body)
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
func generateFeedback(guess, answer string) (string, string) {
	coloredGuess := ""
	feedback := ""
	
	// First check for exact matches to handle duplicate letters correctly
	exactMatches := make(map[int]bool)
	for i := 0; i < len(guess); i++ {
		if guess[i] == answer[i] {
			exactMatches[i] = true
		}
	}
	
	// Count remaining occurrences of each letter in the answer
	remainingLetters := make(map[byte]int)
	for i := 0; i < len(answer); i++ {
		if !exactMatches[i] {
			remainingLetters[answer[i]]++
		}
	}
	
	// Generate feedback
	for i := 0; i < len(guess); i++ {
		if guess[i] == answer[i] {
			// Exact match
			coloredGuess += green + string(guess[i]) + reset
			feedback += greenEmoji
		} else if remainingLetters[guess[i]] > 0 {
			// Letter exists in the word but in wrong position
			coloredGuess += yellow + string(guess[i]) + reset
			feedback += yellowEmoji
			remainingLetters[guess[i]]--
		} else {
			// Letter not in the word
			coloredGuess += gray + string(guess[i]) + reset
			feedback += grayEmoji
		}
	}
	
	return coloredGuess, feedback
}

// Get and increment the game number
func getGameNumber() int {
	var gameNumber int = 1
	
	// Check if the file exists
	if _, err := os.Stat(gameCountFile); err == nil {
		// Read the current game count
		data, err := os.ReadFile(gameCountFile)
		if err == nil {
			if num, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
				gameNumber = num
			}
		}
	}
	
	// Increment and save for next time
	err := os.WriteFile(gameCountFile, []byte(strconv.Itoa(gameNumber+1)), 0644)
	if err != nil {
		fmt.Println("Warning: Could not save game count:", err)
	}
	
	return gameNumber
}

// Display game summary
func displayGameSummary(attempts []string, feedbacks []string, rounds int, won bool, gameNumber int) {
	result := "X/6"
	if won {
		result = fmt.Sprintf("%d/6", rounds)
	}
	
	fmt.Printf("\nWordle %d %s\n\n", gameNumber, result)
	
	for i, attempt := range attempts {
		fmt.Printf("%s : %s\n", attempt, feedbacks[i])
	}
	fmt.Println()
}