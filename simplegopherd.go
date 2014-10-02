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
  "path/filepath"
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

var host *string = flag.String("host", "127.0.0.1", "Host name or IP address to listen")
var port *string = flag.String("port", "70",        "Port number for incoming connections, usually 70")
var base *string = flag.String("base", ".",         "Directory to be published")

func main() {
  flag.Parse()
  *base, _ = filepath.Abs(*base)

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

  raw := make([]byte, 255)
  
  if _, e := sck.Read(raw); e != nil {
    log.Printf("Can't read selector: %s", e.Error())
  } else {
    selector := strings.Trim(string(raw), "\n\r\x00 ")
    request, _ := filepath.Abs(path.Join(*base, selector))

    if strings.Contains(request, "\t") == false {
      
      // Simple file request

      if strings.HasPrefix(request, *base) == false {  // Check if requested object is inside base directory
        // Access violation
      } else {
        // Ok, we can serve
        if stat, e := os.Stat(request); e != nil {
          log.Printf("Can't get attributes of %s: %s", request, e.Error())
        } else {
          if stat.IsDir() {
            // Serve Index
          } else {
            gopher_send(request, sck)
          }
        }
      }

    } else {
      // Search query
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
