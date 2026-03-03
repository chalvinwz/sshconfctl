package prompt

import (
	"bufio"
	"fmt"
	"strings"
)

func Ask(scanner *bufio.Scanner, msg string) string {
	fmt.Print(msg)
	if !scanner.Scan() {
		return ""
	}
	return scanner.Text()
}

func AskYesNo(scanner *bufio.Scanner, question string, defaultYes bool) (bool, error) {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}
	for {
		ans := strings.ToLower(strings.TrimSpace(Ask(scanner, fmt.Sprintf("%s %s ", question, suffix))))
		if ans == "" {
			return defaultYes, nil
		}
		switch ans {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Please answer 'y' or 'n'.")
		}
	}
}
