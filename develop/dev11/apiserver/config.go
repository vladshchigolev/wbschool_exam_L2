package apiserver

// Config ...
type Config struct {
	BindAddr string `json:"bind_addr"`
	LogFile  string `json:"log_file"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogFile:  "data.log",
	}
}
