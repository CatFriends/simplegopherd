package gopher

import "log"
import "path"
import "path/filepath"
import "os"
import "fmt"
import "encoding/csv"
import "strings"
import "io/ioutil"

import "../configuration"

const IndexFileName = "index.csv"
const NewLineSequence = "\n"

const (
  EText      = "0"
  EDirectory = "1"
  EError     = "3"
  EInfo      = "i"
  EGIFImage  = "g"
  EImage     = "I"
  EBinary    = "9"
  ESound     = "s"
)

const empty = ""

func ProcessRequest(selector string) ([]byte) {
  log.Printf("Processing request [%s]", selector, configuration.BaseDirectory())
  if selector == empty {
    return HandleSelector(configuration.BaseDirectory())
  } else {
    return HandleSelector(selector)
  }
}

func ReadIndex(referenceDir string) ([]byte) {
  log.Printf("Reading index of [%s]", referenceDir)
  if f, e := os.Open(path.Join(referenceDir, IndexFileName)); e != nil {
    return []byte(gopherError(e.Error()))
  } else {
    defer f.Close()
    
    index := make([]string, 0)

    reader := csv.NewReader(f)
    if lines, e := reader.ReadAll(); e != nil {
      return []byte(gopherError(e.Error()))
    } else {
      for _, line := range(lines) {
        index = append(index, indexEntry(line[0], referenceDir, line[1]))
      }
    }

    return []byte(strings.Join(index, NewLineSequence) + NewLineSequence)

  }
}

func HandleSelector(selector string) ([]byte) {
  var dirPath string
  if selector != configuration.BaseDirectory() {
    dirPath = path.Join(configuration.BaseDirectory(), selector)
  } else {
    dirPath = selector
  }
  if stat, e := os.Stat(dirPath); e != nil {
    return []byte(gopherError(e.Error()))
  } else {
    if stat.IsDir() == true {
      return ReadIndex(dirPath)
    } else {
      if data, e := ioutil.ReadFile(path.Join(configuration.BaseDirectory(), selector)); e != nil {
        return []byte(gopherError(e.Error()))
      } else {
        return data[:];
      }
    }
  }
}

func gopherEntry(etype, title, url string) (string) {
  return fmt.Sprintf("%s%s\t%s\t%s\t%s", etype, title, url, configuration.HostName(), configuration.PortNumber())
}

func gopherError(reason string) (string) {
  return gopherEntry(EError, reason, "")
}

func indexEntry(title, referenceDir, selector string) (string) {
  if selector == empty {
    return gopherEntry(EInfo, title, empty)
  } else {
    if stat, e := os.Stat(path.Join(referenceDir, selector)); e != nil {
      log.Printf("Error: %s", e.Error())
      return gopherError(e.Error())
    } else {
      if stat.IsDir() == true {
        return gopherEntry(EDirectory, title, selector)
      } else {
        replacer := strings.NewReplacer(configuration.BaseDirectory(), empty)
        entrySelector := replacer.Replace(path.Join(referenceDir, selector))
        log.Printf("Extension of [%s] is [%s]", entrySelector, strings.ToLower(filepath.Ext(selector)))
        switch strings.ToLower(filepath.Ext(selector)) {
        case ".txt": return gopherEntry(EText, title, entrySelector)
        case ".gif": return gopherEntry(EGIFImage, title, entrySelector)
        case ".jpg": return gopherEntry(EImage, title, entrySelector)
        case ".mp3": return gopherEntry(ESound, title, entrySelector)
        default:     return gopherEntry(EBinary, title, entrySelector)
        }
      }
    }
  }
}