package config

type Server struct {
	Port        int    `yaml:"port"`
	QueueSize   int    `yaml:"queue_size"`
	StoreFolder string `yaml:"store_folder"`
}

func (c *Server) Reset() *Server {
	if c.Port <= 0 {
		c.Port = 3800
	}
	if c.QueueSize <= 0 {
		c.QueueSize = 100
	}
	if len(c.StoreFolder) == 0 {
		cfgDir, err := tryGetConfigDir()
		if err == nil {
			c.StoreFolder = cfgDir
		}
	}

	return c
}
