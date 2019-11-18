package disk

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	hostAPI API
}

type API interface {
	IsDiskInitialized(diskID string) (bool, error)
	InitializeDisk(diskID string) error
	PartitionsExist(diskID string) (bool, error)
	CreatePartition(diskID string) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) ListDiskIDs(context context.Context, request *internal.ListDiskIDsRequest, version apiversion.Version) (*internal.ListDiskIDsResponse, error) {
	// TODO: auto-generated stub
	return nil, nil
}

func (s *Server) ListDiskLocations(context context.Context, request *internal.ListDiskLocationsRequest, version apiversion.Version) (*internal.ListDiskLocationsResponse, error) {
	// TODO: auto-generated stub
	return nil, nil
}

func (s *Server) PartitionDisk(context context.Context, request *internal.PartitionDiskRequest, version apiversion.Version) (*internal.PartitionDiskResponse, error) {
	response := &internal.PartitionDiskResponse{}
	diskID := request.DiskId

	log.Infof("Checking if disk %s is initialized", diskID)
	initialized, err := s.hostAPI.IsDiskInitialized(diskID)
	if err != nil {
		return response, err
	}
	if !initialized {
		log.Infof("Initializing disk %s", diskID)
		err = s.hostAPI.InitializeDisk(diskID)
		if err != nil {
			return response, err
		}
	} else {
		log.Infof("Disk %s already initialized", diskID)
	}

	log.Infof("Checking if disk %s is partitioned", diskID)
	paritioned, err := s.hostAPI.PartitionsExist(diskID)
	if err != nil {
		return response, err
	}
	if !paritioned {
		log.Infof("Creating partition on disk %s", diskID)
		err = s.hostAPI.CreatePartition(diskID)
		if err != nil {
			return response, err
		}
	} else {
		log.Infof("Disk %s already partitioned", diskID)
	}

	return response, nil
}

func (s *Server) Rescan(context context.Context, request *internal.RescanRequest, version apiversion.Version) (*internal.RescanResponse, error) {
	// TODO: auto-generated stub
	return nil, nil
}
