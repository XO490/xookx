package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var baseUrl = "https://www.okex.com"

const VERSION = "0.1.1"

func checkErr(err error) {
	if err != nil {
		fmt.Print("upss")
		//fmt.Print(err.Error())
		os.Exit(1)
	}
}

//https://www.okex.com/api/v5/market/index-components?index=BTC-USDT
func apiGetCurrency(index string) string {
	var url = baseUrl + "/api/v5/market/index-components?index=" + index
	//println(url)
	req, err := http.Get(url)
	checkErr(err)

	req.Header.Add("host", baseUrl)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.175 Safari/537.36")

	var jsonObj map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&jsonObj)
	checkErr(err)

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
		"\t-i, --index=BTC-USDT\t| return 'last' > https://www.okx.com/docs-v5/en/#rest-api-market-data-get-index-components\n" +
		"\nerror-code </>:\t> https://www.okx.com/docs-v5/en/#error-code\n" +
		"feedback:\t> https://t.me/xo490")
}

func getArgs() {
	args := os.Args
	if len(args) < 2 {
		help()
	} else {
		index := flag.String("index", "", "--index=BTC-USDT")
		flag.Parse()
		if *index != "" {
			fmt.Println(apiGetCurrency(*index))
		}
	}
}

func main() {
	getArgs()
	//fmt.Print(apiGetCurrency("TONCOIN-USDT"))
}
