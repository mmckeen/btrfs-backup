package btrfs

func DefaultConfig() Config {
	return Config{"/",
		".snapshots",
		false}
}

type Config struct {
	SubvolumePath          string
	SubvolumeDirectoryPath string
	Server                 bool
}

// da getter method de Subvolume
func (c *Config) Subvolume() string {
	return c.SubvolumePath
}

// da getter method de da subvolumeDirectory
func (c *Config) SubvolumeDirectory() string {
	return c.SubvolumeDirectoryPath
}
