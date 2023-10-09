package config

import (
	"flag"
	"log"
)

// ------------------------------------------------------------
//	Struct of configuration vars
type Configuration struct {
	host    string //	Server addres
	baseURL string //	Base URL for create alias
}

//	This value is "true" when Configuration already initialized
var alreadyInitialized bool

// ------------------------------------------------------------
//	Server address get method
func (c *Configuration) Host() string {
	return c.host
}

// ------------------------------------------------------------
//	Base URL get method
func (c *Configuration) BaseURL() string {
	return c.baseURL
}

// ------------------------------------------------------------
//	Set all configuration vars to default values
//	IMPORTANT: Repeated method call will generate a fatal error and terminate the application
func (c *Configuration) initialize() {

	if alreadyInitialized {
		log.Fatal("Repeated method call")
	}
	alreadyInitialized = true

	flag.StringVar(&c.host, "a", "localhost:8080", "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80")
	flag.StringVar(&c.baseURL, "b", "http://localhost:8080", "Response base addres for alias URL.\n\tFor example: http://192.168.1.2")
	flag.Parse()

	if err := checkServerAddres(c.host); err != nil {
		c.host = hostDefault
	}

	if err := checkBaseURL(c.baseURL); err != nil {
		c.baseURL = baseURLDefault
	}
}
