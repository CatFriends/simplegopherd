package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/catfriends/simplegopherd/configuration"
	"github.com/catfriends/simplegopherd/gopher"
)

var (
	configFile = flag.String("config", "", "Configuration file name")
)

func main() {

	// Load configuration from file specified

	flag.Parse()
	if *configFile == "" {
		log.Fatal("You must specify config file name using -config=<file.ini>")
	} else {
		if e := configuration.LoadFromFile(*configFile); e != nil {
			log.Fatal(fmt.Sprintf("Can't load configuration from %s: %s", *configFile, e.Error()))
		}
	}

	// Create server instance

	if lsck, e := net.Listen("tcp", configuration.Binding()); e != nil {
		log.Fatal(fmt.Sprintf("Can't start: %s", e.Error()))
	} else {
		for {
			if conn, e := lsck.Accept(); e != nil {
				log.Printf("Can't serve client: %s", e.Error())
			} else {
				go HandleRequest(conn)
			}
		}
	}

}

func HandleRequest(sck net.Conn) {
	defer sck.Close()

	if selector, e := bufio.NewReader(sck).ReadString('\n'); e != nil {
		log.Printf("Can't get selector string: %s", e.Error())
	} else {
		sck.Write(gopher.ProcessRequest(strings.Trim(selector, "\n\r\t ")))
	}

}
