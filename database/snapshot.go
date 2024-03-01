package database

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type SnapshotInfo struct {
	ID      string    `json:"id"`
	Size    int32     `json:"size"`
	Created time.Time `json:"created"`
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
			Size:    *snapshot.AllocatedStorage,
		}
		snapshots = append(snapshots, snapshotInfo)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Created.Before(snapshots[j].Created)
	})

	return snapshots, nil
}
