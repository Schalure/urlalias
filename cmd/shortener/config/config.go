// Package config describes the configuration for github.com/Schalure/urlalias
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// Application constants
const (
	AppName            = string("github.com/Schalure/urlalias") //	Application name
	hostEnvKey         = string("SERVER_ADDRESS")               //	key for "host" in environment variables
	baseURLEnvKey      = string("BASE_URL")                     //	key for "baseURL" in environment variables
	storageFileEnvKey  = string("FILE_STORAGE_PATH")            //	key for "storageFile" in environment variables
	dbConnectionEnvKey = string("DATABASE_DSN")                 //	key for "dbConnection in environment variables
	enableHTTPSKey     = string("ENABLE_HTTPS")                 //	key for "EnableHTTPS in environment variables
	configFilePathKey  = string("CONFIG")                       //	key for "configFilePath in environment variables
)

// StorageType - enumeration type for Storage
type StorageType int

// Storage enumeration
const (
	MemoryStor StorageType = iota
	FileStor
	DataBaseStor
)

// String returns the name of the Storage type
func (s StorageType) String() string {
	return [...]string{"MemoryStor", "FileStor", "DataBaseStor"}[s]
}

// Default values
const (
	hostDefault           = string("localhost:8080")        //	Host default value
	baseURLDefault        = string("http://localhost:8080") //	Base URL default value
	aliasesFileDefault    = "/tmp/short-url-db.json"        //	Default file name of URLs storage
	usersFileDefault      = "/tmp/users-db.json"
	logToFileDefault      = false //	How to save log default value
	enableHTTPSDefault    = false // Default value for enableHTTPS
	configFilePathDefault = ""    // Default value for configFilePath
)

// Struct of configuration vars
type Configuration struct {
	host        string //	Server addres
	baseURL     string //	Base URL for create alias
	enableHTTPS bool   //	Flag for enable HTTPS

	aliasesFile  string // File name of URLs storage
	usersFile    string
	dbConnection string

	storageType StorageType

	logToFile bool //	true - save log to file, false - print log to console
}

// ConfigurationData for different configuration sources
type ConfigurationData struct {
	host         string
	baseURL      string
	aliasesFile  string
	usersFile    string
	dbConnection string
	enableHTTPS  bool
	configFile   string
}

// Common config variable
var config *Configuration

// NewConfig - constructor of Config type
func NewConfig() *Configuration {

	if config != nil {
		return config
	}

	defaultConfig := ConfigurationData{
		host:         hostDefault,
		baseURL:      baseURLDefault,
		aliasesFile:  aliasesFileDefault,
		usersFile:    usersFileDefault,
		dbConnection: "",
		enableHTTPS:  false,
		configFile:   "",
	}

	flagConfig := parseFlags()
	envConfig := parseEnv()
	fileConfig := ConfigurationData{}

	if envConfig.configFile != "" {
		fileConfig = parseConfigFile(envConfig.configFile)
	} else if flagConfig.configFile != "" {
		fileConfig = parseConfigFile(flagConfig.configFile)
	}

	config.setConfiguration(envConfig, flagConfig, fileConfig, defaultConfig)

	config.chooseStorageType()

	log.Printf("Server address: \"%s\"\n", config.host)
	log.Printf("Base URL: \"%s\"\n", config.host)

	switch config.storageType {
	case DataBaseStor:
		log.Printf("DB conection string: \"%s\"\n", config.dbConnection)
	case FileStor:
		log.Printf("Storage file: \"%s\"\n", config.aliasesFile)
	default:
		log.Print("memory storage is used")
	}

	log.Printf("Save log to file: \"%t\"\n", config.logToFile)
	return config
}

// Host - getter "Configuration.host"
func (c *Configuration) Host() string {
	return c.host
}

// BaseURL - getter "Configuration.baseURL"
func (c *Configuration) BaseURL() string {
	return c.baseURL
}

// AliasesFile - getter "Configuration.AliasesFile"
func (c *Configuration) AliasesFile() string {
	return c.aliasesFile
}

// UsersFile - getter "Configuration.UsersFile"
func (c *Configuration) UsersFile() string {
	return c.usersFile
}

// DBConnection - getter "Configuration.DBConnection"
func (c *Configuration) DBConnection() string {
	return c.dbConnection
}

// StorageType - getter "Configuration.StorageType"
func (c *Configuration) StorageType() StorageType {
	return c.storageType
}

// LogToFile - getter "Configuration.logSaver"
func (c *Configuration) LogToFile() bool {
	return c.logToFile
}

// Getter for enableHTTPS
func (c *Configuration) EnableHTTPS() bool {
	return c.enableHTTPS
}

// setConfiguration sets configurations from different sources. The configurations are installed in priority order.
// The highest priority is the 0th item in the `configs` list. If any fields are left empty, they are filled
// with default values from `defaultConfig`
func (c *Configuration) setConfiguration(defaultConfig ConfigurationData, configs ...ConfigurationData) {

	for _, config := range configs {
		if config.host != "" {
			c.host = config.host
			break
		}
	}
	if c.host == "" {
		c.host = defaultConfig.host
	}

	for _, config := range configs {
		if config.baseURL != "" {
			c.baseURL = config.baseURL
			break
		}
	}
	if c.baseURL == "" {
		c.baseURL = defaultConfig.baseURL
	}

	for _, config := range configs {
		if config.aliasesFile != "" {
			c.aliasesFile = config.aliasesFile
			break
		}
	}
	if c.aliasesFile == "" {
		c.aliasesFile = defaultConfig.aliasesFile
	}

	for _, config := range configs {
		if config.usersFile != "" {
			c.usersFile = config.usersFile
			break
		}
	}
	if c.usersFile == "" {
		c.usersFile = defaultConfig.usersFile
	}

	for _, config := range configs {
		if config.dbConnection != "" {
			c.dbConnection = config.dbConnection
			break
		}
	}
	if c.dbConnection == "" {
		c.dbConnection = defaultConfig.dbConnection
	}

	for _, config := range configs {
		if config.enableHTTPS != true {
			c.enableHTTPS = config.enableHTTPS
			break
		}
	}
	if c.enableHTTPS == true {
		c.enableHTTPS = defaultConfig.enableHTTPS
	}
}

