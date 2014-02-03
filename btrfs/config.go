package btrfs

func DefaultConfig() Config {
	return Config{"/",
		".snapshots",
		false,
		"localhost",
		8000}
}

type Config struct {
	SubvolumePath          string
	SubvolumeDirectoryPath string
	Server                 bool
	DestinationHost        string
	DestinationPort        int
}

func (c *Config) Subvolume() string {
	return c.SubvolumePath
}

func (c *Config) SubvolumeDirectory() string {
	return c.SubvolumeDirectoryPath
}
