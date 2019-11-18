package internal

type DiskIDList struct {
	// List of node metrics.
	IDs []string
}

type DiskLocation struct {
	Adapter string
	Bus     string
	Target  string
	LUNID   string
}

type ListDiskIDsRequest struct {
}

type ListDiskIDsResponse struct {
	// Map of disk device IDs and SCSI IDs associated with each disk device
	Disk_IDs map[string]*DiskIDList
}

type ListDiskLocationsRequest struct {
}

type ListDiskLocationsResponse struct {
	// Map of disk device IDs and <adapter, bus, target, lun ID> associated with each disk device
	DiskLocations map[string]*DiskLocation
}

type PartitionDiskRequest struct {
	// Disk device ID of the disk to partition
	DiskId string
}

type PartitionDiskResponse struct {
}

type RescanRequest struct {
}

type RescanResponse struct {
}
