// cmd/validate-api/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/apierror"
	"github.com/openai/openai-go/option"
)

func main() {
	// --- Command line flags ---
	// These allow the user to override the API key, set a custom timeout, and enable verbose output.
	var (
		apiKey  = flag.String("api-key", "", "OpenAI API key (overrides OPENAI_API_KEY)")
		timeout = flag.Duration("timeout", 10*time.Second, "Request timeout")
		verbose = flag.Bool("verbose", false, "Show request details")
	)
	flag.Parse()

	// --- API Key Resolution ---
	// Priority: --api-key flag > OPENAI_API_KEY env var.
	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key == "" {
		// Fail fast if no API key is provided.
		fmt.Fprintln(os.Stderr, "No API key provided")
		fmt.Fprintln(os.Stderr, " Use:")
		fmt.Fprintln(os.Stderr, "   • --api-key 'sk-...'")
		fmt.Fprintln(os.Stderr, "   • or export OPENAI_API_KEY='sk-...'")
		os.Exit(1)
	}

	// --- Context with Timeout ---
	// Ensures the request does not hang indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// --- Client Initialization ---
	// The client is configured with the resolved API key and timeout.
	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithRequestTimeout(*timeout),
	)

	// --- Verbose Logging ---
	// Show diagnostic info if requested.
	if *verbose {
		fmt.Printf("Testing OpenAI API connection...\n")
		fmt.Printf("Timeout: %v\n", *timeout)
		displayKey := key
		if len(key) > 7 {
			displayKey = key[:7]
		}
		fmt.Printf("API Key: %s...\n", displayKey)
		fmt.Println()
	}

	// --- API Validation ---
	// The core check: attempt to list models. This is a lightweight endpoint and a good proxy for API health.
	if err := validateAPI(ctx, client); err != nil {
		handleError(err)
		os.Exit(1)
	}

	fmt.Println("Your OpenAI configuration is working correctly!")
	if *verbose {
		fmt.Println("🎉")
	}
}

// validateAPI attempts a simple API call to verify connectivity and authentication.
// Returns an error if the call fails for any reason.
func validateAPI(ctx context.Context, client openai.Client) error {
	// Simple test: list models (minimal permissions required, fast response)
	_, err := client.Models.List(ctx)
	return err
}

// handleError provides structured, actionable error messages for common API issues.
// This function distinguishes between API errors (with status codes) and generic/network errors.
func handleError(err error) {
	fmt.Fprintf(os.Stderr, "API Error: %v\n\n", err)

	if apiErr, ok := err.(*apierror.Error); ok {
		// Handle known API error codes with specific guidance.
		switch apiErr.StatusCode {
		case 401:
			fmt.Fprintln(os.Stderr, " Authentication issue:")
			fmt.Fprintln(os.Stderr, "   • Verify your API key is correct")
			fmt.Fprintln(os.Stderr, "   • Check at https://platform.openai.com/api-keys")
			fmt.Fprintln(os.Stderr, "   • Ensure your account has available credits")
		case 403:
			fmt.Fprintln(os.Stderr, " Access denied:")
			fmt.Fprintln(os.Stderr, "   • Check your API key permissions")
			fmt.Fprintln(os.Stderr, "   • Verify your subscription plan")
		case 429:
			fmt.Fprintln(os.Stderr, " Rate limit exceeded:")
			fmt.Fprintln(os.Stderr, "   • Wait a few minutes")
			fmt.Fprintln(os.Stderr, "   • Check your quota at https://platform.openai.com/usage")
		case 500, 502, 503, 504:
			fmt.Fprintln(os.Stderr, " OpenAI server issue:")
			fmt.Fprintln(os.Stderr, "   • OpenAI service is temporarily unavailable")
			fmt.Fprintln(os.Stderr, "   • Try again in a few minutes")
		default:
			// For unhandled status codes, print the code and message for debugging.
			fmt.Fprintf(os.Stderr, " Error code: %d\n", apiErr.StatusCode)
			fmt.Fprintf(os.Stderr, " Message: %s\n", apiErr.Message)
		}
	} else {
		// Handle network and unknown errors.
		errStr := err.Error()
		if strings.Contains(errStr, "Connection refused") {
			fmt.Fprintln(os.Stderr, " Connection issue:")
			fmt.Fprintln(os.Stderr, "   • Check your internet connection")
			fmt.Fprintln(os.Stderr, "   • Verify your proxy/firewall settings")
		} else if strings.Contains(errStr, "timeout") {
			fmt.Fprintln(os.Stderr, " Connection timeout:")
			fmt.Fprintln(os.Stderr, "   • Check your internet connection")
			fmt.Fprintln(os.Stderr, "   • Increase timeout with --timeout 30s")
		} else {
			// Catch-all for unexpected errors.
			fmt.Fprintf(os.Stderr, "Unknown error: %s\n", errStr)
		}
	}

	// Always provide a pointer to the official error documentation for further troubleshooting.
	fmt.Fprintln(os.Stderr, "\n For more help: https://platform.openai.com/docs/guides/error-codes")
}
