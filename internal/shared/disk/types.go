package shared

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
