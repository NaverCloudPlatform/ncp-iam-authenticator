package utils

import (
	"fmt"
	"io"
	"strings"
)

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func PrintLog(writer io.Writer, logs []string) {
	fmt.Fprintln(writer, "## Debug info")
	for _, v := range logs {
		fmt.Fprintln(writer, "# "+v)
	}
	fmt.Fprintln(writer, "")
}
