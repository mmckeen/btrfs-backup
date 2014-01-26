package btrfs

type Args struct {
	Snapshots []string
}

type Snapshot struct {
	name string
	port int
}

type BtrfsRPC struct {
	Driver *Btrfs
}

func (d *BtrfsRPC) ReceiveSnapshot(args *Snapshot, reply *bool) error {
	// a client calls this when it is about to send a snapshot over the wire,
	// providing the destination for that snapshot as well

	return nil
}

func (d *BtrfsRPC) SnapshotsNeeded(args *Args, reply *[]string) error {

	snapshots, err := d.Driver.Subvolumes(d.Driver.BackupConfig)

	if err != nil {
		*reply = make([]string, 1)
		return err
	}

	*reply = difference(args.Snapshots, snapshots)

	return nil
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
