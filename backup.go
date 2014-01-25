package main

import (
	"github.com/mmckeen/btrfs-backup/btrfs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func main() {
	// so that defers work as intended
	os.Exit(realMain())
}

// realMain is executed from main and returns the exit status to exit with.
func realMain() int {
	// If there is no explicit number of Go threads to use, then set it
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Reset the log variables to minimize work in the subprocess
	os.Setenv("BTRFS_BACKUP_LOG", "")
	os.Setenv("BTRFS_BACKUP_LOG_FILE", "")

	err := process()

	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("%s", err)
		return 1
	}

	return 0
}

// Do the majority of the application work,
// spawn off backup jobs, the like
//
// returns exit status for program
func process() error {

	// get default config values
	backupConfig := btrfs.DefaultConfig()

	backupConfig = btrfs.Config{"/home", ".snapshots"}

	// TODO: parse command line args

	// create drivers
	btrfs_driver := new(btrfs.Btrfs)

	// validate
	err := validateConfig(backupConfig, btrfs_driver)

	if err != nil {
		return err
	}

	return nil
}

// validate the config object
func validateConfig(backupConfig btrfs.Config, driver *btrfs.Btrfs) error {

	// create b

	// check to see if subvolume exists
	// do other sanity checks
	err := driver.Prepare(backupConfig)
	if err != nil {
		return err
	}

	subvols, err := driver.Subvolumes(backupConfig)
	if err != nil {
		return err
	}

	for i := 0; i < len(subvols); i++ {
		log.Printf("%s", subvols[i])
	}

	_, err2 := driver.Snapshot(backupConfig, "/")
	if err2 != nil {
		return err2
	}

	return nil
}

// Print some basic application info
func info() {
	log.SetOutput(os.Stderr)

	log.Printf("Btrfs Backup Target OS/Arch: %s %s", runtime.GOOS, runtime.GOARCH)
	log.Printf("Built with Go Version: %s", runtime.Version())
}

// logOutput determines where we should send logs (if anywhere).
func logOutput() (logOutput io.Writer, err error) {
	logOutput = ioutil.Discard
	if os.Getenv("BTRFS_BACKUP_LOG") != "" {
		logOutput = os.Stderr

		if logPath := os.Getenv("BTRFS_BACKUP_LOG_PATH"); logPath != "" {
			var err error
			logOutput, err = os.Create(logPath)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}
