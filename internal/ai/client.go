package ai

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "workzen-be/internal/grpc/ai"
)

const retryPolicy = `{
	"methodConfig": [{
		"name": [{"service": "ai.AIService"}],
		"waitForReady": true,
		"retryPolicy": {
			"MaxAttempts": 3,
			"InitialBackoff": "0.2s",
			"MaxBackoff": "2s",
			"BackoffMultiplier": 2,
			"RetryableStatusCodes": ["UNAVAILABLE"]
		}
	}]
}`

type Client struct {
	conn    *grpc.ClientConn
	service pb.AIServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewAIServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) AnalyzeCV(ctx context.Context, cvText string) (*pb.CVAnalysisResponse, error) {
	requestCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	return c.service.AnalyzeCV(requestCtx, &pb.AnalyzeCVRequest{CvText: cvText})
}

func (c *Client) MatchJob(ctx context.Context, cvText, jdText string) (*pb.JobMatchResponse, error) {
	requestCtx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	return c.service.MatchJob(requestCtx, &pb.MatchJobRequest{CvText: cvText, JdText: jdText})
}
