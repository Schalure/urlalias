package config

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	AppName = string("github.com/Schalure/urlalias")
	//Host    = string("localhost:8080")
)

// сразу используем глобальные переменные
var (
    Host = flag.String("a", "localhost:8080", "Server IP addres and port for server starting.\n\tFor example: 192.168.1.2:80")
    ResponseBaseAddres = flag.String("b", "localhost:8080", "Response base addres for alias URL.\n\tFor example: 192.168.1.2")
)


// ------------------------------------------------------------
//	Initialize config optipns.
//	It must finish without errors or panic will ensue.
func MustInit(){

	flag.Parse()


	args := strings.Split(*Host, ":")
	if len(args) != 2{
		log.Panicf("Server IP addres and port in not right format: %s. For example: 192.168.1.2:80\n", *Host)
	}


	if(args[0] != "localhost"){
		if net.ParseIP(args[0]) == nil{
			log.Panicf("Server IP addres in not right format: %s. For example: 192.168.1.2:port\n", args[0])
		}

	}

	if _, err := strconv.Atoi(args[1]); err != nil{
		log.Panicf("Server port in not right format: %s. For example: addresIP:80\n", args[1])
	}


	if(*ResponseBaseAddres != "localhost"){
		if net.ParseIP(*ResponseBaseAddres) == nil{
			log.Panicf("Response base addres in not right format: %s. For example: 192.168.1.2\n", *ResponseBaseAddres)
		}
	}

	log.Printf("Server host: %s\t Response base addres: %s\n", *Host, *ResponseBaseAddres)
}