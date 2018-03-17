package main

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
)

const configFile = "config.json"

type Configuration struct {
    Port    int `json:"port"`
}

func loadConfiguration() Configuration {
  raw, fileerr := ioutil.ReadFile(configFile)
  if fileerr != nil {
    panic(fmt.Sprintf("unable to load configuration from %s. Reason: %s", configFile, fileerr))
  }
  configuration := Configuration{}
  err := json.Unmarshal(raw, &configuration)
  if err != nil {
    panic(fmt.Sprintf("Unable to decode configuration from %s. Reason: %s", configFile, err))
  }
  return configuration
}

