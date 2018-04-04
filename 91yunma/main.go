package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

var taskConfig = map[int]int{
	500: 6,
	300: 5,
	200: 2,
	100: 1,
}

var conf struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Money    int    `json:"money"`
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

	if _, ok := taskConfig[conf.Money]; ok == false {
		panic(fmt.Sprintf("config error. money cannot be %d", conf.Money))
	}
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
	js := fmt.Sprintf("getTask(%d)", taskConfig[conf.Money])
	for {
		var ret interface{}
		err := pg.RunScript(js, nil, &ret)
		log.Printf("get task result: %+v\n", ret)
		if err != nil {
			fmt.Println("get task error. ", err)
			continue
		}
		for i := 0; i < 10; i++ {
			flag, _ := pg.HTML()
			if strings.Contains(flag, "获取代充订单") {
				break
			}
			time.Sleep(time.Second)
		}

		if isAuto == false {
			elem := pg.FindByName("auto")
			elem.Click()
			isAuto = true
		}
		elem := pg.FindByID("submit")
		elem.Click()

		flag, _ := pg.HTML()
		if strings.Contains(flag, "扫码代充") {
			for {
				log.Println("等待充值...")
				time.Sleep(2 * time.Second)
				flag, _ := pg.HTML()
				if strings.Contains(flag, "扫码代充") == false {
					break
				}
			}
		} else {
			itv := RandInt(2000) + 4000
			time.Sleep(time.Duration(itv) * time.Millisecond)
		}
	}
}
