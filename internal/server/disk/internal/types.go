package internal

type DiskLocation struct {
	Adapter string
	Bus     string
	Target  string
	LUNID   string
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
