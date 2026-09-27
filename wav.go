package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func isWAVName(name string) bool {
	return strings.EqualFold(filepath.Ext(name), ".wav")
}

func checkHeader(r io.Reader) error {
	var hdr [12]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return fmt.Errorf("not a WAV file (too short)")
	}
	if !bytes.Equal(hdr[0:4], []byte("RIFF")) || !bytes.Equal(hdr[8:12], []byte("WAVE")) {
		return fmt.Errorf("not a WAV file (missing RIFF/WAVE header)")
	}
	return nil
}

func validateWAV(path string) error {
	if !isWAVName(path) {
		return fmt.Errorf("%s: not a .wav file", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := checkHeader(f); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}
