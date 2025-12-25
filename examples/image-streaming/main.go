package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nordlys-Labs/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	fmt.Println("Starting image streaming example...")

	stream := client.Images.GenerateStreaming(ctx, openai.ImageGenerateParams{
		Model:         openai.ImageModelGPTImage1,
		Prompt:        "A cute baby sea otter",
		N:             openai.Int(1),
		Size:          openai.ImageGenerateParamsSize1024x1024,
		PartialImages: openai.Int(3),
	})

	for stream.Next() {
		event := stream.Current()

		switch variant := event.AsAny().(type) {
		case openai.ImageGenPartialImageEvent:
			fmt.Printf("  Partial image %d/3 received\n", variant.PartialImageIndex+1)
			fmt.Printf("   Size: %d characters (base64)\n", len(variant.B64JSON))

			// Save partial image to file
			filename := fmt.Sprintf("partial_%d.png", variant.PartialImageIndex+1)
			if err := saveBase64Image(variant.B64JSON, filename); err != nil {
				panic(fmt.Errorf("failed to save partial image: %w", err))
			}
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("   💾 Saved to: %s\n", absPath)
		case openai.ImageGenCompletedEvent:
			fmt.Printf("\n✅ Final image completed!\n")
			fmt.Printf("   Size: %d characters (base64)\n", len(variant.B64JSON))

			// Save final image to file
			filename := "final_image.png"
			if err := saveBase64Image(variant.B64JSON, filename); err != nil {
				panic(fmt.Errorf("failed to save final image: %w", err))
			}
			absPath, _ := filepath.Abs(filename)
			fmt.Printf("   💾 Saved to: %s\n", absPath)

		default:
			fmt.Printf("Received unknown event type: %+v\n", event)
		}
	}

	if err := stream.Err(); err != nil {
		panic(fmt.Errorf("error during streaming: %w", err))
	}
}

func saveBase64Image(b64Data, filename string) error {
	imageData, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	if err := os.WriteFile(filename, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
