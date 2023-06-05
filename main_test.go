package main

import (
	"testing"
)

func TestCheckFileStatus(t *testing.T) {

	err := CheckFileStatus("fate.w3x", "application/octet-stream")

	if err != nil {
		t.Error(err)
	}

}

func TestLoadMapsFolder(t *testing.T) {

	LoadMapsFolder()
}
