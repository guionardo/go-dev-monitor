package config

type Server struct {
	Port        int    `yaml:"port"`
	QueueSize   int    `yaml:"queue_size"`
	StoreFolder string `yaml:"store_folder"`
}

func (c *Server) Reset() {
	if c.Port <= 0 {
		c.Port = 3800
	}
	if c.QueueSize <= 0 {
		c.QueueSize = 100
	}
	if len(c.StoreFolder) == 0 {

	}
}
