package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var analyCode map[string]string

func init() {
	analyCode = make(map[string]string)
	analyCode["&#xe602;"] = "1"
	analyCode["&#xe603;"] = "0"
	analyCode["&#xe604;"] = "3"
	analyCode["&#xe605;"] = "2"
	analyCode["&#xe606;"] = "4"
	analyCode["&#xe607;"] = "5"
	analyCode["&#xe608;"] = "6"
	analyCode["&#xe609;"] = "9"
	analyCode["&#xe60a;"] = "7"
	analyCode["&#xe60b;"] = "8"
	analyCode["&#xe60c;"] = "4"
	analyCode["&#xe60d;"] = "0"
	analyCode["&#xe60e;"] = "1"
	analyCode["&#xe60f;"] = "5"
	analyCode["&#xe610;"] = "2"
	analyCode["&#xe611;"] = "3"
	analyCode["&#xe612;"] = "6"
	analyCode["&#xe613;"] = "7"
	analyCode["&#xe614;"] = "8"
	analyCode["&#xe615;"] = "9"
	analyCode["&#xe616;"] = "0"
	analyCode["&#xe617;"] = "2"
	analyCode["&#xe618;"] = "1"
	analyCode["&#xe619;"] = "4"
	analyCode["&#xe61a;"] = "3"
	analyCode["&#xe61b;"] = "5"
	analyCode["&#xe61c;"] = "7"
	analyCode["&#xe61d;"] = "8"
	analyCode["&#xe61e;"] = "9"
	analyCode["&#xe61f;"] = "6"
}

func GetData(url string) (postNum, nickName string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	rand.Seed(time.Now().Unix())
	s := strconv.Itoa(rand.Intn(1000))

	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/" + s  + ".36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cookie", "_ga=GA1.2.685263550.1587277283; _gid=GA1.2.143250871.1587911549; tt_webid=6820028204934923790; _ba=BA0.2-20200301-5199e-c7q9NP0laGm7KfaPfGcH")
	res, err := client.Do(req)
	if err == nil {
		b, _ := ioutil.ReadAll(res.Body)
		result := string(b)
		var itemIDRegexp= regexp.MustCompile(`<div class="user-tab active tab get-list" data-type="post">作品<span class="num">(.*?)</span>`)
		ids := itemIDRegexp.FindStringSubmatch(result)
		if len(ids) == 0 {
			itemIDRegexp = regexp.MustCompile(`<div class="user-tab tab get-list" data-type="post">作品<span class="num">(.*?)</span>`)
			ids = itemIDRegexp.FindStringSubmatch(result)
		}
		if len(ids) > 0 {
			var numRegExp = regexp.MustCompile(`<i class="icon iconfont tab-num"> (.*?) </i>`)
			nums := numRegExp.FindAllStringSubmatch(ids[0], -1)
			for i := 0; i < len(nums); i++ {
				postNum += analyCode[nums[i][1]]
			}
		}

		var nickNameRegexp = regexp.MustCompile(`<p class="nickname">(.*?)</p>`)

		nickNames := nickNameRegexp.FindStringSubmatch(result)
		if len(nickNames) > 0 {
			nickName = nickNames[1]
		}
	}
	return
}

func read3(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}
	return string(fd), nil
}

type Config struct {
	Delay int `json:"delay"`
}


func ParserConfig(data string) Config {
	var config Config
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}
	return config
}

func main()  {
	file, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	info := log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime)

	data, err := read3("./config.json")
	if err != nil {
		panic(err)
	}
	config := ParserConfig(data)

	//读取账号
	user, err := read3("./user.txt")
	uids := strings.Split(user, "\r\n")
	for i := 0; i < len(uids); i++ {
		post, nickName := GetData(uids[i])
		info.Println("昵称:" + nickName + " 作品数:" + post + " 用户主页:" + uids[i])
		time.Sleep(time.Duration(config.Delay) * time.Second)
	}
}
