// +build linux darwin

package parse

import (
	"os"
	"path/filepath"

	"github.com/containers/storage/pkg/unshare"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/devices"
	"github.com/pkg/errors"
)

func DeviceFromPath(device string) ([]configs.Device, error) {
	var devs []configs.Device
	src, dst, permissions, err := Device(device)
	if err != nil {
		return nil, err
	}
	if unshare.IsRootless() && src != dst {
		return nil, errors.Errorf("Renaming device %s to %s is not supported in rootless containers", src, dst)
	}
	srcInfo, err := os.Stat(src)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting info of source device %s", src)
	}

	if !srcInfo.IsDir() {

		dev, err := devices.DeviceFromPath(src, permissions)
		if err != nil {
			return nil, errors.Wrapf(err, "%s is not a valid device", src)
		}
		dev.Path = dst
		devs = append(devs, *dev)
		return devs, nil
	}

	// If source device is a directory
	srcDevices, err := devices.GetDevices(src)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting source devices from directory %s", src)
	}
	for _, d := range srcDevices {
		d.Path = filepath.Join(dst, filepath.Base(d.Path))
		d.Permissions = permissions
		devs = append(devs, *d)
	}
	return devs, nil
}
