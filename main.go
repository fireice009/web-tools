package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var counter int
var total int
var wg sync.WaitGroup

func main() {
	rst := map[string]int{}
	b, err := ioutil.ReadFile("src.csv")
	if err != nil {
		panic(err)
	}
	r2 := csv.NewReader(strings.NewReader(string(b)))
	ss, _ := r2.ReadAll()
	for _, item := range ss {
		if len(item) == 0 {
			continue
		}
		rst[item[0]] = -1
	}

	ch := make(chan string, 128)
	defer close(ch)
	total = len(rst)
	wg.Add(total)
	go product(rst, ch)
	for i := 0; i < 3; i++ {
		go custom(rst, ch)
	}
	wg.Wait()

	f, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(f)
	for url, wgh := range rst {
		w.Write([]string{url, strconv.Itoa(wgh)})
	}
	w.Flush()
}

func product(rst map[string]int, ch chan string) {
	for url, wgh := range rst {
		if wgh >= 0 {
			continue
		}
		ch <- url
	}
}

func custom(rst map[string]int, ch <-chan string) {
	for url := range ch {
		w, err1 := weight(url)
		if err1 != nil {
			fmt.Println(err1)
			w = -1
		}
		rst[url] = w
		counter++
		fmt.Printf("%4d/%d  %4.1f%%  %s=%d\n", counter, total, float32(counter)/float32(total)*100, url, w)
		wg.Done()
	}
}

func weight(domain string) (int, error) {
	url := fmt.Sprintf("http://baidurank.aizhan.com/api/br?domain=%s&style=text", domain)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	idx := strings.Index(strBody, ">")
	return strconv.Atoi(strBody[idx+1 : idx+2])
}
