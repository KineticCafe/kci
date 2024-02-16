// Package database provides AWS RDS helper methods.
package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type DatabaseInfo struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	MultiAZ          bool           `json:"multi_az"`
	SnapshotsEnabled bool           `json:"snapshots_enabled"`
	Snapshots        []SnapshotInfo `json:"snapshots"`
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
