package configuration

import "code.google.com/p/gcfg"

type configurationDataType struct {
  Network struct {
    HostName string
    PortNumber string
  }
  Site struct {
    BaseDirectory string
    IndexFileName string
  }
  Gopher struct {
    NewLineSequence string
  }
}

var configurationData configurationDataType

func LoadFromFile(name string) (error) {
  return gcfg.ReadFileInto(&configurationData, name)
}

func HostName() (string)        { return configurationData.Network.HostName }
func PortNumber() (string)      { return configurationData.Network.PortNumber }
func Binding() (string)         { return fmt.Sprintf("%s:%s", configurationData.Network.HostName, configurationData.Network.PortNumber) }
func BaseDirectory() (string)   { return configurationData.Site.BaseDirectory }
func IndexFileName() (string)   { return configurationData.Site.IndexFileName }
func NewLineSequence() (string) { return configurationData.Gopher.NewLineSequence }