package ec2_instance

import "strconv"

// FilterFunc is a type that defines a function that can be used to
// filter EC2 instances. The function should take an EC2Instance as input and return
// a boolean value indicating whether the instance matches the filter criteria.
// Return true if the EC2Instance meets the condition(s) specified in the function,
// false otherwise.
type FilterFunc func(instance EC2Instance) bool

// IsOld is a filter function that checks if an EC2Instance is older than one year.
// It converts the AMI_Age field of the instance into an integer and checks if it is greater than 365.
// Returns true if the instance is older than 1 year, false otherwise.
// Note: This function does not handle errors from the conversion of AMI_Age to integer.
// Ensure AMI_Age is properly formatted as an integer string.
func IsOld(instance EC2Instance) bool {
	age, _ := strconv.Atoi(instance.InstanceAge)
	return age > 90
}

func IsRunningFilter(instance EC2Instance) bool {
	return instance.Status == "running"
}

// Filter modifies the EC2InstanceManager's Instances slice in-place to only contain instances
// that satisfy the provided filter function. The filter function should take an EC2Instance as an
// argument and return a boolean indicating whether the instance meets the desired criteria.
// Instances for which the filter function returns true are kept, others are removed.
func (mgr *EC2InstanceManager) Filter(filterFunc FilterFunc) {
	j := 0

	for _, instance := range mgr.Instances {
		if filterFunc(instance) {
			mgr.Instances[j] = instance
			j++
		}
	}

	mgr.Instances = mgr.Instances[:j]
}
