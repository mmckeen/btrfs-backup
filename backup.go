package main

import (
	"flag"
	"github.com/mmckeen/btrfs-backup/btrfs"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
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

	// parse command line args
	subvolume_source := flag.String("subvolume", btrfs.DefaultConfig().SubvolumePath, "Subvolume to back up.")
	subvolume_destination_directory := flag.String("destination_subvolume", btrfs.DefaultConfig().SubvolumeDirectoryPath,
		"A relative path off of the subvolume path that will come to store snapshots.")
	server := flag.Bool("server", btrfs.DefaultConfig().Server, "Whether to enable listening as a backup server.")
	flag.Parse()

	log.Printf("%s", *server)
	// header info
	info()

	// set backup configuration
	backupConfig := btrfs.Config{*subvolume_source, *subvolume_destination_directory, *server}

	// create drivers
	btrfs_driver := new(btrfs.Btrfs)
	btrfs_driver.BackupConfig = backupConfig

	// validate
	err := validateConfig(backupConfig, btrfs_driver)

	if err != nil {
		return err
	}

	// start server if asked
	RPC := new(btrfs.BtrfsRPC)
	RPC.Driver = btrfs_driver

	if backupConfig.Server {
		rpc.Register(RPC)
		rpc.HandleHTTP()

		l, e := net.Listen("tcp", ":1234")
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)

	} else {
		// otherwise we are a client.  Query the client for a list of snapshots to send!
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}

		// Synchronous call
		subvols, err := btrfs_driver.Subvolumes(backupConfig)
		args := btrfs.Args{subvols}
		var reply []string
		err = client.Call("BtrfsRPC.SnapshotsNeeded", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}

		for i := 0; i < len(subvols); i++ {
			log.Printf("%s\n", subvols[i])
		}

	}

	return nil
}

// validate the config object
func validateConfig(backupConfig btrfs.Config, driver *btrfs.Btrfs) error {

	// check to see if subvolume exists
	// do other sanity checks
	err := driver.Prepare(backupConfig)
	if err != nil {
		return err
	}

	// do initial testing of system by listing subvolumes
	// and perform an initial snapshot for purposes of use later
	subvols, err := driver.Subvolumes(backupConfig)
	if err != nil && subvols == nil {
		return err
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
