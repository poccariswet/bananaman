package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func ConcatAACFile(ctx context.Context, aacDir, filename string) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 1)
	)

	files, err := ioutil.ReadDir(aacDir)
	if err != nil {
		return err
	}

	tempdir, err := TempDiraac()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempdir)

	for _, f := range files {
		wg.Add(1)
		go func(fname, tempdir string) {
			defer wg.Done()

			buf, err := AACtoByte(fname, aacDir)
			if err != nil {
				errChan <- err
			}

			c := 0
			for i, _ := range buf {
				if fmt.Sprintf("%x", buf[i]) == "5c" && fmt.Sprintf("%x", buf[i+1]) == "ff" {
					c = i + 1
					break
				}
			}

			if err := createAAC(filepath.Join(tempdir, fname), buf[c:]); err != nil {
				errChan <- err
			}
		}(f.Name(), tempdir)
	}
	select {
	case err := <-errChan:
		return err
	default:
	}
	wg.Wait()

	if err := RemakeAAC(ctx, filename, tempdir); err != nil {
		return err
	}
	return nil
}

func createAAC(name string, bf []byte) error {
	wf, err := os.Create(name)
	if err != nil {
		return err
	}
	defer wf.Close()

	wf.Write(bf)

	return nil
}

func TempDiraac() (string, error) {
	aacDir, err := ioutil.TempDir(homepath, "output")
	if err != nil {
		return "", err
	}
	return aacDir, nil
}

func AACtoByte(fname, aacDir string) ([]byte, error) {
	fpath := filepath.Join(aacDir, fname)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	f, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, f.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func RemakeAAC(ctx context.Context, filename, fpath string) error {
	files, err := ioutil.ReadDir(fpath)
	if err != nil {
		return err
	}

	var res []byte
	for _, f := range files {
		res = append(res, f.Name()...)
		res = append(res, '|')
	}
	name := fmt.Sprintf("concat:%s", string(res[:len(res)-1]))

	aacFile := filepath.Join(homepath, "RadioOutput", filename)
	cmd := exec.CommandContext(ctx, cmdPath, "-i", name, "-c", "copy", fmt.Sprintf("%s.aac", aacFile))
	cmd.Dir = fpath
	cmd.Run()

	return nil
}
