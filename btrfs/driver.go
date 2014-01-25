package btrfs

// Driver interface, designed to enable addition
// of new storage drivers in the future

type Driver interface {

	// Prepares configuration and sets variables
	Prepare() error

	// return mount location of subvolume in question
	Subvolume() string

	// Return slice of subvolumes part of this filesystem
	// Returns results of `btrfs subvolume list <Subvolume()>`
	Subvolumes() []string

	// Create a new filesystem snapshot, will always
	// store the snapshot under <Subvolume>/.snapshots
	// returns snapshot location under success, nil if not
	Snapshot() (string, error)
}
