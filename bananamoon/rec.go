package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const (
	startLayout     = "20060102"
	startTimeLayout = "20060102150405"
)

var (
	aacFile    string
	aacCopyDir string
)

type Stations []Station

type Station struct {
	ID    string   `xml:"id,attr"`
	Name  string   `xml:"name"`
	Scd   Scd      `xml:"scd,omitempty"`
	Progs Programs `xml:"progs,omitempty"`
}

type Scd struct {
	Progs Programs `xml:"progs"`
}

type Programs struct {
	Date  string    `xml:"date"`
	Progs []Program `xml:"prog"`
}

type Program struct {
	Ft       string `xml:"ft,attr"`
	To       string `xml:"to,attr"`
	Ftl      string `xml:"ftl,attr"`
	Tol      string `xml:"tol,attr"`
	Dur      string `xml:"dur,attr"`
	Title    string `xml:"title"`
	SubTitle string `xml:"sub_title"`
	Desc     string `xml:"desc"`
	Pfm      string `xml:"pfm"`
	Info     string `xml:"info"`
	URL      string `xml:"url"`
}

type stationsData struct {
	XMLName     xml.Name `xml:"radiko"`
	XMLStations struct {
		XMLName  xml.Name `xml:"stations"`
		Stations Stations `xml:"station"`
	} `xml:"stations"`
}

func (c *Client) GetStartTime(ctx context.Context, stationID string, start time.Time) (*Program, error) {
	var err error
	if stationID == "" {
		return nil, errors.New("Station ID is nothing")
	}

	stations, err := c.GetStations(ctx, start)
	if err != nil {
		return nil, err
	}

	location, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("func ProgramDate time.LoadLocation err : %v \n", err)
	}
	localtime := start.In(location)
	ft := localtime.Format(startTimeLayout)

	var prog *Program
	for _, s := range stations {
		if s.ID == stationID {
			for _, p := range s.Progs.Progs {
				if p.Ft == ft {
					prog = &p
					break
				}
			}
		}
	}
	if prog == nil {
		return nil, errors.New("program is not found")
	}
	return prog, nil
}

func (c *Client) GetStations(ctx context.Context, start time.Time) (Stations, error) {
	stationEndpoint := path.Join("v3", "/program/date", ProgramDate(start), fmt.Sprintf("JP13.xml"))
	u := *c.URL
	u.Path = path.Join(c.URL.Path, stationEndpoint)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("pragma", "no-cache")
	var sdata stationsData
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = xml.Unmarshal(r, &sdata); err != nil {
		return nil, err
	}
	return sdata.stations(), nil
}

func (s *stationsData) stations() Stations {
	return s.XMLStations.Stations
}

func ProgramDate(start time.Time) string {
	var err error
	location, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("func ProgramDate time.LoadLocation err : %v \n", err)
	}
	localtime := start.In(location)
	h := localtime.Hour()
	if h >= 0 && h <= 4 {
		localtime = localtime.Add(-24 * time.Hour)
	}
	return localtime.Format(startLayout)
}

func AACDownload(list []string, aacDir string) error {
	var (
		judg bool
		wg   sync.WaitGroup
		ch   = make(chan struct{}, 64)
	)

	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			var err error
			for i := 0; i < 4; i++ {
				ch <- struct{}{}
				err = listup(link, aacDir)
				<-ch
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("Failed to download: %s", err)
				judg = true
			}
		}(v)
	}
	wg.Wait()

	if judg {
		return errors.New("Lack of aac files")
	}
	return nil
}

func listup(link, aacDir string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(link)
	file, err := os.Create(filepath.Join(aacDir, filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
