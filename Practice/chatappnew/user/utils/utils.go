package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetTerminalInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	return input
}
