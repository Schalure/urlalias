// Configuration appliation github.com/Schalure/urlalias
package config

import (
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
	AppName        = string("github.com/Schalure/urlalias") //	Application name
	hostDefault    = string("localhost:8080")               //	Host default value
	baseURLDefault = string("http://localhost:8080")             //	Base URL default value
	hostEnvKey     = string("SERVER_ADDRESS")               //	key for "host" in environment variables
	baseURLEnvKey  = string("BASE_URL")                     //	key for "baseURL" in environment variables
)

// ------------------------------------------------------------
//
//	Configuration vars for application
var (
	Config Configuration
)

// ------------------------------------------------------------
//
//	Initialize config options.
func Initialize() {

	Config.initialize()

	//	get host from environment variables
	if hostFromEnv, ok := os.LookupEnv(hostEnvKey); ok {
		if err := checkServerAddres(hostEnvKey); err == nil {
			Config.host = hostFromEnv
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", hostEnvKey, hostFromEnv)
		}
	}

	//	get baseURL from environment variables
	if baseURLFromEnv, ok := os.LookupEnv(baseURLEnvKey); ok {
		if err := checkBaseURL(baseURLFromEnv); err == nil {
			Config.baseURL = baseURLFromEnv
		} else {
			log.Printf("The environment variable \"%s\" is written in the wrong format: %s", baseURLEnvKey, baseURLFromEnv)
		}
	}

	log.Printf("Server address: \"%s\", Base URL: \"%s\"", Config.host, Config.baseURL)
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

	if args[0] != "localhost" {
		if net.ParseIP(args[0]) == nil {
			return fmt.Errorf("ip addres in not right format: %s. for example: 192.168.1.2:port", args[0])
		}
	}

	if _, err := strconv.Atoi(args[1]); err != nil {
		return fmt.Errorf("ort in not right format: %s. for example: addres:80", args[1])
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
	} else {
		if strs[0] != "http://" || strs[0] != "https://"{
			return fmt.Errorf("Base URL in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
		}
		args := strings.Split(strs[1], ":")

		if args[0] != "localhost" {
			if net.ParseIP(args[0]) == nil {
				return fmt.Errorf("ip addres in not right format: %s. for example: http://192.168.1.2:port", baseURLFromOpt)
			}
		}

		if _, err := strconv.Atoi(args[1]); err != nil {
			return fmt.Errorf("ort in not right format: %s. for example: http://addres:80", baseURLFromOpt)
		}
	}
	return nil
}

// type Options struct {
// 	Host    string `env:"SERVER_ADDRESS"` //	Server addres
// 	BaseURL string `env:"BASE_URL"`       //	Base URL for create alias
// }

// func InitEnvOpt() *Options {

// 	var opt Options

// 	//	Get env vars
// 	if err := env.Parse(&opt); err != nil {
// 		log.Fatal(err)
// 	}
// 	return &opt
// }

// Host: *flag.String("a", "localhost:8080", "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80"),
// BaseURL: *flag.String("b", "http://localhost:8080", "Response base addres for alias URL.\n\tFor example: 192.168.1.2"),
