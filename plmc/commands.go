package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// CmdPlot get plot for one or more processes
func CmdPlot(filename string) error {
	resp, err := http.Get(fmt.Sprintf("%s/plot", PLMUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	fmt.Print(filename, " written")
	if err != nil {
		return err
	}
	return nil
}
