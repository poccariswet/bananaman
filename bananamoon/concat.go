package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func ConcatAACFile(ctx context.Context, aacDir, filename string) error {
	files, err := ioutil.ReadDir(aacDir)
	if err != nil {
		return err
	}

	tfile, err := ioutil.TempFile(aacDir, "aac")
	if err != nil {
		return err
	}
	defer os.Remove(tfile.Name())

	for _, f := range files {
		path := fmt.Sprintf("file '%s'\n", filepath.Join(aacDir, f.Name()))
		if _, err := tfile.WriteString(path); err != nil {
			return err
		}
	}

	if err := concat(ctx, filename, aacDir, tfile.Name()); err != nil {
		return err
	}

	return nil
}

func concat(ctx context.Context, filename, fpath, tempfile string) error {
	aacFile := filepath.Join(homepath, "RadioOutput", filename)
	cmd := exec.CommandContext(
		ctx,
		cmdPath,
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		tempfile,
		"-c",
		"copy",
		aacFile,
	)
	cmd.Dir = fpath
	cmd.Run()

	return nil
}

func TempDiraac() (string, error) {
	aacDir, err := ioutil.TempDir(homepath, "output")
	if err != nil {
		return "", err
	}
	return aacDir, nil
}
