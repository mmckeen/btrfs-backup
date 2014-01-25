package btrfs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// A library responsible for abstracting out the various btrfs-tools commands
// so that one can easily retrieve a list of subvolumes, create snapshots, and
// do other tasks

type Btrfs struct {
	subvolume string
}

func (d *Btrfs) Prepare(subvolume string) error {
	// test to see if subvolume is a valid btrfs file system

	var stderr bytes.Buffer
	cmd := exec.Command("btrfs", "subvolume", "show", subvolume)
	cmd.Stderr = &stderr

	log.Printf("Checking to see if valid subvolume: %s", subvolume)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error checking for a valid subvolume: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	// all good, is a valid subvolume
	log.Printf("Valid subvolume found: %s", subvolume)

	return nil

}
