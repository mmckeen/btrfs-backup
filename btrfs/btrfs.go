package btrfs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// A library responsible for abstracting out the various btrfs-tools commands
// so that one can easily retrieve a list of subvolumes, create snapshots, and
// do other tasks

type Btrfs struct {
	BackupConfig Config
}

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

func (d *Btrfs) SendSubvolume(subvolume string) error {
	// receive should already be started on other host, proceded with the sending!

	subvolumes := strings.Split(subvolume, " ")

	subvolume = subvolumes[8]

	// make sure we only send back snapshots that we've made
	if !(strings.Contains(subvolume, d.BackupConfig.SubvolumeDirectoryPath+"/") && strings.Contains(subvolume, "btrfs_backup")) {
		return nil
	}

	// trim back even more to only include the actually snapshot name
	subvolumes = strings.Split(subvolume, "/")

	subvolume = subvolumes[1]

	if subvolume == d.BackupConfig.SubvolumeDirectory() {
		return nil
	}

	log.Printf("Sending snapshot: %s", subvolume)

	var stderr bytes.Buffer

	cmd := exec.Command("btrfs", "send", subvolume)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error sending subvolume: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	return nil
}

func (d *Btrfs) Snapshot(config Config, srcSnapshot string) (string, error) {

	snapshot_dir := config.Subvolume() + "/" + config.SubvolumeDirectory()

	t := time.Now()

	timestamp := t.Format("20060102150405")

	snapshot := snapshot_dir + "/btrfs_backup_" + timestamp

	log.Printf("Making sure of the source directory: %s", srcSnapshot)

	var stderr bytes.Buffer
	cmd := exec.Command("btrfs", "subvolume", "show", srcSnapshot)
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error checking for a valid subvolume: %s\nStderr: %s",
			err, stderr.String())
		return "", err
	}

	log.Printf("Valid subvolume found: %s", srcSnapshot)

	log.Printf("Making sure of the destination directory: %s", snapshot_dir)

	create_destination := false
	var stderr3 bytes.Buffer
	cmd = exec.Command("btrfs", "subvolume", "show", snapshot_dir)
	cmd.Stderr = &stderr3

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		create_destination = true
	}

	if create_destination == true {

		log.Printf("Error in verifying the destination subvolume %s, trying to create it", snapshot_dir)

		new_volume, err_create := d.CreateSubvolume(snapshot_dir)

		if new_volume == snapshot_dir && err_create == nil {
			log.Printf("New subvolume created: %s", snapshot_dir)
		} else {

			err_create = fmt.Errorf("Error creating new subvolume. \nStderr: %s",
				err_create)

			return "", err_create
		}

	}

	log.Printf("Valid subvolume found: %s", snapshot_dir)

	// all good, is a valid subvolume source and destination
	log.Printf("Creating snapshot in: %s", snapshot)

	var stderr2 bytes.Buffer

	cmd = exec.Command("btrfs", "subvolume", "snapshot", "-r", srcSnapshot, snapshot)
	cmd.Stderr = &stderr2

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error creating snapshot: %s\nStderr: %s",
			err, stderr2.String())
		return "", err
	}

	cmd = exec.Command("sync", snapshot_dir)
	cmd.Stderr = &stderr2

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error running filesystem sync: %s\nStderr: %s",
			err, stderr2.String())
		return "", err
	}

	log.Printf("Snapshot created: %s from %s", snapshot, srcSnapshot)

	return snapshot, nil
}

func (d *Btrfs) CreateSubvolume(subvolumeDir string) (string, error) {

	var stderr bytes.Buffer
	cmd := exec.Command("btrfs", "subvolume", "create", subvolumeDir)
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err2 := cmd.Wait(); err2 != nil {
		err2 = fmt.Errorf("Error while creating destination subvolume: %s\nStderr: %s",
			err2, stderr.String())
		return "", err2
	}

	return subvolumeDir, nil

}
