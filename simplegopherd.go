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

func main() {
  configurationFileName := flag.String("config", "", "INI configuration file name")
  flag.Parse()

  if e := configuration.LoadFromFile(*configurationFileName); e != nil {
    log.Fatal(e)
  }

  if lsck, e := net.Listen("tcp", fmt.Sprintf("%s:%s", configuration.HostName(), configuration.PortNumber())); e != nil {
    log.Fatal(fmt.Sprintf("Can't start server: %s", e.Error()))
  } else {
    log.Printf("Gopher is now listening to %s port %s", configuration.HostName(), configuration.PortNumber())
    for {
      if conn, e := lsck.Accept(); e != nil {
        log.Printf("Unable to serve incomming connection!")
      } else {
        log.Printf("Accepting client %s", conn.RemoteAddr().String())
        go HandleRequest(conn)
      }
    }
  }

}

func HandleRequest(sck net.Conn) {
  defer sck.Close()

  if selector, e := bufio.NewReader(sck).ReadString('\n'); e != nil {
    log.Printf("Bad selector %s", e.Error())
  } else {
    selector = strings.Trim(selector, "\n\r\t ")
    sck.Write(gopher.ProcessRequest(selector))
  }

}
