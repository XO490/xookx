package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var baseUrl = "https://www.okex.com"

const VERSION = "0.1.2"

func checkErr(err error, debug bool) {
	if err != nil {
		if debug {
			fmt.Print(err.Error())
		} else {
			fmt.Print("upss")
		}
		os.Exit(1)
	}
}

//https://www.okex.com/api/v5/market/index-components?index=BTC-USDT
func apiGetCurrency(index string) string {
	var url = baseUrl + "/api/v5/market/index-components?index=" + index
	//println(url)
	req, err := http.Get(url)
	checkErr(err, false)

	req.Header.Add("host", baseUrl)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.175 Safari/537.36")

	var jsonObj map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&jsonObj)
	checkErr(err, false)

	if jsonObj["code"] == "0" {
		//index := jsonObj["data"].(map[string]interface{})["index"].(string)
		lastFloat, _ := strconv.ParseFloat(jsonObj["data"].(map[string]interface{})["last"].(string), 32)
		last := fmt.Sprintf("%.2f", lastFloat)

		//fmt.Println(last)
		//fmt.Println(index + ": " + last)
		//fmt.Println(reflect.TypeOf(last))

		return last
	} else {
		return "</>: " + jsonObj["code"].(string)
	}
}

func help() {
	fmt.Print("xoOKX\n" +
		VERSION + "\n" +
		"---------\n" +
		"usage:\n" +
		"\t--index=BTC-USDT\t| return 'last' > https://www.okx.com/docs-v5/en/#rest-api-market-data-get-index-components\n" +
		"\t--sound\t\t\t| default sound notification (use --sound OR --custom-sound)\n" +
		"\t--custom-sound\t\t| your custom sound notification: --custom-sound=btc.mp3 (use --sound OR --custom-sound)\n" +
		"\nerror-code </>:\t> https://www.okx.com/docs-v5/en/#error-code\n" +
		"feedback:\t> https://t.me/xo490")
}

func notification(filename string) {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err, true)

	file, err := ioutil.ReadFile(pwd + "/media/" + filename)
	checkErr(err, true)

	dec, data, err := minimp3.DecodeFull(file)
	checkErr(err, true)

	context, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024)
	checkErr(err, true)

	var player = context.NewPlayer()
	player.Write(data)

	<-time.After(time.Second)

	defer context.Close()
	defer dec.Close()

	err = player.Close()
	checkErr(err, true)
}

func getArgs() {
	args := os.Args
	if len(args) < 2 {
		help()
	} else {
		index := flag.String("index", "", "--index=BTC-USDT")
		sound := flag.Bool("sound", false, "--sound")
		customSound := flag.String("custom-sound", "", "--custom-sound=up.mp3")

		flag.Parse()

		if *index != "" {
			fmt.Println(apiGetCurrency(*index))
		}

		if *sound {
			notification("default.mp3")
		}

		if *customSound != "" {
			notification(*customSound)
		}
	}
}

func main() {
	getArgs()
	//notification("up.mp3")
	//fmt.Print(apiGetCurrency("TONCOIN-USDT"))
}
