package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"workzen-be/internal/ai"
)

func main() {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	client, err := ai.NewClient(addr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	cvText := "John Doe, Software Engineer with 5 years of experience in Go and Python."
	jdText := "Senior Backend Engineer with experience in Go."

	fmt.Println("Testing AnalyzeCV...")
	resp, err := client.AnalyzeCV(context.Background(), cvText)
	if err != nil {
		fmt.Printf("AnalyzeCV failed: %v\n", err)
	} else {
		fmt.Printf("Summary: %s\n", resp.Summary)
		fmt.Printf("Fit Score: %d\n", resp.FitScore)
	}

	fmt.Println("\nTesting MatchJob...")
	matchResp, err := client.MatchJob(context.Background(), cvText, jdText)
	if err != nil {
		fmt.Printf("MatchJob failed: %v\n", err)
	} else {
		fmt.Printf("Score: %d\n", matchResp.Score)
		fmt.Printf("Verdict: %s\n", matchResp.Verdict)
	}
}
