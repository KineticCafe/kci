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

		snapshotsResp, err := mgr.Client.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
			DBInstanceIdentifier: dbInstance.DBInstanceIdentifier,
		})
		if err != nil {
			return fmt.Errorf("Unable to list DB snapshots for %s: %v", *dbInstance.DBInstanceIdentifier, err)
		}

		for _, snapshot := range snapshotsResp.DBSnapshots {
			snapshotInfo := SnapshotInfo{
				ID:      *snapshot.DBSnapshotIdentifier,
				Created: *snapshot.SnapshotCreateTime,
			}
			dbInfo.Snapshots = append(dbInfo.Snapshots, snapshotInfo)
		}

		sort.Slice(dbInfo.Snapshots, func(i, j int) bool {
			return dbInfo.Snapshots[i].Created.Before(dbInfo.Snapshots[j].Created)
		})

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

// FetchAMIAge retrieves the AMI age of the given EC2 instance.
/*
func (instance *EC2Instance) FetchAMIAge(client *ec2.Client) error {
	describeImagesInput := &ec2.DescribeImagesInput{
		ImageIds: []string{instance.AMI_ID},
	}

	imagesResp, err := client.DescribeImages(context.Background(), describeImagesInput)
	if err != nil {
		return fmt.Errorf("failed to describe image %s, %v", instance.AMI_ID, err)
	}

	if len(imagesResp.Images) == 0 {
		return fmt.Errorf("no image found with id %s", instance.AMI_ID)
	}

	amiCreationDate, err := time.Parse(time.RFC3339, *imagesResp.Images[0].CreationDate)
	if err != nil {
		return fmt.Errorf("failed to parse creation date of image %s, %v", instance.AMI_ID, err)
	}

	instance.AMI_Age = fmt.Sprintf("%d", int(time.Since(amiCreationDate).Hours()/24))

	return nil
}

// FetchAMIAge recursively calls FetchAMIAge on all Instances.
func (mgr *EC2InstanceManager) FetchAMIAge() error {
	for i := range mgr.Instances {
		err := mgr.Instances[i].FetchAMIAge(mgr.Client)
		if err != nil {
			return fmt.Errorf("unable to scan AMI Ages: %w", err)
		}
	}

	return nil
}

// FetchInstances connects to an AWS and fetches descriptions of all
// ec2 instances. These instances will be available in the Instances
// field.
func (mgr *EC2InstanceManager) FetchInstances(filter string) error {
	// empty in case of multiple runs
	mgr.Instances = []EC2Instance{}

	var filters []types.Filter
	if filter != "" {
		filter = "*" + filter + "*"

		filters = append(filters, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{filter},
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	resp, err := mgr.Client.DescribeInstances(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to describe instances: %w", err)
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			var name string

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}

			instanceStruct := EC2Instance{
				ID:          *instance.InstanceId,
				Name:        name,
				AMI_ID:      *instance.ImageId,
				InstanceAge: fmt.Sprintf("%d", int(time.Since(*instance.LaunchTime).Hours()/24)),
				AMI_Age:     "N/A", // not scanned initially
				IsSSM:       false, // This needs actual check
				Status:      string(instance.State.Name),
				PublicIP:    aws.ToString(instance.PublicIpAddress),
				PrivateIP:   aws.ToString(instance.PrivateIpAddress),
			}
			mgr.Instances = append(mgr.Instances, instanceStruct)
		}
	}

	return nil
}

//manager.Filter(ec2_instance.RunningInstances)
*/
