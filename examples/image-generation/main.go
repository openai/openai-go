package main

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()

	ctx := context.Background()

	prompt := "A cute robot in a forest of trees."

	print("> ")
	println(prompt)
	println()

	// GPT Image returns base64-encoded image data by default.
	image, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt: prompt,
		Model:  openai.ImageModelGPTImage1,
		N:      openai.Int(1),
	})
	if err != nil {
		panic(err)
	}
	println("Image Base64 Length:")
	println(len(image.Data[0].B64JSON))
	println()

	imageBytes, err := base64.StdEncoding.DecodeString(image.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	dest := "./image.png"
	println("Writing image to " + dest)
	err = os.WriteFile(dest, imageBytes, 0755)
	if err != nil {
		panic(err)
	}
}
