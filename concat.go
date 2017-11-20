package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func ConcatAACFile(ctx context.Context, aacDir, filename string) error {
	var err error
	aacCopyDir, err = TempDiraac()
	if err != nil {
		log.Fatalf("temDir err: %v\n", err)
	}
	defer os.RemoveAll(aacCopyDir)

	err = concat(ctx, aacDir)
	if err != nil {
		return err
	}
	err = ResultAAC(ctx, filename)
	if err != nil {
		return err
	}
	return nil
}

//concat aac files
func concat(ctx context.Context, aacDir string) error {
	files, err := ioutil.ReadDir(aacDir)
	if err != nil {
		return err
	}
	var res []string
	for _, f := range files {
		res = append(res, f.Name())
	}

	ConcatAAC(ctx, res, aacDir)
	return nil
}

// concat AAC file using ffmpeg
func ConcatAAC(ctx context.Context, files []string, aacDir string) {
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)

		go func(ctx context.Context, file, aacDir string, num int) {
			defer wg.Done()
			var fname string
			name := fmt.Sprintf("concat:%s", file)

			if num >= 0 && num < 10 {
				fname = fmt.Sprintf("0000%d.aac", num)
			} else if num >= 10 && num < 100 {
				fname = fmt.Sprintf("000%d.aac", num)
			} else if num >= 100 && num < 1000 {
				fname = fmt.Sprintf("00%d.aac", num)
			} else if num >= 1000 && num < 10000 {
				fname = fmt.Sprintf("0%d.aac", num)
			}
			fname = filepath.Join(aacCopyDir, fname)
			cmd := exec.CommandContext(ctx, "ffmpeg", "-i", name, "-c", "copy", fname)
			cmd.Dir = aacDir
			cmd.Run()

		}(ctx, file, aacDir, i)

		wg.Wait()
	}
}

func ResultAAC(ctx context.Context, filename string) error {
	files, err := ioutil.ReadDir(aacCopyDir)
	if err != nil {
		return err
	}

	var res []byte
	for _, f := range files {
		res = append(res, f.Name()...)
		res = append(res, '|')
	}
	name := fmt.Sprintf("concat:%s", string(res[:len(res)-1]))

	aacFile = filepath.Join(RadikoPath, "RadioOutput", filename)
	cmd := exec.CommandContext(ctx, cmdPath, "-i", name, "-c", "copy", aacFile)
	cmd.Dir = aacCopyDir
	cmd.Run()

	return nil
}

func TempDiraac() (string, error) {
	aacDir, err := ioutil.TempDir(RadikoPath, "aac")
	if err != nil {
		return "", err
	}
	return aacDir, nil
}
