package gopher

import "log"
import "path"
import "path/filepath"
import "os"
import "fmt"
import "strings"
import "io/ioutil"
import "bytes"
import "bufio"

import "github.com/catfriends/simplegopherd/configuration"

const (
	EDirectory = "1"
	EError     = "3"
	EInfo      = "i"
	EBinary    = "9"
)

var extensionType = map[string]string{
	".txt":  "0",
	".diz":  "0",
	".nfo":  "0",
	".gif":  "g",
	".jpg":  "I",
	".jpeg": "I",
	".png":  "I",
	".mp3":  "s",
	".wav":  "s",
	".mid":  "s",
}

const empty = ""
const nl = "\n"

// Process Gopher request from the client.
func ProcessRequest(sck net.Conn) {
	if selector, e := bufio.NewReader(sck).ReadString('\n'); e != nil {
		gopherError(e.Error(), sck)
	} else {
		HandleSelector(getSelector(sck), sck)
	}
}

// Process selector given by request.
func HandleSelector(selector string, sck net.Conn) {
	var fsObject = getFSObjectName(selector)

	// Get information about the filesystem object
	if info, e := os.Stat(fsObject); e != nil {
		gopherError(e.Error(), sck)
	} else {
		if info.IsDir() == true {
			return ReadIndex(fsObject)
		} else {
			ReadFile(fsObject, sck)
		}
	}
}

// Find name of the filesystem object to be processed.
func getFSObjectName(selector string) (string)
	var fsObject string
	if selector == empty {
		fsObject = configuration.BaseDirectory()
	} else {
		fsObject = path.Join(configuration.BaseDirectory(), selector)
	}

  log.Printf("FSO for selector [%s] is [%s]", selector, fsObject)
  return fsObject
}

// Find and process Index file for specified directory.
func ReadIndex(referenceDir string) ([]byte) {
  var indexFile = path.Join(referenceDir, configuration.IndexFileName())

  log.Printf("Processing Index [%s]", indexFile)

	if f, e := os.Open(indexFile); e != nil {
		return []byte(gopherError(e.Error()))
	} else {
		defer f.Close()

		var index bytes.Buffer

    scanner := bufio.NewScanner(f);
    for scanner.Scan() {
    	title, selector := processIndexLine(scanner.Text())
      entry := indexEntry(title, referenceDir, selector)
      log.Printf("  - %s", entry)
      index.WriteString(entry + nl)
    }

    index.WriteString(nl)
		return index.Bytes()

	}
}

// Read file in binary mode.
func ReadFile(name string, sck net.Conn) {
  log.Printf("Reading file [%s] in binary mode", name)
	if data, e := ioutil.ReadFile(name); e != nil {
		gopherError(e.Error(), sck)
	} else {
		sck.Write(data[:])
	}
}

func gopherEntry(etype, title, url string) string {
	return fmt.Sprintf("%s%s\t%s\t%s\t%s", etype, title, url, configuration.HostName(), configuration.PortNumber())
}

func gopherError(reason string, socket net.Conn) {
  log.Printf("Error: %s", reason)
	socket.Write([]byte(gopherEntry(EError, reason, empty)))
}

func indexEntry(title, referenceDir, selector string) string {
	if selector == empty {
		return gopherEntry(EInfo, title, empty)
	} else {
		if stat, e := os.Stat(path.Join(referenceDir, selector)); e != nil {
			return gopherError(e.Error())
		} else {
			if stat.IsDir() == true {
				return gopherEntry(EDirectory, title, selector)
			} else {
				replacer := strings.NewReplacer(configuration.BaseDirectory(), empty)
				entrySelector := replacer.Replace(path.Join(referenceDir, selector))
				if extensionType, e := extensionType[strings.ToLower(filepath.Ext(selector))]; e != true {
					return gopherEntry(EBinary, title, entrySelector)
				} else {
					return gopherEntry(extensionType, title, entrySelector)
				}
			}
		}
	}
}

// Process single line from index file
func processIndexLine(line string) (string, string) {
  parts := strings.Split(line, "\t")
  switch len(parts) {
  case 0:
  	return empty, empty
  case 1:
  	return parts[0], empty
  default:
  	return parts[0], parts[1]
  }
}