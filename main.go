package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/files"
	"github.com/xyproto/ollamaclient/v2"
)

const (
	versionString    = "translate 1.1.1"
	translationModel = "mixtral:instruct"
)

func main() {

	stdinText := "hello"
	if files.DataReadyOnStdin() {
		data, err := io.ReadAll(os.Stdin)
		if err == nil { // success
			stdinText = string(data)
		}
	} else if len(os.Args) > 1 {
		var xs []string
		for _, arg := range os.Args {
			if !strings.HasPrefix(arg, "-") {
				xs = append(xs, strings.TrimSpace(arg))
			}
		}
		if len(xs) > 0 {
			stdinText = strings.Join(xs, " ")
			stdinText = strings.TrimPrefix(stdinText, "\"")
			stdinText = strings.TrimSuffix(stdinText, "\"")
		}
	}

	// Extract the base language from the LANG environment variable
	locale := env.Str("LANG", "en_US")

	// Construct the prompt with the language's display name and input from STDIN
	prompt := fmt.Sprintf("Translate the following text to the locale %s. Only output the translated text and nothing else. Add no commentary! Add no explanations! Only generate the translation, and nothing else. Here is the full text that needs to be translated: %s", locale, stdinText)

	oc := ollamaclient.New()
	oc.ModelName = translationModel

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
	translatedText = strings.TrimPrefix(translatedText, "«")
	translatedText = strings.TrimSuffix(translatedText, "»")
	translatedText = strings.TrimSpace(translatedText)

	fmt.Printf("%s\n", translatedText)
}
