package configuration

import "code.google.com/p/gcfg"

type configurationDataType struct {
  Network struct {
    HostName string
    PortNumber string
  }
  Site struct {
    BaseDirectory string
  }
}

var configurationData configurationDataType

func LoadFromFile(name string) (error) {
  return gcfg.ReadFileInto(&configurationData, name)
}

func BaseDirectory() (string) { return configurationData.Site.BaseDirectory }
func HostName() (string)      { return configurationData.Network.HostName }
func PortNumber() (string)    { return configurationData.Network.PortNumber }