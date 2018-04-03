package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

var conf struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Task     int    `json:"task"`
}

func RandInt(n int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return rnd.Intn(n)
}

func init() {
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

	log.Println(runtime.GOOS)
}

func main() {
	cd := agouti.ChromeDriver()
	// cd := agouti.PhantomJS()
	// cd := agouti.Selenium()
	// cd := agouti.Selendroid("./selenium-server-standalone-3.11.0.jar")

	err := cd.Start()
	if err != nil {
		fmt.Println("Start webDriver:", err)
		return
	}
	defer cd.Stop()

	// pg, err := cd.NewPage(agouti.Browser("chrome"))
	pg, err := cd.NewPage(agouti.Browser("safari"))
	if err != nil {
		fmt.Println("NewPage:", err)
		return
	}

	err = pg.SetPageLoad(10000)
	if err != nil {
		fmt.Println("SetPageLoad:", err)
		return
	}

	err = pg.SetScriptTimeout(3000)
	if err != nil {
		fmt.Println("SetScriptTimeout:", err)
		return
	}

	err = pg.Navigate("http://mf.91yunma.cn")
	if err != nil {
		fmt.Println("Navigate:", err)
		return
	}

	elem := pg.FindByID("username")
	elem.SendKeys(conf.UserName)
	elem = pg.FindByID("password")
	elem.SendKeys(conf.Password)
	elem = pg.FindByID("vcode")
	elem.Click()

	// elem = driver.find_element_by_id("btnLogin")
	// elem.click()
	for {
		flag, _ := pg.URL()
		if strings.Contains(flag, "admin/desktop") {
			log.Println("got " + flag)
			break
		}
		time.Sleep(time.Second)
	}

	err = pg.Navigate("http://mf.91yunma.cn/admin/qpay/get_tasks")
	if err != nil {
		fmt.Println("Navigate:", err)
		return
	}

	for {
		flag, _ := pg.Title()
		if strings.Contains(flag, "待处理订单") {
			log.Println("got " + flag)
			break
		}
		time.Sleep(time.Second)
	}

	isAuto := false
	js := fmt.Sprintf("getTask(%d)", conf.Task)
	for {
		var ret interface{}
		err := pg.RunScript(js, nil, &ret)
		log.Printf("getTask result: %+v\n", ret)
		if err != nil {
			fmt.Println("getTask error. ", err)
			continue
		}

		time.Sleep(time.Second)
		flag, _ := pg.HTML()
		if strings.Contains(flag, "扫码代充") {
			log.Println("等待充值 10 秒...")
			time.Sleep(10 * time.Second)
			continue
		}

		if isAuto == false {
			elem := pg.FindByName("auto")
			if elem == nil {
				continue
			}
			elem.Click()
			isAuto = true
		}

		elem := pg.FindByID("submit")
		elem.Click()
		itv := RandInt(2000) + 4000
		time.Sleep(time.Duration(itv) * time.Millisecond)
	}
}
