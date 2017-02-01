package v2action

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
)

type ApplicationInstanceStatusState ccv2.ApplicationInstanceStatusState

type ApplicationInstance struct {
	// CPU is the instance's CPU utilization percentage.
	CPU float64

	// Disk is the instance's disk usage in bytes.
	Disk int

	// DiskQuota is the instance's allowed disk usage in bytes.
	DiskQuota int

	// ID is the instance ID.
	ID int

	// Memory is the instance's memory usage in bytes.
	Memory int

	// MemoryQuota is the instance's allowed memory usage in bytes.
	MemoryQuota int

	// State is the instance's state.
	State ApplicationInstanceStatusState

	// Uptime is the number of seconds the instance has been running.
	Uptime int

	// Details are arbitrary information about the instance.
	Details string

	// Since is the Unix time stamp that represents the time the instance was
	// created.
	Since float64
}

// ApplicationInstancesNotFoundError is returned when a requested application is not
// found.
type ApplicationInstancesNotFoundError struct {
	ApplicationGUID string
}

func (e ApplicationInstancesNotFoundError) Error() string {
	return fmt.Sprintf("Application instances '%s' not found.", e.ApplicationGUID)
}

func (instance ApplicationInstance) StartTime() time.Time {
	return time.Now().Add(-1 * time.Duration(instance.Uptime) * time.Second)
}

func (actor Actor) GetApplicationInstancesByApplication(guid string) ([]ApplicationInstance, Warnings, error) {
	var allWarnings Warnings

	ccAppInstanceStatuses, warnings, err := actor.CloudControllerClient.GetApplicationInstanceStatusesByApplication(guid)
	allWarnings = append(allWarnings, warnings...)

	if _, ok := err.(cloudcontroller.ResourceNotFoundError); ok {
		return nil, allWarnings, ApplicationInstancesNotFoundError{ApplicationGUID: guid}
	} else if err != nil {
		return nil, allWarnings, err
	}

	ccAppInstances, warnings, err := actor.CloudControllerClient.GetApplicationInstancesByApplication(guid)
	allWarnings = append(allWarnings, warnings...)

	if _, ok := err.(cloudcontroller.ResourceNotFoundError); ok {
		return nil, allWarnings, ApplicationInstancesNotFoundError{ApplicationGUID: guid}
	} else if err != nil {
		return nil, allWarnings, err
	}

	appInstances := []ApplicationInstance{}

	for id, appInstance := range ccAppInstanceStatuses {
		nextInstance := ApplicationInstance{
			CPU:         appInstance.CPU,
			Disk:        appInstance.Disk,
			DiskQuota:   appInstance.DiskQuota,
			ID:          id,
			Memory:      appInstance.Memory,
			MemoryQuota: appInstance.MemoryQuota,
			State:       ApplicationInstanceStatusState(appInstance.State),
			Uptime:      appInstance.Uptime,
		}

		// TODO: should we error if instance in instanceStatuses but not in instances map?
		if _, ok := ccAppInstances[id]; ok {
			nextInstance.Details = ccAppInstances[id].Details
			nextInstance.Since = ccAppInstances[id].Since
		}

		appInstances = append(appInstances, nextInstance)
	}

	return appInstances, allWarnings, err
}
