package main

import (
	"os"
	"path/filepath"
	"testing"
)

func defaultConfig(t *testing.T) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Fatal("Environment GOPATH needs to be set for this test")
	}
	return filepath.Join(gopath, "src", "github.com", "midstar", "plm", "plm.config")
}

func TestConfig(t *testing.T) {
	config := LoadConfiguration(defaultConfig(t))
	assertEqualsInt(t, "config.Port", 12124, config.Port)
	assertEqualsInt(t, "config.FastLogTimeMs", 6000, config.FastLogTimeMs)
	assertEqualsInt(t, "config.SlowLogFactor", 10, config.SlowLogFactor)
	assertEqualsInt(t, "config.FastLogSize", 600, config.FastLogSize)
	assertEqualsInt(t, "config.SlowLogSize", 1440, config.SlowLogSize)
}

func TestConfigInvalidFile(t *testing.T) {
	config := LoadConfiguration("dont_exist.properties")
	assertEqualsInt(t, "config.Port", 12124, config.Port)
	assertEqualsInt(t, "config.FastLogTimeMs", 3000, config.FastLogTimeMs)
	assertEqualsInt(t, "config.SlowLogFactor", 20, config.SlowLogFactor)
	assertEqualsInt(t, "config.FastLogSize", 1200, config.FastLogSize)
	assertEqualsInt(t, "config.SlowLogSize", 1440, config.SlowLogSize)
}

func TestLoadPropertyInt(t *testing.T) {
	properties := make(map[string]string)
	properties["mkey"] = "invalid"
	value := getPropertyInt(properties, "mkey", 3)
	assertEqualsInt(t, "mKey value", 3, value)
}

func TestLoadProperties(t *testing.T) {
	properties, errprop := LoadPropertyFile(defaultConfig(t))
	if errprop != nil {
		t.Fatal(errprop)
	}
	assertEqualsInt(t, "Size of properties", 5, len(properties))
	assertEqualsStr(t, "Value of property port", "12124", properties["port"])
	assertEqualsStr(t, "Value of property fastLogTimeMs", "6000", properties["fastLogTimeMs"])
	assertEqualsStr(t, "Value of property slowLogFactor", "10", properties["slowLogFactor"])
	assertEqualsStr(t, "Value of property fastLogSize", "600", properties["fastLogSize"])
	assertEqualsStr(t, "Value of property slowLogSize", "1440", properties["slowLogSize"])
}

func TestLoadPropertiesInvalidFile(t *testing.T) {
	properties, errprop := LoadPropertyFile("dont_exist.properties")
	if errprop == nil {
		t.Fatal("Expected an error when loading properties")
	}
	assertEqualsInt(t, "Size of properties", 0, len(properties))
}
