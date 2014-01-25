package btrfs

// Driver interface, designed to enable addition
// of new storage drivers in the future
type Driver interface {

	// Prepares configuration and sets variables
	Prepare(interface{}) error

	// Return slice of subvolumes part of this filesystem
	// Returns results of `btrfs subvolume list <Subvolume()>`
	Subvolumes(interface{}) ([]string, error)

	// Create a new filesystem snapshot, will always
	// store the snapshot under <Subvolume>/<SnapshotsDirectory>
	// returns snapshot location under success, error if not
	Snapshot(interface{}, string) (string, error)

	// Function to create a subbvolume of the specified directory
	// Assumes that the directory doesn't already exists and does
	// No error checking, the caller should do this.
	CreateSubvolume(string) (string, error)
}
