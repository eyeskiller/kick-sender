package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// PrintBanner prints the application banner
func PrintBanner() {
	fmt.Printf("\n%s%s", Bold, Cyan)
	fmt.Println("╭────────────────────────────────────╮")
	fmt.Println("│         KICK INTERVAL              │")
	fmt.Println("│       Automated Chat Sender        │")
	fmt.Println("╰────────────────────────────────────╯")
	fmt.Printf("%s\n", Reset)
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Printf("%s✓%s %s\n", Green, Reset, msg)
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Printf("%s✗%s %s\n", Red, Reset, msg)
}

// PrintInfo prints an info message
func PrintInfo(msg string) {
	fmt.Printf("%s→%s %s\n", Blue, Reset, msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Printf("%s!%s %s\n", Yellow, Reset, msg)
}

// Prompt prompts the user for input
func Prompt(message string) string {
	fmt.Printf("%s%s%s\n", Bold, message, Reset)
	fmt.Print("> ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptYesNo prompts for a yes/no confirmation
func PromptYesNo(message string) bool {
	fmt.Printf("%s%s%s [Y/n] ", Bold, message, Reset)
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	
	return input == "" || input == "y" || input == "yes"
}

// PromptChoice prompts the user to select from a list of choices
func PromptChoice(message string, choices []string) int {
	fmt.Printf("%s%s%s\n", Bold, message, Reset)
	for i, choice := range choices {
		fmt.Printf("[%d] %s\n", i+1, choice)
	}
	fmt.Print("Select: ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	var choice int
	fmt.Sscanf(input, "%d", &choice)
	
	if choice < 1 || choice > len(choices) {
		return -1
	}
	return choice - 1
}

// PrintSendResult prints the result of a send attempt
func PrintSendResult(t time.Time, err error) {
	timestamp := t.Format("15:04:05")
	
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") {
			fmt.Printf("%s  %sRETRY%s      %s\n", timestamp, Yellow, Reset, err.Error())
		} else {
			fmt.Printf("%s  %sFAILED%s     %s\n", timestamp, Red, Reset, err.Error())
		}
	} else {
		fmt.Printf("%s  %sSENT%s       Message delivered\n", timestamp, Green, Reset)
	}
}

// PrintConfig prints configuration summary
func PrintConfig(account, channel, message, interval string) {
	fmt.Println("\nConfiguration:")
	fmt.Println("────────────────────────────────────")
	fmt.Printf("Account:  %s\n", account)
	fmt.Printf("Channel:  %s\n", channel)
	fmt.Printf("Message:  %s\n", message)
	fmt.Printf("Interval: %s\n", interval)
	fmt.Println("────────────────────────────────────")
}

// PrintRunning prints the running status
func PrintRunning(account, channel, interval string) {
	fmt.Printf("\n%s%sKICK INTERVAL — RUNNING%s\n\n", Bold, Cyan, Reset)
	fmt.Printf("Account:  %s\n", account)
	fmt.Printf("Channel:  %s\n", channel)
	fmt.Printf("Interval: %s\n\n", interval)
	fmt.Println("────────────────────────────────────")
}

// PrintSummary prints a summary of sent messages
func PrintSummary(sent, failed int) {
	fmt.Println("────────────────────────────────────")
	fmt.Printf("\nMessages sent: %d\n", sent)
	fmt.Printf("Failed attempts: %d\n\n", failed)
}
