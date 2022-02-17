package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajimehoshi/oto"
	"github.com/spf13/viper"
	"github.com/tosone/minimp3"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	App struct {
		debug bool
	}
	Headers struct {
		head        string
		contentType string
		userAgent   string
	}
}

var CONFIG = Config{}
var DEBUG = true
var PWD string

func readConfig() {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err, DEBUG)
	PWD = pwd

	y := viper.New()
	y.SetConfigFile(PWD + "/config.yaml")
	err = y.ReadInConfig()
	if err != nil {
		createConfig()
	}

	CONFIG.Headers.head = y.GetString("headers.head")
	CONFIG.Headers.contentType = y.GetString("headers.content-type")
	CONFIG.Headers.userAgent = y.GetString("headers.user-agent")
	CONFIG.App.debug = y.GetBool("app.debug")
	DEBUG = CONFIG.App.debug
}

func createConfig() {
	data := "app:\n" +
		"  debug: false\n\n" +
		"headers:\n" +
		"  head: https://www.okex.com\n" +
		"  content-type: application/json\n" +
		"  user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.88 Safari/537.36 Vivaldi/5.1.2567.39\n"

	err := ioutil.WriteFile(PWD+"/config.yaml", []byte(data), 0755)
	checkErr(err, DEBUG)
	fmt.Println("Config created in path: ", PWD)
}

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

func getNumOfBuild() string {
	//pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//checkErr(err, DEBUG)

	build, err := os.ReadFile(PWD + "/.build")
	if err != nil {
		err := os.WriteFile(PWD+"/.build", []byte("0"), 0755)
		checkErr(err, DEBUG)

		build, err = os.ReadFile(PWD + "/.build")
		checkErr(err, DEBUG)
	}

	oldbuild, _ := strconv.Atoi(string(build))
	newbuild := strconv.Itoa(oldbuild + 1)

	err = os.WriteFile(PWD+"/.build", []byte(newbuild), 0755)
	checkErr(err, DEBUG)
	return newbuild
}

func apiGetCurrency(index string) string {
	/*
		example url: https://www.okex.com/api/v5/market/index-components?index=BTC-USDT
		api documentation: https://www.okx.com/docs-v5/en/#rest-api-market-data-get-index-components
	*/
	var url = CONFIG.Headers.head + "/api/v5/market/index-components?index=" + index
	req, err := http.Get(url)
	checkErr(err, DEBUG)

	req.Header.Add("host", CONFIG.Headers.head)
	req.Header.Add("content-type", CONFIG.Headers.contentType)
	req.Header.Add("user-agent", CONFIG.Headers.userAgent)

	var jsonObj map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&jsonObj)
	checkErr(err, false)

	if jsonObj["code"] == "0" {
		lastFloat, _ := strconv.ParseFloat(jsonObj["data"].(map[string]interface{})["last"].(string), 32)
		last := fmt.Sprintf("%.2f", lastFloat)

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
		"\t--sound\t\t\t| default sound notification: --sound (use --sound OR --custom-sound)\n" +
		"\t--custom-sound\t\t| your custom sound notification: --custom-sound=btc.mp3 (use --sound OR --custom-sound)\n" +
		"\t--config\t\t\t| create config.yaml file\n" +
		"\nerror-code </>:\t> https://www.okx.com/docs-v5/en/#error-code\n" +
		"feedback:\t> https://t.me/xo490")
}

func notification(filename string) {
	//pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//checkErr(err, DEBUG)

	file, err := ioutil.ReadFile(PWD + "/media/" + filename)
	checkErr(err, DEBUG)

	dec, data, err := minimp3.DecodeFull(file)
	checkErr(err, DEBUG)

	context, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024)
	checkErr(err, DEBUG)

	var player = context.NewPlayer()
	player.Write(data)

	<-time.After(time.Second)

	defer context.Close()
	defer dec.Close()

	err = player.Close()
	checkErr(err, DEBUG)
}

func getArgs() {
	args := os.Args
	if len(args) < 2 {
		help()
	} else {
		index := flag.String("index", "", "--index=BTC-USDT")
		sound := flag.Bool("sound", false, "--sound")
		customSound := flag.String("custom-sound", "", "--custom-sound=up.mp3")
		config := flag.Bool("config", false, "--config")

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

		if *config {
			createConfig()
		}
	}
}

/*
	Building app
*/

var VERSION string

func buildApp(release bool, version string, buildNum string) {
	readConfig()

	if release {
		if buildNum != "" {
			VERSION = version + " | build " + buildNum + " | release"
		} else {
			fmt.Println("For RELEASE=true buildNum can't be empty")
			os.Exit(1)
		}
	}

	if release != true {
		buildNum = getNumOfBuild()
		VERSION = version + " | build " + getNumOfBuild() + " | debug"
	}

	getArgs()
}

func main() {
	buildApp(true, "0.1.3", "71")
}
