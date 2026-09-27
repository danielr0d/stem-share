package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

func receive(ctx context.Context, code, outDir string, force bool) error {
	var c wormhole.Client
	msg, err := c.Receive(ctx, code)
	if err != nil {
		return err
	}
	if msg.Type != wormhole.TransferDirectory {
		msg.Reject()
		return fmt.Errorf("expected a set of stems, got %s", msg.Type)
	}

	tmp, err := os.CreateTemp("", "stemshare-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	fmt.Printf("Receiving %d file(s)\n", msg.FileCount)
	bar := pb.Full.Start64(int64(msg.TransferBytes))
	bar.Set(pb.Bytes, true)
	_, err = io.Copy(tmp, bar.NewProxyReader(msg))
	bar.Finish()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	return extract(tmp.Name(), outDir, force)
}

func extract(zipPath, outDir string, force bool) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := filepath.Base(f.Name)
		if !isWAVName(name) {
			fmt.Printf("skipped %s (not a .wav)\n", name)
			continue
		}
		dst := filepath.Join(outDir, name)
		if err := extractFile(f, dst, force); err != nil {
			fmt.Printf("skipped %s: %v\n", name, err)
			continue
		}
		fmt.Println("saved", dst)
	}
	return nil
}

func extractFile(f *zip.File, dst string, force bool) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	var hdr bytes.Buffer
	if err := checkHeader(io.TeeReader(r, &hdr)); err != nil {
		return err
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if force {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	out, err := os.OpenFile(dst, flags, 0o644)
	if os.IsExist(err) {
		return fmt.Errorf("already exists (use -force)")
	}
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, io.MultiReader(&hdr, r)); err != nil {
		out.Close()
		os.Remove(dst)
		return err
	}
	return out.Close()
}
