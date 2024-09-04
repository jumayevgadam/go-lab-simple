package main

import "fmt"

// Let's create some struct named Config, now wanna show u
// pointer to struct is best practise with some simple examples
// Let's GO!!!!

// Config struct is
type Config struct {
	Host string
	Port int
}

// function to setdefaultconfiguration
func setdefaultconfiguration(config *Config) {
	config.Host = "localhost"
	config.Port = 9000
}

// function to override configuration
func ovverrideConfig(config *Config) {
	config.Host = "product.server.com"
	config.Port = 8000
}

func main() {
	config := Config{}

	// Setting default configuration
	setdefaultconfiguration(&config)
	fmt.Printf("Default Config: %+v\n", config)

	// Overriding the configuration
	ovverrideConfig(&config)
	fmt.Printf("Overridden config: %+v\n", config)
}
