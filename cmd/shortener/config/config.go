// Configuration appliation github.com/Schalure/urlalias
package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
)

// ------------------------------------------------------------
//	Struct of configuration vars
type Configuration struct{
	//	Server addres
	//	It can be use from environment variables or from app flags
	//		IMPORTANT: VALUE FROM ENVIRONMENT VARIABLES HAS FIRST PRIORITY
	Host string `env:"SERVER_ADDRESS"`
	//	Base URL for create alias
	//	It can be use from environment variables or from app flags
	//		IMPORTANT: VALUE FROM ENVIRONMENT VARIABLES HAS FIRST PRIORITY
	BaseURL string `env:"BASE_URL"`
}


// ------------------------------------------------------------
//	Application constants
const (
	//	Application name
	AppName = string("github.com/Schalure/urlalias")
)


// ------------------------------------------------------------
//	Configuration vars for application
var (
	//	Application configuration struct
	Config Configuration
)


// ------------------------------------------------------------
//	Initialize config optipns.
//	It must finish without errors or panic will ensue.
func MustInit(){

    hostFromOpt := flag.String("a", "localhost:8080", "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80")
    baseURLFromOpt := flag.String("b", "http://localhost:8080", "Response base addres for alias URL.\n\tFor example: 192.168.1.2")
	flag.Parse()

	//	Check host, if error set default value
	if err := checkServerAddres(*hostFromOpt); err != nil{
		log.Printf("ERROR: Server IP addres and port in not right format: %s. For example: 192.168.1.2:80\n", *hostFromOpt)
		*hostFromOpt = "localhost:8080"
	}

	//	Check baseURL, if error set default value
	if err := checkBaseURL(*baseURLFromOpt); err != nil{
		log.Printf("ERROR: Response base addres in not right format: %s. For example: 192.168.1.2:80\n", *baseURLFromOpt)
		*baseURLFromOpt = "http://localhost:8080"
	}

	//	Get env vars
	if err := env.Parse(&Config); err != nil{
		log.Fatal(err)
	}

	//	Set Host
	if Config.Host != ""{
		if err := checkServerAddres(Config.Host); err != nil{
			Config.Host = *hostFromOpt
			log.Println(err.Error())
		}
	}else{
		Config.Host = *hostFromOpt
	}

	//	Set BaseURL
	if Config.BaseURL != ""{
		if err := checkBaseURL(Config.BaseURL); err != nil{
			Config.BaseURL = *baseURLFromOpt
			log.Println(err.Error())
		}
	}else{
		Config.BaseURL = *baseURLFromOpt
	}

	log.Printf("Server addres: \"%s\", Base URL: \"%s\"", Config.Host, Config.BaseURL)
}


// ------------------------------------------------------------
//	Check format IP addres and port.
//	Input:
//		addres string - for example 127.0.0.1:8080
//	Output:
//		err error
func checkServerAddres(addres string) error{
	args := strings.Split(addres, ":")

	if(args[0] != "localhost"){
		if net.ParseIP(args[0]) == nil{
			return fmt.Errorf("ip addres in not right format: %s. for example: 192.168.1.2:port", args[0])
		}
	}

	if _, err := strconv.Atoi(args[1]); err != nil{
		return fmt.Errorf("ort in not right format: %s. for example: addres:80", args[1])
	}
	return nil
}

// ------------------------------------------------------------
//	Check format base URL IP addres and port.
//	Input:
//		addres string - for example https://127.0.0.1:8080
//	Output:
//		err error
func checkBaseURL(baseURLFromOpt string) error{

	var strs = strings.SplitAfterN(baseURLFromOpt, "//", 2)
	if len(strs) != 2{
		return fmt.Errorf("ip addres in not right format: %s. for example: 192.168.1.2:port", baseURLFromOpt)
	}else{
		args := strings.Split(strs[1], ":")

		if(args[0] != "localhost"){
			if net.ParseIP(args[0]) == nil{
				return fmt.Errorf("ip addres in not right format: %s. for example: 192.168.1.2:port", args[0])
			}
		}

		if _, err := strconv.Atoi(args[1]); err != nil{
			return fmt.Errorf("ort in not right format: %s. for example: addres:80", args[1])
		}		
	}
	return nil
}