// chooseStorageType choose storage type
func (c *Configuration) chooseStorageType() {

	if c.dbConnection != "" {
		c.storageType = DataBaseStor
	} else if c.aliasesFile != "" {
		c.storageType = FileStor
	} else {
		c.storageType = MemoryStor
	}
}

// parseFlags parses flags method of "Config" type
func parseFlags() ConfigurationData {

	configData := ConfigurationData{}

	configData.host = *flag.String("a", hostDefault, "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80")
	configData.baseURL = *flag.String("b", baseURLDefault, "Response base addres for alias URL.\n\tFor example: http://192.168.1.2")
	configData.dbConnection = *flag.String("d", "", "data base connection string")
	configData.enableHTTPS = *flag.Bool("s", enableHTTPSDefault, "Variant HTTP connect: true - HTTPS, false - HTTP")
	configData.configFile = *flag.String("c", configFilePathDefault, "Configuration file path")

	storageFile := ""
	flag.Func("f", "File name of URLs storage. Specify the full name of the file", func(s string) error {

		//	TODO need to parce file path to check valid
		if s == "" {
			storageFile = aliasesFileDefault
		} else {
			storageFile = s
		}
		return nil
	})

	flag.Parse()

	if err := checkServerAddres(configData.host); err != nil {
		configData.host = hostDefault
	}

	if err := checkBaseURL(configData.baseURL); err != nil {
		configData.baseURL = baseURLDefault
	}

	configData.aliasesFile = storageFile
	configData.usersFile = storageFile + "-users"

	return configData
}

// parseEnv parses environment variables method of "Config" type
func parseEnv() ConfigurationData {

	configData := ConfigurationData{}

	if host, ok := os.LookupEnv(hostEnvKey); ok {
		if err := checkServerAddres(host); err == nil {
			configData.host = host
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", hostEnvKey, host)
		}
	}

	//	get baseURL from environment variables
	if baseURL, ok := os.LookupEnv(baseURLEnvKey); ok {
		if err := checkBaseURL(baseURL); err == nil {
			configData.baseURL = baseURL
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", baseURLEnvKey, baseURL)
		}
	}

	//	get storage file from environment variables
	if storageFile, ok := os.LookupEnv(storageFileEnvKey); ok {
		configData.aliasesFile = storageFile
		configData.usersFile = storageFile + "-users"
	}

	//	get storage file from environment variables
	if dbConnection, ok := os.LookupEnv(dbConnectionEnvKey); ok {
		configData.dbConnection = dbConnection
	}

	if enableHTTPS, ok := os.LookupEnv(enableHTTPSKey); ok {
		if isEnableHTTPS, err := strconv.ParseBool(enableHTTPS); err != nil {
			configData.enableHTTPS = isEnableHTTPS
		}
	}

	if configFile, ok := os.LookupEnv(configFilePathKey); ok {
		configData.configFile = configFile
	}

	return configData
}

// setConfigurationFromFile reads configuration file and set configuration parametrs
func parseConfigFile(filePath string) ConfigurationData {

	type ConfigurationFileData struct {
		Host         string `json:"server_address"`
		BaseURL      string `json:"base_url"`
		DBFile       string `json:"file_storage_path"`
		ConnectionDB string `json:"database_dsn"`
		EnableHTTPS  bool   `json:"enable_https"`
	}
	var configurationFileData ConfigurationFileData
	configData := ConfigurationData{}

	if filePath == "" {
		return ConfigurationData{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ConfigurationData{}
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.Size() == 0 {
		return ConfigurationData{}
	}

	data := make([]byte, fileInfo.Size())

	if _, err = file.Read(data); err != nil {
		return ConfigurationData{}
	}

	if err := json.Unmarshal(data, &configurationFileData); err != nil {
		return ConfigurationData{}
	}

	if err := checkServerAddres(configurationFileData.Host); err == nil {
		configData.host = configurationFileData.Host
	}
	if err := checkBaseURL(configurationFileData.BaseURL); err == nil {
		configData.baseURL = configurationFileData.BaseURL
	}
	if configurationFileData.DBFile != "" {
		configData.aliasesFile = configurationFileData.DBFile
		configData.usersFile = configurationFileData.DBFile + "-users"
	}
	if configurationFileData.ConnectionDB != "" {
		configData.dbConnection = configurationFileData.ConnectionDB
	}
	configData.enableHTTPS = configurationFileData.EnableHTTPS

	return configData
}

// checkServerAddres checks format IP addres and port.
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

// checkBaseURL checks format base URL IP addres and port.
func checkBaseURL(baseURLFromOpt string) error {

	var strs = strings.SplitAfterN(baseURLFromOpt, "//", 2)
	if len(strs) != 2 {
		return fmt.Errorf("ip addres in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
	}

	if strs[0] != "http://" && strs[0] != "https://" {
		return fmt.Errorf("base url in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
	}

	return nil
}
