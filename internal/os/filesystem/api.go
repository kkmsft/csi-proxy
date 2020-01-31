package filesystem

import (
	// "fmt"
	"fmt"
	"k8s.io/klog"
	"os"
	// "os/exec"
	// "runtime"
)

// Implements the Filesystem OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/filesystem/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (APIImplementor) PathExists(path string) (bool, error) {
	return pathExists(path)
}

func (APIImplementor) Mkdir(path string) error {
	klog.Infof("going to create path: %s", path)
	return os.MkdirAll(path, 0755)
}

func (APIImplementor) Rmdir(path string, force bool) error {
	if force {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func (APIImplementor) LinkPath(src string, tgt string) error {
	klog.Infof("going to call the link path - src: (old)  %s, tgt: (new) %s", src, tgt)
	err := os.Symlink(src, tgt)
	if err != nil {
		klog.Info("Error from symlink: %+v", err)
		return err
	}
	return nil
}

func (APIImplementor) IsLikelyNotMountPoint(tgt string) (bool, error) {
	// TODO: Reuse the code in mount_windows under k8s.io/kubernetes/pkg/util/mount
	// This code is same except the pathExists usage.
	klog.Info("Here going to lstat")
	stat, err := os.Lstat(tgt)
	if err != nil {
		klog.Infof("Here returning: %v ", err)
		return true, err
	}

	klog.Info("going to check the link")
	// If its a link and it points to an existing file then its a mount point.
	if stat.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(tgt)
		if err != nil {
			return true, fmt.Errorf("readlink error: %v", err)
		}
		exists, err := pathExists(target)
		if err != nil {
			return true, err
		}
		return !exists, nil
	}

	return true, nil
}
