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

  if lsck, e := net.Listen("tcp", configuration.Binding()); e != nil {
    log.Fatal(fmt.Sprintf("Can't start server: %s", e.Error()))
  } else {
    log.Printf("Gopher is now up at %s", configuration.Binding())
    for {
      if conn, e := lsck.Accept(); e != nil {
        log.Printf("Can't serve connection: %s", e.Error())
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
    log.Printf("Can't get selector string: %s", e.Error())
  } else {
    sck.Write(gopher.ProcessRequest(strings.Trim(selector, "\n\r\t ")))
  }

}
