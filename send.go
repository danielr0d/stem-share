package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

func send(ctx context.Context, paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("no files given")
	}

	dir := "stems-" + time.Now().Format("20060102-150405")
	seen := map[string]bool{}
	var entries []wormhole.DirectoryEntry

	for _, p := range paths {
		if err := validateWAV(p); err != nil {
			return err
		}
		name := filepath.Base(p)
		if seen[name] {
			return fmt.Errorf("duplicate file name %q", name)
		}
		seen[name] = true
		path := p
		entries = append(entries, wormhole.DirectoryEntry{
			Path:   dir + "/" + name,
			Mode:   0o644,
			Reader: func() (io.ReadCloser, error) { return os.Open(path) },
		})
	}

	var bar *pb.ProgressBar
	progress := func(sent, total int64) {
		if bar == nil {
			bar = pb.Full.Start64(total)
			bar.Set(pb.Bytes, true)
		}
		bar.SetCurrent(sent)
	}

	var c wormhole.Client
	code, result, err := c.SendDirectory(ctx, dir, entries, wormhole.WithProgress(progress))
	if err != nil {
		return err
	}

	fmt.Printf("Sending %d file(s). On the other computer run:\n\n    stemshare receive %s\n\n", len(entries), code)

	select {
	case r := <-result:
		if bar != nil {
			bar.Finish()
		}
		if !r.OK {
			return r.Error
		}
		fmt.Println("Done.")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
