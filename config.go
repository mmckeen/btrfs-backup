package main

func defaultConfig() config {
	return config{"/",
		".snapshots"}
}

type config struct {
	subvolume          string
	subvolumeDirectory string
}

// da getter method de Subvolume
func (c *config) Subvolume() string {
	return defaultConfig().subvolume
}
