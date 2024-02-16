package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type RDSManager struct {
	Databases []DatabaseInfo
	Client    *rds.Client
}

// NewManagerWithClient creates a new RDSManager with a supplied aws client.
func NewManagerWithClient(client *rds.Client) *RDSManager {
	return &RDSManager{
		Databases: []DatabaseInfo{},
		Client:    client,
	}
}

// NewManager creates a new RDSManager using the default AWS config and client
func NewManager() (*RDSManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := rds.NewFromConfig(cfg)

	mgr := NewManagerWithClient(client)

	return mgr, nil
}
