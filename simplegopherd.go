package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "log"
  "net"
  "os"
  "path"
  "strings"
)

const ClError = "3"
const ClInfo = "i"
const ClDirectory = "1"
const ClBinary = "9"

const IndexFN = "index.cat"

var extensions = map[string]string{
  ".txt": "0",
  ".diz": "0",
  ".nfo": "0",
  ".gif": "g",
  ".jpg": "I",
  ".png": "I",
  ".mp3": "s",
  ".wav": "s",
  ".mid": "s",
}

var host *string = flag.String("host", "127.0.0.1", "Bind to specified address")
var port *string = flag.String("port", "70", "Bind to specified port")
var base *string = flag.String("base", ".", "Directory to serve files from")

func main() {
  flag.Parse()

  log.Printf("Starting Gopher at %s:%s using %s", *host, *port, *base)

  if listener, e := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port)); e != nil {
    log.Fatalf("Can't create listener: %s", e.Error())
  } else {
    log.Printf("Waiting for client connections...")
    for {
      if client, e := listener.Accept(); e != nil {
        log.Printf("Failed to accept client: %s", e.Error())
      } else {
        go gopher(client)
      }
    }
  }

}

func gopher(sck net.Conn) {
  defer sck.Close()

  if selector, e := bufio.NewReader(sck).ReadString('\n'); e != nil {
    log.Printf("Can't read selector string: %s", e.Error())
  } else {
    selector = strings.Trim(selector, "\n\r\t ")

    fso := path.Join(*base, selector)

    if stat, e := os.Stat(fso); e != nil {
      log.Printf("Can't get stats of %s: %s", fso, e.Error())
    } else {

      if stat.IsDir() == true {
        gopher_index(path.Join(fso, IndexFN), sck)
      } else {
        if strings.HasSuffix(fso, IndexFN) {
          gopher_index(fso, sck)
        } else {
          gopher_send(fso, sck)
        }
      }

    }

  }

}

func gopher_index(fileName string, sck net.Conn) {
  if indexFile, e := os.Open(fileName); e != nil {
    log.Printf("Can't read Index %s: %s", fileName, e.Error())
    gopher_error(e.Error(), sck)
  } else {
    defer indexFile.Close()

    scanner := bufio.NewScanner(indexFile)

    for scanner.Scan() {
      record := strings.Split(scanner.Text(), "\t")

      switch len(record) {
      case 0:
        gopher_entry(ClInfo, "", "", sck)
      case 1:
        gopher_entry(ClInfo, record[0], "", sck)
      case 2:
        gopher_entry(getClass(fileName, record[1]), record[0], record[1], sck)
      }

    }

  }
}

func getClass(base, selector string) string {
  var fileName = path.Join(path.Dir(base), selector)
  if stat, e := os.Stat(fileName); e != nil {
    log.Printf("Can't get stats of %s: %s", fileName, e.Error())
    return ClInfo
  } else {
    if stat.IsDir() {
      return ClDirectory
    } else {
      if class, e := extensions[strings.ToLower(path.Ext(fileName))]; e != true {
        return ClBinary
      } else {
        return class
      }
    }
  }
}

func gopher_generate() {}

func gopher_send(fileName string, sck net.Conn) {
  if file, e := os.Open(fileName); e != nil {
    log.Printf("Can't read file %s: %s", fileName, e.Error())
    gopher_error(e.Error(), sck)
  } else {
    if _, e := io.Copy(sck, file); e != nil {
      log.Printf("Failed to send file %s: %s", fileName, e.Error())
    }
  }
}

func gopher_error(message string, sck net.Conn) {
  gopher_entry(ClError, message, "", sck)
}

func gopher_entry(class, text, selector string, sck net.Conn) {
  sck.Write([]byte(fmt.Sprintf("%s%s\t%s\t%s\t%d\n", class, text, selector, host, port)))
}
