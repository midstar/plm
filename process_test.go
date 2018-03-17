// process unit tests
package main

import (
	"testing"
)

func TestProcessUpdate(t *testing.T) {
  processMap := NewProcessMap()
  processMap.Update()
}