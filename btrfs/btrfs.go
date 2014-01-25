package btrfs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// A library responsible for abstracting out the various btrfs-tools commands
// so that one can easily retrieve a list of subvolumes, create snapshots, and
// do other tasks

type Btrfs struct{}

func (d *Btrfs) Prepare(config Config) error {
	// test to see if subvolume is a valid btrfs file system

	var stderr bytes.Buffer
	cmd := exec.Command("btrfs", "subvolume", "show", config.Subvolume())
	cmd.Stderr = &stderr

	log.Printf("Checking to see if valid subvolume: %s", config.Subvolume())
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error checking for a valid subvolume: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	// all good, is a valid subvolume
	log.Printf("Valid subvolume found: %s", config.Subvolume())

	return nil

}

func (d *Btrfs) Subvolumes(config Config) ([]string, error) {

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command("btrfs", "subvolume", "list", config.Subvolume())
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	log.Printf("Getting subvolume list for: %s", config.Subvolume())
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error getting subvolume list: %s\nStderr: %s",
			err, stderr.String())
		return nil, err
	}

	output := stdout.String()

	// partition output into slice
	subvols := strings.Split(output, "\n")[0 : len(strings.Split(output, "\n"))-1]

	return subvols, nil
}

func (d *Btrfs) Snapshot(config Config) (string, error) {
	return "", nil
}
