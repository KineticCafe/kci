// Package database provides AWS RDS helper methods.
package database

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type SnapshotInfo struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
}

type DatabaseInfo struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	MultiAZ          bool           `json:"multi_az"`
	SnapshotsEnabled bool           `json:"snapshots_enabled"`
	Snapshots        []SnapshotInfo `json:"snapshots"`
}

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

// FetchDatabases connects to an AWS and fetches descriptions of all
// RDS databases.
func (mgr *RDSManager) Fetch(filter string) error {
	mgr.Databases = []DatabaseInfo{}

	resp, err := mgr.Client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("unable to describe databases: %w", err)
	}

	for _, dbInstance := range resp.DBInstances {
		if len(filter) > 0 && !strings.Contains(aws.ToString(dbInstance.DBInstanceIdentifier), filter) {
			continue
		}

		dbInfo := DatabaseInfo{
			Name:             aws.ToString(dbInstance.DBName),
			ID:               aws.ToString(dbInstance.DBInstanceIdentifier),
			SnapshotsEnabled: aws.ToInt32(dbInstance.BackupRetentionPeriod) > 0,
			MultiAZ:          *dbInstance.MultiAZ,
		}

		snapshots, err := mgr.FetchSnapshots(dbInfo.ID)
		if err != nil {
			return fmt.Errorf("could not load snapshots: %v", err)
		}

		dbInfo.Snapshots = snapshots

		mgr.Databases = append(mgr.Databases, dbInfo)
	}

	return nil
}

func (mgr *RDSManager) FetchSingle(identifier string) ([]DatabaseInfo, error) {
	return []DatabaseInfo{}, nil
}

func (mgr *RDSManager) FetchSnapshots(identifier string) ([]SnapshotInfo, error) {
	snapshots := []SnapshotInfo{}

	snapshotsResp, err := mgr.Client.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(identifier),
	})

	if err != nil {
		return snapshots, fmt.Errorf("Unable to list DB snapshots for %s: %v", identifier, err)
	}

	for _, snapshot := range snapshotsResp.DBSnapshots {
		snapshotInfo := SnapshotInfo{
			ID:      *snapshot.DBSnapshotIdentifier,
			Created: *snapshot.SnapshotCreateTime,
		}
		snapshots = append(snapshots, snapshotInfo)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Created.Before(snapshots[j].Created)
	})

	return snapshots, nil
}

func (db *DatabaseInfo) LatestSnapshot() (SnapshotInfo, error) {
	if !db.SnapshotsEnabled {
		return SnapshotInfo{}, fmt.Errorf("cannot return latest snapshot if snapshots disabled")
	}

	if len(db.Snapshots) < 0 {
		return SnapshotInfo{}, fmt.Errorf("cannot return latest snapshot if no snapshots registered")
	}

	return db.Snapshots[len(db.Snapshots)-1], nil
}

func (db *DatabaseInfo) LatestSnapshotID() string {
	if db.SnapshotsEnabled && len(db.Snapshots) > 0 {
		return db.Snapshots[len(db.Snapshots)-1].ID
	}

	return ""
}
