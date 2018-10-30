package main

import (
	//pb "github.com/alexhunt7/gofigure/proto"
	//"golang.org/x/net/context"
	"os"
	"testing"
)

func TestParseFileMode(t *testing.T) {
	tests := []struct {
		smode string
		fmode os.FileMode
	}{
		{"0777", os.FileMode(uint(511))},
		{"777", os.FileMode(uint(511))},
		{"0644", os.FileMode(uint(420))},
		{"0000", os.FileMode(uint(0))},
		{"0002", os.FileMode(uint(2))},
		{"0400", os.FileMode(uint(256))},
		{"0600", os.FileMode(uint(384))},
	}

	for _, tt := range tests {
		omode, _ := parseFileMode(tt.smode)
		if omode != tt.fmode {
			t.Errorf("TestParseFileMode, %s: %v != %v", tt.smode, omode, tt.fmode)
		}
	}
}

func TestParseFileModeError(t *testing.T) {
	tests := []string{
		"asdf",
		"",
		"!",
		"10.4",
		"9",
	}

	for _, tt := range tests {
		_, err := parseFileMode(tt)
		if err == nil {
			t.Errorf("TestParseFileModeError did not error as expected")
		}
	}
}
