//go:build linux || darwin

package tun

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Create creates a TUN device at the path specified.
func (t *Tun) Create(path string) error {
	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0751); err != nil {
		return err
	}

	const (
		major = 10
		minor = 200
	)
	dev := unix.Mkdev(major, minor)
	err := unix.Mknod(path, unix.S_IFCHR, int(dev))
	if err != nil {
		return fmt.Errorf("creating TUN device file node: %w", err)
	}

	fd, err := unix.Open(path, 0, 0)
	if err != nil {
		return fmt.Errorf("unix opening TUN device file: %w", err)
	}

	const nonBlocking = true
	err = unix.SetNonblock(fd, nonBlocking)
	if err != nil {
		return fmt.Errorf("setting non block to TUN device file descriptor: %w", err)
	}

	return nil
}
