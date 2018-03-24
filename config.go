package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// DefaultConfigFile default configuration file
const DefaultConfigFile = "plm.config"

// Configuration holds parameters that are configurable.
type Configuration struct {
	Port          int
	FastLogTimeMs int
	SlowLogFactor int
	FastLogSize   int
	SlowLogSize   int
}

// LoadConfiguration loads configuration from file and returns a
// Configuration. Default values will be used if configuration
// cannot be found.
func LoadConfiguration(fileName string) *Configuration {
	p, errprop := LoadPropertyFile(fileName)
	if errprop != nil {
		fmt.Println(errprop)
	}
	// Create configuration with default values
	configuration := Configuration{
		Port:          getPropertyInt(p, "port", 9090),
		FastLogTimeMs: getPropertyInt(p, "fastLogTimeMs", 3000),
		SlowLogFactor: getPropertyInt(p, "slowLogFactor", 20),
		FastLogSize:   getPropertyInt(p, "fastLogSize", 1200),
		SlowLogSize:   getPropertyInt(p, "slowLogSize", 1440)}

	return &configuration
}

func getPropertyInt(properties map[string]string, key string, defaultValue int) int {
	value, hasKey := properties[key]
	if !hasKey {
		return defaultValue
	}
	intValue, valueerr := strconv.Atoi(value)
	if valueerr != nil {
		fmt.Printf("Property %s does not have a valid integer value. Using default %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}

// LoadPropertyFile loads property files of the same format as found in Java
// property files and returns a map of strings.
func LoadPropertyFile(fileName string) (map[string]string, error) {
	properties := make(map[string]string)
	b, fileerr := ioutil.ReadFile(fileName)
	if fileerr != nil {
		return properties, fmt.Errorf("unable to load properties from %s. Reason: %s", fileName, fileerr)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		tLine := strings.TrimSpace(line)
		if len(tLine) > 0 && !strings.HasPrefix(tLine, "#") {
			parts := strings.SplitN(tLine, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if len(key) > 0 {
					properties[key] = value
				}
			}
		}
	}
	return properties, nil
}
