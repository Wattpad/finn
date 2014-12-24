package finn

// BaseQueue contains common methods used by the different queue implementations
type BaseQueue struct {
	config QueueConfig
}

// SetConfig sets up configuration values. Defaults fill in for empty values in the passed-in QueueConfig.
func (self *BaseQueue) SetConfig(config QueueConfig, defaults QueueConfig) {
	if config == nil {
		self.config = defaults
		return
	}

	for key, value := range defaults {
		if config[key] == "" {
			config[key] = value
		}
	}

	self.config = config
}
