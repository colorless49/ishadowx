package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

var url string
var h bool

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.StringVar(&url, "url", "https://fast.ishadowx.net", "需要爬虫的网站地址。")
}
func main() {
	flag.Parse()

	if h {
		flag.Usage()
		fmt.Scanf("%s")
		return
	}

	cms := clientMultiServer{}
	cms.LocalPort = 1080
	var brookCmd string
	ssrs := make([]ssr, 0, 20)
	//Cookie:_ga=GA1.2.104527843.1524205208; _gid=GA1.2.190482586.1524789574
	//https://fast.ishadowx.net/  isx.yt     dwz.pm/x

	doc, err := newDocument(url)
	if err != nil {
		fmt.Println("create Document error.")
		return
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
	//time.Sleep(5000000000) //5s
	fmt.Scanf("%s")
}

func newDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Safari/537.36 LBBROWSER")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return goquery.NewDocumentFromReader(resp.Body)
	} else {
		return nil, errors.New("返回码不是200。")
	}

	//return goquery.NewDocumentFromResponse(resp)
}
