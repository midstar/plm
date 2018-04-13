package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

type program struct {
	plm        *PLM
	workingDir string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	p.plm = CreatePLM(p.workingDir)
	go p.run()
	return nil
}
func (p *program) run() {
	p.plm.Start()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.plm.Stop()
	return nil
}

// Main method can be runned as a "normal" console application AND as a
// service. See install.bat for how to install the service on Windows.
//
// The application takes one optional argument which is the working
// directory, i.e. where logs, configs and templates are found.
// If not specified the current directory is used.
func main() {
	workingDir := ""
	if len(os.Args) > 1 {
		workingDir = os.Args[1]
	}

	svcConfig := &service.Config{
		Name:        "plm",
		DisplayName: "Process Load Monitor",
		Description: "Process Load Monitor Service",
	}

	prg := &program{workingDir: workingDir}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
