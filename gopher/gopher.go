package gopher

import "log"
import "path"
import "path/filepath"
import "os"
import "fmt"
import "encoding/csv"
import "strings"
import "io/ioutil"

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
func ProcessRequest(selector string) []byte {
	if selector == empty {
		return HandleSelector(empty)
	} else {
		return HandleSelector(selector)
	}
}

// Process selector given by request.
func HandleSelector(selector string) []byte {
	var fsObject string

	// Find name of the filesystem object to be processed
	if selector == empty {
		fsObject = configuration.BaseDirectory()
	} else {
		fsObject = path.Join(configuration.BaseDirectory(), selector)
	}

  log.Printf("FSO for selector [%s] is [%s]", selector, fsObject)

	// Get information about the filesystem object
	if stat, e := os.Stat(fsObject); e != nil {
		return []byte(gopherError(e.Error()))
	} else {
		if stat.IsDir() == true {
			return ReadIndex(fsObject)
		} else {
			return ReadFile(fsObject)
		}
	}
}

// Find and process Index file for specified directory.
func ReadIndex(referenceDir string) ([]byte) {
  var indexFile = path.Join(referenceDir, configuration.IndexFileName())

  log.Printf("Processing Index [%s]", indexFile)

	if f, e := os.Open(indexFile); e != nil {
		return []byte(gopherError(e.Error()))
	} else {
		defer f.Close()

		index := make([]string, 0)

		reader := csv.NewReader(f)
		if lines, e := reader.ReadAll(); e != nil {
			return []byte(gopherError(e.Error()))
		} else {
			for _, line := range lines {
        entry := indexEntry(line[0], referenceDir, line[1])
				index = append(index, entry)

        log.Printf("  - %s", entry)

			}
		}

		return []byte(strings.Join(index, nl) + nl)

	}
}

// Read file in binary mode.
func ReadFile(name string) ([]byte) {
  log.Printf("Reading file [%s] in binary mode", name)

	if data, e := ioutil.ReadFile(name); e != nil {
		return []byte(gopherError(e.Error()))
	} else {
		return data[:]
	}
}

func gopherEntry(etype, title, url string) string {
	return fmt.Sprintf("%s%s\t%s\t%s\t%s", etype, title, url, configuration.HostName(), configuration.PortNumber())
}

func gopherError(reason string) string {
  log.Printf("Error: %s", reason)
	return gopherEntry(EError, reason, "")
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
