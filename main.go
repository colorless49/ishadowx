package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ssr struct {
	IP       string
	Port     string
	Password string
	Method   string
}

type clientMultiServer struct {
	LocalPort      int         `json:"local_port"`
	ServerPassword [][3]string `json:"server_password"`
}

func main() {
	cms := clientMultiServer{}
	cms.LocalPort = 1080
	var brookCmd string
	ssrs := make([]ssr, 0, 20)

	doc, err := newDocument("https://free.ishadowx.net")
	if err != nil {
		fmt.Println("create Document error.")
	}
	// Find the review items
	doc.Find(".portfolio-item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		ssr := ssr{}
		s.Find("H4").Each(func(j int, e *goquery.Selection) {
			switch j {
			case 0:
				ssr.IP = e.Find("span").First().Text()
				fmt.Println("IP:" + ssr.IP)
			case 1:
				ssr.Port = strings.TrimSpace(e.Find("span").First().Text())
				fmt.Println("Port:" + ssr.Port)
			case 2:
				ssr.Password = strings.TrimSpace(e.Find("span").First().Text())
				fmt.Println("Password:" + ssr.Password)
			case 3:
				ssr.Method = strings.Split(e.Text(), ":")[1]
				fmt.Println("Method:" + ssr.Method)
			}
		})
		ssrs = append(ssrs, ssr)
		fmt.Println("--------------------------------")
	})

	for _, v := range ssrs {
		if brookCmd == "" && v.Method == "aes-256-cfb" {
			fmt.Print("brook CMD ----------:")
			fmt.Println(v)
			brookCmd = "brook ssclient -l 127.0.0.1:1080 -i 127.0.0.1 -s " + v.IP + ":" + v.Port + " -p " + v.Password + " --http"
		}

		fmt.Println(v)
		cms.ServerPassword = append(cms.ServerPassword, [3]string{v.IP + ":" + v.Port, v.Password, v.Method})
	}
	jsonstr, _ := json.Marshal(cms)
	if ioutil.WriteFile("client-multi-server.json", jsonstr, 0644) != nil {
		fmt.Println("写入client-multi-server.json失败!")
	}
	if ioutil.WriteFile("brook.bat", []byte(brookCmd), 0644) != nil {
		fmt.Println("写入brook.bat失败!")
	}
	time.Sleep(5000000000)
}

func newDocument(url string) (*goquery.Document, error) {
	// Load the URL
	res, e := http.Get(url) //根据url获取该网页的内容  res
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	return goquery.NewDocumentFromResponse(res)
}
