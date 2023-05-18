package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"xoOKX/cmd/text"
)

type Config struct {
	App struct {
		Debug   bool   `yaml:"debug"`
		Timeout int    `yaml:"timeout"`
		Proxy   string `yaml:"proxy"`
	}
	Headers struct {
		Host        string `yaml:"host"`
		ContentType string `yaml:"content-type"`
		UserAgent   string `yaml:"user-agent"`
		ApiUrl      string `yaml:"api-url"`
	}
}

var (
	config   = Config{}
	debug    = true
	version  = "0.1.4"
	pwd      = filepath.Dir(os.Args[0])
	yamlName = "config.yaml"
)

func readConfig() {
	file, err := os.ReadFile(filepath.Join(pwd, yamlName))
	if err != nil {
		log.Printf("[Error] readConfig ReadFile> %v\n", err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Printf("[Error] readConfig Unmarshal> %v\n", err)
	}
}

func createConfig() {
	data := "app:\n" +
		"  debug: false\n" +
		"  timeout: 5\n" +
		"#  proxy: \"http://195.154.43.182:58099\"\n" +
		"#  proxy url format - http://user:password@host:port\n\n" +
		"headers:\n" +
		"  host: \"https://www.okx.com\"\n" +
		"  content-type: \"application/json\"\n" +
		"  user-agent: \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.88 Safari/537.36 Vivaldi/5.1.2567.39\"\n" +
		"  api-url: \"/api/v5/market/index-components?index=\"\n"

	err := os.WriteFile(yamlName, []byte(data), 0755)
	checkErr(err, debug)
	fmt.Println("Config created.")
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

func apiGetCurrency(index string) string {
	/*
		example url: https://www.okx.com/api/v5/market/index-components?index=BTC-USDT
		api documentation: https://www.okx.com/docs-v5/en/#rest-api-market-data-get-index-components
	*/
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	var urlApi = config.Headers.Host + config.Headers.ApiUrl + index
	//var urlApi = config.Headers.Host + "/api/v5/market/index-components?index=" + index
	request, err := http.NewRequest(http.MethodGet, urlApi, nil)
	checkErr(err, debug)

	request.Header.Set("Host", config.Headers.Host)
	request.Header.Set("Content-Type", config.Headers.ContentType)
	request.Header.Set("User-Agent", config.Headers.UserAgent)

	client := &http.Client{
		Timeout: time.Second * time.Duration(config.App.Timeout),
	}

	if config.App.Proxy != "" {
		proxyUrl, _ := url.Parse(config.App.Proxy)
		proxy := &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client.Transport = proxy
	}

	response, _ := client.Do(request)

	var jsonObj map[string]interface{}
	if response != nil {
		err = json.NewDecoder(response.Body).Decode(&jsonObj)
		checkErr(err, debug)

		if jsonObj["code"] == "0" {
			lastFloat, _ := strconv.ParseFloat(jsonObj["data"].(map[string]interface{})["last"].(string), 32)
			last := fmt.Sprintf("%.2f", lastFloat)
			return last
		} else {
			return "</>: " + jsonObj["code"].(string)
		}
	}
	return "--"
}

func notification(filename string) {
	fileBytes, err := os.ReadFile("media/" + filename)
	checkErr(err, debug)

	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	samplingRate := 44100
	numOfChannels := 2
	audioBitDepth := 2

	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan

	player := otoCtx.NewPlayer(decodedMp3)

	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}

func getArgs() {
	args := os.Args
	if len(args) < 2 {
		help()
	} else {
		index := flag.String("index", "", text.Index)
		sound := flag.Bool("sound", false, text.Sound)
		customSound := flag.String("custom-sound", "", text.CustomSound)
		config := flag.Bool("config", false, text.Config)

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

func help() {
	fmt.Print("xoOKX\n" +
		version + "\n" +
		"----\n" +
		"usage of " + os.Args[0] + ":\n" +
		"--index=BTC-USDT\n\t" + text.Index +
		"\n\n--sound\n\t" + text.Sound +
		"\n\n--custom-sound\n\t" + text.CustomSound +
		"\n\n--config\n\t" + text.Config +
		"\n\nerror-code </>: " + text.ErrorCode +
		"\nfeedback: " + text.Contact)
}

func main() {
	readConfig()
	getArgs()
}
