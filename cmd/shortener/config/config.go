package config

import (
	"flag"
	"fmt"
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

	if err := checkNetAddres(*Host); err != nil{
		log.Panicf("Server IP addres and port in not right format: %s. For example: 192.168.1.2:80\n", *Host)
	}

	if err := checkNetAddres(*ResponseBaseAddres); err != nil{
		log.Panicf("Response base addres in not right format: %s. For example: 192.168.1.2:80\n", *ResponseBaseAddres)
	}

	log.Printf("Server host: %s\t Response base addres: %s\n", *Host, *ResponseBaseAddres)
}


// ------------------------------------------------------------
//	Check format IP addres and port.
//	Input:
//		addres string - for example 127.0.0.1:8080
//	Output:
//		err error
func checkNetAddres(addres string) error{
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