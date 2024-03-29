package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/files"
	"github.com/xyproto/ollamaclient"
)

const versionString = "translate 1.0.1"

func main() {

	stdinText := "hello"
	if files.DataReadyOnStdin() {
		data, err := io.ReadAll(os.Stdin)
		if err == nil { // success
			stdinText = string(data)
		}
	}

	// Extract the base language from the LANG environment variable
	locale := env.Str("LANG", "en_US")

	// Construct the prompt with the language's display name and input from STDIN
	//prompt := "Translate the following text to the locale \"" + locale + "\" (and only output the translated text): " + stdinText
	prompt := "Translate the following text to the locale " + locale + " (and only output the translated text): " + stdinText

	oc := ollamaclient.NewWithModel("mixtral:instruct")

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(versionString)
		return
	}

	oc.Verbose = len(os.Args) > 1 && os.Args[1] == "-v"

	if oc.Verbose {
		fmt.Println("Prompt: " + prompt)
	}

	if err := oc.PullIfNeeded(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var sb strings.Builder
	lines := strings.Split(strings.TrimSpace(output), "\n")
	lastIndex := len(lines) - 1
	for i, line := range lines {
		if (i == 0 || i == lastIndex) && (strings.Contains(line, locale) || strings.Contains(strings.ReplaceAll(line, "\\", ""), locale) || strings.Contains(line, "translat") || (strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")"))) {
			continue
		}
		sb.WriteString(line + "\n")
	}
	translatedText := strings.TrimSpace(sb.String())

	fmt.Printf("%s\n", translatedText)
}
