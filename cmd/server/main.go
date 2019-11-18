package main

import (
	diskapi "github.com/kubernetes-csi/csi-proxy/internal/os/disk"
	filesystemapi "github.com/kubernetes-csi/csi-proxy/internal/os/filesystem"
	"github.com/kubernetes-csi/csi-proxy/internal/server"
	disksrv "github.com/kubernetes-csi/csi-proxy/internal/server/disk"
	filesystemsrv "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
	srvtypes "github.com/kubernetes-csi/csi-proxy/internal/server/types"
	flag "github.com/spf13/pflag"
)

var (
	kubeletCSIPluginsPath = flag.String("kubelet-csi-plugins-path", `C:\var\lib\kubelet\plugins`, "Absolute path of the Kubelet plugin directory in the host file system")
	kubeletPodPath        = flag.String("kubelet-pod-path", `C:\var\lib\kubelet\pods`, "Absolute path of the kubelet pod directory in the host file system")
)

func main() {
	flag.Parse()
	apiGroups, err := apiGroups()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(apiGroups...)
	if err := s.Start(nil); err != nil {
		panic(err)
	}
}

// apiGroups returns the list of enabled API groups.
func apiGroups() ([]srvtypes.APIGroup, error) {
	fssrv, err := filesystemsrv.NewServer(*kubeletCSIPluginsPath, *kubeletPodPath, filesystemapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}
	disksrv, err := disksrv.NewServer(diskapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}
	return []srvtypes.APIGroup{
		fssrv,
		disksrv,
	}, nil
}
