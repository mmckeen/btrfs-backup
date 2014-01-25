package btrfs

func DefaultConfig() Config {
	return Config{"/",
		".snapshots"}
}

type Config struct {
	subvolume          string
	subvolumeDirectory string
}

// da getter method de Subvolume
func (c *Config) Subvolume() string {
	return DefaultConfig().subvolume
}
