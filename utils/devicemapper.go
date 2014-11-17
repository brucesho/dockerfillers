package utils

import (
	"fmt"
	"os"
	"syscall"
)

func DeviceMapperMount(driverRootDir, imageId, mountPoint string) error {
	deviceFile, err := deviceMapperGetDeviceFile(driverRootDir, imageId)
	if err != nil {
		return err
	}

	deviceFileExists, err := fileExists(deviceFile)
	if err != nil {
		return err
	}

	if !deviceFileExists {
		fmt.Printf("Creating device file: %s\n", deviceFile)
	}

	return nil
}

func deviceMapperGetDeviceFile(driverRootDir, imageId string) (string, error) {
	st, err := os.Stat(driverRootDir)
	if err != nil {
		return "", err
	}
	sysSt := st.Sys().(*syscall.Stat_t)
	// "reg-" stands for "regular file".
	// In the future we might use "dev-" for "device file", etc.
	// docker-maj,min[-inode] stands for:
	//	- Managed by docker
	//	- The target of this device is at major <maj> and minor <min>
	//	- If <inode> is defined, use that file inside the device as a loopback image. Otherwise use the device itself.
	devicePrefix := fmt.Sprintf("docker-%d:%d-%d", major(sysSt.Dev), minor(sysSt.Dev), sysSt.Ino)
	return fmt.Sprintf("/dev/mapper/%s-%s", devicePrefix, imageId), nil
}

func major(device uint64) uint64 {
	return (device >> 8) & 0xfff
}

func minor(device uint64) uint64 {
	return (device & 0xff) | ((device >> 12) & 0xfff00)
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
