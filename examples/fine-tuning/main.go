package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Nordlys-Labs/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	fmt.Println("==> Uploading file")

	data, err := os.Open("./fine-tuning-data.jsonl")
	if err != nil {
		panic(err)
	}
	file, err := client.Files.New(ctx, openai.FileNewParams{
		File:    data,
		Purpose: openai.FilePurposeFineTune,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Uploaded file with ID: %s\n", file.ID)

	fmt.Println("Waiting for file to be processed")
	for {
		file, err = client.Files.Get(ctx, file.ID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("File status: %s\n", file.Status)
		if file.Status == "processed" {
			break
		}
		time.Sleep(time.Second)
	}

	fmt.Println("")
	fmt.Println("==> Starting fine-tuning")
	fineTune, err := client.FineTuning.Jobs.New(ctx, openai.FineTuningJobNewParams{
		Model:        openai.FineTuningJobNewParamsModelGPT3_5Turbo,
		TrainingFile: file.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Fine-tuning ID: %s\n", fineTune.ID)

	fmt.Println("")
	fmt.Println("==> Track fine-tuning progress:")

	events := make(map[string]openai.FineTuningJobEvent)

	for fineTune.Status == "running" || fineTune.Status == "queued" || fineTune.Status == "validating_files" {
		fineTune, err = client.FineTuning.Jobs.Get(ctx, fineTune.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println(fineTune.Status)

		page, err := client.FineTuning.Jobs.ListEvents(ctx, fineTune.ID, openai.FineTuningJobListEventsParams{
			Limit: openai.Int(100),
		})
		if err != nil {
			panic(err)
		}

		for i := len(page.Data) - 1; i >= 0; i-- {
			event := page.Data[i]
			if _, exists := events[event.ID]; exists {
				continue
			}
			events[event.ID] = event
			timestamp := time.Unix(int64(event.CreatedAt), 0)
			fmt.Printf("- %s: %s\n", timestamp.Format(time.Kitchen), event.Message)
		}

		time.Sleep(5 * time.Second)
	}
}
