package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3/scripts/internal/gosupport"
)

func main() {
	root := flag.String("root", ".", "repository root")
	flag.Parse()

	findings := gosupport.CheckReleaseNote(*root)
	if len(findings) == 0 {
		fmt.Println("The generated release changelog contains the approved Go-version note.")
		return
	}

	fmt.Fprintln(os.Stderr, "The generated release is missing required Go-version communication:")
	for _, finding := range findings {
		fmt.Fprintf(os.Stderr, "  - %s\n", finding)
	}
	os.Exit(1)
}
