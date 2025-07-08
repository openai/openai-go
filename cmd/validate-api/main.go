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
	// Command line options
	var (
		apiKey  = flag.String("api-key", "", "OpenAI API key (overrides OPENAI_API_KEY)")
		timeout = flag.Duration("timeout", 10*time.Second, "Request timeout")
		verbose = flag.Bool("verbose", false, "Show request details")
	)
	flag.Parse()

	// Get API key
	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "No API key provided")
		fmt.Fprintln(os.Stderr, " Use:")
		fmt.Fprintln(os.Stderr, "   • --api-key 'sk-...'")
		fmt.Fprintln(os.Stderr, "   • or export OPENAI_API_KEY='sk-...'")
		os.Exit(1)
	}

	// Create client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithRequestTimeout(*timeout),
	)

	// Verbose mode
	if *verbose {
		fmt.Printf("Testing OpenAI API connection...\n")
		fmt.Printf("Timeout: %v\n", *timeout)
		displayKey := key
		if len(key) > 7 {
			displayKey = key[:7]
		}
		fmt.Printf("🔑 API Key: %s...\n", displayKey)
		fmt.Println()
	}

	// Test connection
	if err := validateAPI(ctx, client); err != nil {
		handleError(err)
		os.Exit(1)
	}

	fmt.Println("✅ API connection successful")
	if *verbose {
		fmt.Println("🎉 Your OpenAI configuration is working correctly!")
	}
}

func validateAPI(ctx context.Context, client openai.Client) error {
	// Simple test: list models
	_, err := client.Models.List(ctx)
	return err
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "❌ API Error: %v\n\n", err)

	if apiErr, ok := err.(*apierror.Error); ok {
		switch apiErr.StatusCode {
		case 401:
			fmt.Fprintln(os.Stderr, "🔐 Authentication issue:")
			fmt.Fprintln(os.Stderr, "   • Verify your API key is correct")
			fmt.Fprintln(os.Stderr, "   • Check at https://platform.openai.com/api-keys")
			fmt.Fprintln(os.Stderr, "   • Ensure your account has available credits")
		case 403:
			fmt.Fprintln(os.Stderr, "🚫 Access denied:")
			fmt.Fprintln(os.Stderr, "   • Check your API key permissions")
			fmt.Fprintln(os.Stderr, "   • Verify your subscription plan")
		case 429:
			fmt.Fprintln(os.Stderr, "⏱️  Rate limit exceeded:")
			fmt.Fprintln(os.Stderr, "   • Wait a few minutes")
			fmt.Fprintln(os.Stderr, "   • Check your quota at https://platform.openai.com/usage")
		case 500, 502, 503, 504:
			fmt.Fprintln(os.Stderr, "🌐 OpenAI server issue:")
			fmt.Fprintln(os.Stderr, "   • OpenAI service is temporarily unavailable")
			fmt.Fprintln(os.Stderr, "   • Try again in a few minutes")
		default:
			fmt.Fprintf(os.Stderr, "�� Error code: %d\n", apiErr.StatusCode)
			fmt.Fprintf(os.Stderr, "📝 Message: %s\n", apiErr.Message)
		}
	} else {
		errStr := err.Error()
		if strings.Contains(errStr, "connection refused") {
			fmt.Fprintln(os.Stderr, "🌐 Connection issue:")
			fmt.Fprintln(os.Stderr, "   • Check your internet connection")
			fmt.Fprintln(os.Stderr, "   • Verify your proxy/firewall settings")
		} else if strings.Contains(errStr, "timeout") {
			fmt.Fprintln(os.Stderr, "⏰ Connection timeout:")
			fmt.Fprintln(os.Stderr, "   • Check your internet connection")
			fmt.Fprintln(os.Stderr, "   • Increase timeout with --timeout 30s")
		} else {
			fmt.Fprintf(os.Stderr, "❓ Unknown error: %s\n", errStr)
		}
	}

	fmt.Fprintln(os.Stderr, "\n💡 For more help: https://platform.openai.com/docs/guides/error-codes")
}
