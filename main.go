package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies/shadowsocks"
	"github.com/deanxv/yescaptcha-go"
	"github.com/deanxv/yescaptcha-go/req"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	regExp         = `(\d{1,2})h(\d{1,2})m(\d{1,2})s`
	apiKey         = "xxx"
	yescaptchaWait = 80 * time.Second
)

func drip() (time.Duration, error) {
	cli := yescaptcha.NewClient(apiKey, "33989", "https://api.yescaptcha.com")
	task := req.TurnstileTaskProxylessRequest{
		WebsiteURL: "https://artio.faucet.berachain.com",
		Type:       "TurnstileTaskProxylessM1",
		WebsiteKey: "0x4AAAAAAARdAuciFArKhVwt",
	}
	resp, err := cli.CreateTurnstileTaskProxyless(&task)
	if err != nil {
		return 0, (err)
	}
	time.Sleep(yescaptchaWait)
	result, err := cli.GetTaskResult(resp.TaskId)
	if err != nil {
		return 0, (err)
	}
	authorization := "Bearer " + result.Solution.Token
	log.Println(authorization)
	//url := "https://artio-80085-faucet-api-recaptcha.berachain.com/api/claim?address=0x961dfB987266e3D5029713E7af4F989a41eB961A"
	url := "https://artio-80085-faucet-api-cf.berachain.com/api/claim?address=0x961dfB987266e3D5029713E7af4F989a41eB961A"
	marshal, err := json.Marshal(map[string]string{"address": "0x961dfB987266e3D5029713E7af4F989a41eB961A"})
	if err != nil {
		return 0, (err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	if err != nil {
		return 0, (err)
	}
	request.Header.Set("authority", "artio-80085-faucet-api-recaptcha.berachain.com")
	request.Header.Set("accept", "*/*")
	request.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("content-type", "text/plain;charset=UTF-8")
	request.Header.Set("origin", "https://artio.faucet.berachain.com")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("referer", "https://artio.faucet.berachain.com/")
	request.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	request.Header.Set("Authorization", authorization)
	httpCli := (&http.Client{})
	do, err := httpCli.Do(request)
	if err != nil {
		return 0, (err)
	}
	if do.StatusCode == http.StatusOK {
		log.Println("drop ok")
		d, _ := time.ParseDuration("8h")
		return d, nil
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return 0, (err)
	}
	s := string(all)
	log.Println(s)
	return estimatedTime(s)
}

func estimatedTime(s string) (est time.Duration, err error) {
	reg := regexp.MustCompile(regExp)
	match := reg.FindStringSubmatch(s)
	if len(match) < 1 {
		return 0, errors.New("parse time failed: " + s)
	}
	est, err = time.ParseDuration(match[0])
	if err != nil {
		return 0, errors.New("parseDuration failed: " + err.Error())
	}
	est -= yescaptchaWait
	//for i, sd := range match[1:] {
	//	n, _ := strconv.Atoi(sd)
	//	switch i {
	//	case 0:
	//		est += time.Hour * time.Duration(n)
	//	case 1:
	//		est += time.Minute * time.Duration(n)
	//	case 2:
	//		est += time.Second * time.Duration(n)
	//	}
	//}
	return
}

func main() {
	for {
		d, err := drip()
		if err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(d)
	}
}

func SSClientYaml(cli *http.Client) *http.Client {
	dialer, err := shadowsocks.NewDialerWithCfg(proxies.StringResolver, "ss.yaml")
	if err != nil {
		panic(err)
	}
	return proxies.NewHttpClient(cli, dialer)
}
