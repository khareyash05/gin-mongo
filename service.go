package main

import "context"

type WorkflowService interface {
	CalculateCoverage(ctx context.Context) error
	GeneratePRTests(ctx context.Context) error
}
