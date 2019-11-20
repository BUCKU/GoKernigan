package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main(){
	start := time.Now()
	ch := make(chan string)

	sites := getAlexaUrls("alexa_top_100")

	for _, url := range sites {
		go fetch(url, ch)
	}
	f, err := os.OpenFile("SitesData.txt", os.O_CREATE|os.O_APPEND, 0644)
	if err!=nil{
		fmt.Printf("%v", err)
	}
	defer f.Close()
	for range sites {
		str := <-ch
		f.WriteString(str)
		fmt.Printf(str)
	}
	secs := time.Since(start).Seconds()
	fmt.Printf("Elapsed: %.2fs", secs)
}

func fetch(url string, ch chan<- string) {
	start := time.Now()

	if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://"){
		url = "http://"+url
	}
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("fetch: %v\n", err)
		return
	}

	b, err := io.Copy(ioutil.Discard,resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("copy: %v\n", err)
	}
	resp.Body.Close()
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("\n----------\nUrl:%s,\nStatusCode: %d,\nSize: %d,\nTime: %.2fs\n----------", url, resp.StatusCode, b, secs)
}

func getAlexaUrls(path string) []string {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err!= nil {
		panic("Cannot read alexa top 100 file!")
	}

	var rets []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		s := scanner.Text()
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "<a href=\"/siteinfo/"){
			s = strings.TrimPrefix(s, "<a href=\"/siteinfo/")
			i := strings.Index(s, "\"")
			s = s[:i]
			rets = append(rets, s)
		}
	}

	defer f.Close()
	return rets
}