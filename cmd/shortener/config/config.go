// Configuration appliation github.com/Schalure/urlalias
package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------------------------------
//
//	Application constants
const (
	AppName           = string("github.com/Schalure/urlalias") //	Application name
	hostEnvKey        = string("SERVER_ADDRESS")               //	key for "host" in environment variables
	baseURLEnvKey     = string("BASE_URL")                     //	key for "baseURL" in environment variables
	storageFileEnvKey = string("FILE_STORAGE_PATH")            //	key for "storageFile" in environment variables
)

// ------------------------------------------------------------
//
//	Default values
const (
	hostDefault        = string("localhost:8080")        //	Host default value
	baseURLDefault     = string("http://localhost:8080") //	Base URL default value
	storageFileDefault = "/tmp/short-url-db.json"        //	Default file name of URLs storage
	logToFileDefault   = false                           //	How to save log default value
)

// ------------------------------------------------------------
//
//	Struct of configuration vars
type Configuration struct {
	host           string //	Server addres
	baseURL        string //	Base URL for create alias
	storageFile    string // File name of URLs storage
	useStorageFile bool
	logToFile      bool //	true - save log to file, false - print log to console
}

// Common config variable
var config *Configuration

// ------------------------------------------------------------
//
//	Constructor of Config type
//	Output:
//		*Config
func NewConfig() *Configuration {

	if config != nil {
		return config
	}

	config = new(Configuration)

	//	Fill default values
	config.host = hostDefault
	config.baseURL = baseURLDefault
	config.logToFile = logToFileDefault

	config.parseFlags()
	config.parseEnv()

	log.Printf("Server address: \"%s\"\n", config.host)
	log.Printf("Base URL: \"%s\"\n", config.host)
	if config.useStorageFile {
		log.Printf("Storage file: \"%s\"\n", config.storageFile)
	}
	log.Printf("Save log to file: \"%t\"\n", config.logToFile)
	return config
}

// ------------------------------------------------------------
//
//	Getter "Configuration.host"
//	Output:
//		c.host string
func (c *Configuration) Host() string {
	return c.host
}

// ------------------------------------------------------------
//
//	Getter "Configuration.baseURL"
//	Output:
//		c.baseURL string
func (c *Configuration) BaseURL() string {
	return c.baseURL
}

func (c *Configuration) StorageFile() string {
	return c.storageFile
}

// ------------------------------------------------------------
//
//	Getter "Configuration.logSaver"
//	Output:
//		c.baseURL string
func (c *Configuration) LogToFile() bool {
	return bool(c.logToFile)
}

// ------------------------------------------------------------
//
//	Parse flags method of "Config" type
func (c *Configuration) parseFlags() {

	host := flag.String("a", hostDefault, "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80")
	baseURL := flag.String("b", baseURLDefault, "Response base addres for alias URL.\n\tFor example: http://192.168.1.2")
	storageFile := ""
	useStorageFile := false
	logToFile := flag.Bool("l", logToFileDefault, "Variant of logger: true - save log to file, false - print log to console")

	flag.Func("f", "File name of URLs storage. Specify the full name of the file", func(s string) error {

		if s == "" {
			storageFile = storageFileDefault
		} else {
			storageFile = s
		}
		useStorageFile = true

		return nil
	})

	flag.Parse()

	if err := checkServerAddres(*host); err == nil {
		c.host = *host
	}

	if err := checkBaseURL(*baseURL); err == nil {
		c.baseURL = *baseURL
	}

	c.logToFile = *logToFile

	if useStorageFile {
		c.storageFile = storageFile
	}
}

// ------------------------------------------------------------
//
//	Parse environment variables method of "Config" type
func (c *Configuration) parseEnv() {

	if host, ok := os.LookupEnv(hostEnvKey); ok {
		if err := checkServerAddres(host); err == nil {
			c.host = host
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", hostEnvKey, host)
		}
	}

	//	get baseURL from environment variables
	if baseURL, ok := os.LookupEnv(baseURLEnvKey); ok {
		if err := checkBaseURL(baseURL); err == nil {
			c.baseURL = baseURL
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", baseURLEnvKey, baseURL)
		}
	}

	//	get storage file from environment variables
	if storageFile, ok := os.LookupEnv(storageFileEnvKey); ok {
		c.storageFile = storageFile
	}

}

// ------------------------------------------------------------
//
//	Check format IP addres and port.
//	Input:
//		addres string - for example 127.0.0.1:8080
//	Output:
//		err error
func checkServerAddres(addres string) error {

	args := strings.Split(addres, ":")
	if len(args) != 2 {
		return fmt.Errorf("ip addres and port in not right format: %s. for example: 192.168.1.2:port", addres)
	}

	if args[0] != "localhost" && net.ParseIP(args[0]) == nil {
		return fmt.Errorf("ip addres in not right format: %s. for example: 192.168.1.2:port", args[0])
	}

	if _, err := strconv.Atoi(args[1]); err != nil {
		return fmt.Errorf("port in not right format: %s. for example: addres:80", args[1])
	}
	return nil
}

// ------------------------------------------------------------
//
//	Check format base URL IP addres and port.
//	Input:
//		addres string - for example https://127.0.0.1:8080
//	Output:
//		err error
func checkBaseURL(baseURLFromOpt string) error {

	var strs = strings.SplitAfterN(baseURLFromOpt, "//", 2)
	if len(strs) != 2 {
		return fmt.Errorf("ip addres in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
	}

	if strs[0] != "http://" && strs[0] != "https://" {
		return fmt.Errorf("base url in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
	}

	args := strings.Split(strs[1], ":")
	if len(args) != 2 {
		return fmt.Errorf("ip addres and port in not right format: %s. for example: 192.168.1.2:port", strs[1])
	}

	if args[0] != "localhost" && net.ParseIP(args[0]) == nil {
		return fmt.Errorf("ip addres in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
	}

	if _, err := strconv.Atoi(args[1]); err != nil {
		return fmt.Errorf("ort in not right format: %s. for example: http://addres:80", baseURLFromOpt)
	}

	return nil
}
