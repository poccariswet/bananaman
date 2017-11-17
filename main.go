package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	timeLayout = "20060102150405"
	filetype   = "aac"
)

var (
	location *time.Location
	homepath = os.Getenv("HOME")
)

func init() {
	var err error
	location, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("time package LoadLoaction err : %v\n", err)
	}

	root := filepath.Join(homepath, "RadioOutput")
	_, err = os.Stat(root)
	if err != nil {
		os.Mkdir(root, 0777)
		fmt.Println("Made a Dir at" + root)
	}
}

func main() {
	var ID, start, radiname string
	flag.StringVar(&ID, "id", "", "id")
	flag.StringVar(&start, "s", "", "start")
	flag.StringVar(&radiname, "file", "", "file")
	flag.Parse()
	if ID == "" {
		fmt.Println("Please input stationID, like '-id=TBS'")
		os.Exit(1)
	}
	if start == "" {
		fmt.Printf("Please input start time you wanna listen to radio name,\nlike 2017/11/11/01:00 -> 20171111010000\n")
		os.Exit(1)
	}
	if radiname == "" {
		fmt.Println("Please input filename,\nlike '-file=bananamoonGOLD'")
	}
	go spinner(100 * time.Millisecond)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	//とりあえず
	areaID := "JP13"
	startTime, err := time.ParseInLocation(timeLayout, start, location)
	if err != nil {
		log.Fatalf("time pacakge ParseInLocation err : %v\n", err)
	}

	// client 設定
	client, err := GetClient(ctx, areaID)
	if err != nil {
		log.Fatalf("Failed to GetClient: %v\n", err)
	}

	// トークンの取得と設定
	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		log.Fatalf("client AuthorizeToken func err : %v\n", err)
	}

	// M3U8ファイルの取得
	uri, err := client.CreateM3U8Playlist(ctx, ID, startTime)
	if err != nil {
		log.Fatalf("Failed playlist.m3u8 err: %v\n", err)
	}

	// list(M3U8ファイルの中にある、aaclist)
	list, err := Getlist(uri)
	if err != nil {
		log.Fatalf("Faild M3U8 list err : %v\n", err)
	}

	//一時的にディレクトリを生成(listから取得したaacファイルを保存用)
	aacDir, err := TempDiraac()
	if err != nil {
		log.Fatalf("tempAACdir err : %v\n", err)
	}
	defer os.RemoveAll(aacDir)

	//ilistからaacファイルを取得し、TempDiraac()で作成したディレクトリに入れる
	if err := AACDownload(list, aacDir); err != nil {
		log.Fatalf("BulkDownload err : %v\n", err)
	}

	filename := fmt.Sprintf("%s_%s.%s", radiname, start, filetype)
	//生成された、aacファイルをconcatする
	if err := ConcatAACFile(ctx, aacDir, filename); err != nil {
		log.Fatalf("ConcatAACFile err : %v\n", err)
	}

	fmt.Printf("Success output %s\n", filename)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
