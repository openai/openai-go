package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

func main() {
	root := flag.String("root", ".", "repository root")
	feedURL := flag.String("feed-url", gosupport.DefaultReleaseFeed, "official Go JSON release feed")
	reportPath := flag.String("report", "", "optional path for the Markdown report")
	flag.Parse()

	result, err := gosupport.CheckRepository(
		*root,
		*feedURL,
		time.Now(),
		&http.Client{Timeout: 30 * time.Second},
	)
	if err != nil {
		result.Findings = append(result.Findings, err.Error())
	}
	markdown := result.Markdown()
	fmt.Print(markdown)

	if *reportPath != "" {
		if writeErr := os.WriteFile(*reportPath, []byte(markdown), 0o644); writeErr != nil {
			fmt.Fprintln(os.Stderr, writeErr)
			os.Exit(1)
		}
	}
	if len(result.Findings) > 0 {
		os.Exit(1)
	}
}
