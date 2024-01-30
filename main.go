package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies/shadowsocks"
	"github.com/deanxv/yescaptcha-go"
	"github.com/deanxv/yescaptcha-go/req"
	"io"
	"net/http"
	"time"
)

func main() {
	apiKey := "xxx"
	cli := yescaptcha.NewClient(apiKey, "33989", "https://api.yescaptcha.com")
	task := req.NoCaptchaTaskProxylessRequest{
		WebsiteURL: "https://artio.faucet.berachain.com",
		Type:       "RecaptchaV3TaskProxylessM1",
		WebsiteKey: "6LfOA04pAAAAAL9ttkwIz40hC63_7IsaU2MgcwVH",
		PageAction: "driptoken",
	}
	resp, err := cli.CreateNoCaptchaTaskProxyless(&task)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(40 * time.Second)
	result, err := cli.GetTaskResult(resp.TaskId)
	if err != nil {
		fmt.Println(err)
		return
	}
	authorization := "Bearer " + result.Solution.GRecaptchaResponse
	fmt.Println(authorization)
	url := "https://artio-80085-faucet-api-recaptcha.berachain.com/api/claim?address=0x426D2B685259d2bB75F2fe312D9b79289b3C5DD3"

	marshal, err := json.Marshal(map[string]string{"address": "0x426D2B685259d2bB75F2fe312D9b79289b3C5DD3"})
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	if err != nil {
		fmt.Println(err)
		return
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
	httpCli := SSClientYaml(&http.Client{})
	do, err := httpCli.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(do.Status)
	all, err := io.ReadAll(do.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(all))
}

func SSClientYaml(cli *http.Client) *http.Client {
	dialer, err := shadowsocks.NewDialerWithCfg(proxies.StringResolver, "ss.yaml")
	if err != nil {
		panic(err)
	}
	return proxies.NewHttpClient(cli, dialer)
}
