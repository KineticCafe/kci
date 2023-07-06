package ec2_instance

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Scan through all instances and set the IsSSM flag.
func (mgr *EC2InstanceManager) FetchSSMDetails() error {
	// TODO move this to struct as it should be shared
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	ssmOutput, err := ssmClient.DescribeInstanceInformation(ctx, &ssm.DescribeInstanceInformationInput{})
	if err != nil {
		return fmt.Errorf("cannot describe SSM instance information, %v", err)
	}

	ssmManaged := make(map[string]bool)
	for _, instance := range ssmOutput.InstanceInformationList {
		ssmManaged[*instance.InstanceId] = true
	}

	for i := range mgr.Instances {
		mgr.Instances[i].IsSSM = ssmManaged[mgr.Instances[i].ID]
	}

	return nil
}
