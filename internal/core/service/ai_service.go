package service

import (
	"context"
	"errors"
	"workzen-be/internal/ai"
	pb "workzen-be/internal/grpc/ai"
)

var ErrAIServiceUnavailable = errors.New("AI service is not available")

type AIService interface {
	AnalyzeCV(ctx context.Context, cvText string) (*pb.CVAnalysisResponse, error)
	MatchJob(ctx context.Context, cvText, jdText string) (*pb.JobMatchResponse, error)
}

type aiService struct {
	aiClient *ai.Client
}

func NewAIService(aiClient *ai.Client) AIService {
	return &aiService{
		aiClient: aiClient,
	}
}

func (s *aiService) AnalyzeCV(ctx context.Context, cvText string) (*pb.CVAnalysisResponse, error) {
	if s.aiClient == nil {
		return nil, ErrAIServiceUnavailable
	}
	return s.aiClient.AnalyzeCV(ctx, cvText)
}

func (s *aiService) MatchJob(ctx context.Context, cvText, jdText string) (*pb.JobMatchResponse, error) {
	if s.aiClient == nil {
		return nil, ErrAIServiceUnavailable
	}
	return s.aiClient.MatchJob(ctx, cvText, jdText)
}
