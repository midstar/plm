// process unit tests
package main

import (
	"github.com/midstar/proci"
	"testing"
)

type ProcessInterface interface {
}

func TestProcessUpdate(t *testing.T) {
	processMap := NewProcessMap(proci.Proci{})
	processMap.Update()
}
