package disk

import (
	"fmt"
	"os/exec"
)

// Implements the OS API calls related to Disk Devices. All code here should be very simple
// pass-through to the OS APIs or cmdlets. Any logic around the APIs/cmdlet invocation
// should go in internal/server/filesystem/disk.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) IsDiskInitialized(diskID string) (bool, error) {
	cmd := fmt.Sprintf("Get-Disk -Number %s | Where partitionstyle -eq 'raw'", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking initialized status of disk %s: %v, %v", diskID, out, err)
	}
	if len(out) == 0 {
		// disks with raw initializtion not detected
		return true, nil
	}
	return false, nil
}

func (APIImplementor) InitializeDisk(diskID string) error {
	cmd := fmt.Sprintf("Initialize-Disk -Number %s -PartitionStyle MBR", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error initializing disk %s: %v, %v", diskID, out, err)
	}
	return nil
}

func (APIImplementor) PartitionsExist(diskID string) (bool, error) {
	cmd := fmt.Sprintf("Get-Partition | Where DiskNumber -eq %s", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking presence of partitions on disk %s: %v, %v", diskID, out, err)
	}
	if len(out) > 0 {
		// disk has paritions in it
		return true, nil
	}
	return false, nil
}

func (APIImplementor) CreatePartition(diskID string) error {
	cmd := fmt.Sprintf("New-Partition -DiskNumber %s -UseMaximumSize", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating parition on disk %s: %v, %v", diskID, out, err)
	}
	return nil
}
