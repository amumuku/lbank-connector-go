package sve

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"github.com/LBank-exchange/lbank-connector-go/pkg" // 使用完整路径
)

type HttpService struct {
	c         *Client
	ReqObj    *http.Request
	RespObj   *http.Response
	Body      string
	CostTime  int64
	Method    string
	Headers   map[string]string
	Text      string
	Content   []byte
	IsEchoReq bool
	isDebug   bool

	Error error

	EchoStr         string `json:"echostr"`
	Timestamp       string `json:"timestamp"`
	SignatureMethod string `json:"signature_method"`
}

var defaultHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

var tr = &http.Transport{
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	MaxIdleConnsPerHost: 2000,
}

type KwArgs func(hs *HttpService)

func WithHeaders(headers map[string]string) KwArgs {
	return func(hs *HttpService) {
		hs.Headers = headers
	}
}
func WithDebug(debug bool) KwArgs {
	return func(hs *HttpService) {
		hs.isDebug = debug
	}
}

func WithParams(params map[string]string) KwArgs {
	urlParams := url.Values{}
	for k, v := range params {
		urlParams.Add(k, v)
	}
	return func(hs *HttpService) {
		hs.ReqObj.URL.RawQuery = urlParams.Encode()
	}
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (hs *HttpService) Get(url, data string, kwargs ...KwArgs) *http.Response {
	var newUrl string
	if len(data) > 0 {
		newUrl = fmt.Sprintf("%s?%s", url, data)
	} else {
		newUrl = url
	}
	text, err := hs.DoHttpRequest("GET", newUrl, "", kwargs...)
	hs.Error = err
	hs.Text = text
	if err != nil {
		hs.c.debug("[reqErr] %s\n" + err.Error())
	}
	return hs.RespObj
}

func (hs *HttpService) Post(url, json string, kwargs ...KwArgs) *http.Response {
	text, err := hs.DoHttpRequest("POST", url, json, kwargs...)
	hs.Text = text
	if err != nil {
		hs.c.Logger.Error("[reqErr] %s\n" + err.Error())
	}
	return hs.RespObj
}

func (hs *HttpService) IsPrintReq(isEchoReq bool) *HttpService {
	hs.IsEchoReq = isEchoReq
	return hs
}

func (hs *HttpService) DoHttpRequest(method, url, body string, kwargs ...KwArgs) (string, error) {
	hs.BuildHeader()
	client := hs.BuildClient()
	req, err := hs.BuildRequest(method, url, body)
	if err != nil {
		return "", err
	}
	hs.ReqObj = req

	for _, kw := range kwargs {
		kw(hs)
	}

	if len(hs.Headers) > 0 {
		hs.BuildRequestHeaders(hs.ReqObj, hs.Headers)
	} else {
		hs.BuildRequestHeaders(hs.ReqObj, defaultHeaders)
	}
	startTime := time.Now()
	respObj, err := client.Do(hs.ReqObj)
	hs.RespObj = respObj
	elapsed := time.Since(startTime).Nanoseconds() / int64(time.Millisecond)
	hs.CostTime = elapsed
	if hs.IsEchoReq || hs.isDebug || hs.c.Debug {
		hs.PrintReqInfo(hs.ReqObj)
	}
	if err != nil {
		hs.c.Logger.Error("[reqErr] %s\n" + err.Error())
		hs.PrintReqInfo(hs.ReqObj)
		return "", err
	}
	defer respObj.Body.Close()
	content, err := io.ReadAll(respObj.Body)
	hs.Content = content
	if err != nil {
		hs.c.Logger.Error("[RespErr]%s\n" + err.Error())
		hs.PrintReqInfo(hs.ReqObj)
		hs.PrintRespInfo(content, elapsed)
		return "", err
	}
	if hs.isDebug || hs.c.Debug {
		hs.PrintRespInfo(content, elapsed)
	}
	return string(content), nil
}

func (hs *HttpService) BuildRequest(method, url string, body string) (req *http.Request, err error) {
	hs.Body = body
	b := hs.BuildBody(body)
	req, err = http.NewRequest(method, url, b)
	if err != nil {
		hs.c.Logger.Error("BuildRequestErr" + err.Error())
		return nil, err
	}
	return req, nil
}

func (hs *HttpService) BuildClient() *http.Client {
	client := &http.Client{Timeout: 3 * 60 * time.Second, Transport: tr}
	return client
}

func (hs *HttpService) BuildRequestHeaders(req *http.Request, headers map[string]string) *HttpService {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return hs
	//{"Content-Type": 'application/x-www-form-urlencoded',
	//            "signature_method": self.sign_method,
	//            'timestamp': self.timestamp,
	//            'echostr': self.random_str}
}

func (hs *HttpService) PrintReqInfo(req *http.Request) {
	s := fmt.Sprintf("\n    [ReqHeaders]:%v", req.Header) + fmt.Sprintf("\n    [ReqMethod]:%s", req.Method) +
		fmt.Sprintf("\n    [ReqUrl]:%s", req.URL) + fmt.Sprintf("\n    [ReqBody]:%s", hs.Body)
	hs.c.debug(s)
}

func (hs *HttpService) BuildBody(body string) *strings.Reader {
	hs.Body = body
	return strings.NewReader(body)
}

func (hs *HttpService) PrintRespInfo(resInfo []byte, costTime int64) *HttpService {
	costFloat := float64(costTime) / 1.0e9
	formatCostTime := fmt.Sprintf("%.3f", costFloat)
	hs.CostTime = costTime / 1e6
	r, _ := pkg.PrettyPrint(resInfo)
	s := fmt.Sprintf("\n    [RespHttpCode]:%d", hs.RespObj.StatusCode) + fmt.Sprintf("\n    [RespCost]:%sSecond",
		formatCostTime) + fmt.Sprintf("\n    [RespBody]:%s", r)
	hs.c.debug(s)
	return hs
}

func (hs *HttpService) PrettyPrint(resInfo []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, resInfo, "", " "); err != nil {
		return string(resInfo), err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func (hs *HttpService) Map2String(body map[string]interface{}) string {
	return pkg.Map2JsonString(body)
}

// Json https://github.com/tidwall/gjson
func (hs *HttpService) Json() gjson.Result {
	return gjson.Parse(hs.Text)
}

func (hs *HttpService) InitTsAndStr() {
	hs.Timestamp = pkg.Timestamp()
	hs.EchoStr = pkg.RandomStr()
}
func (hs *HttpService) BuildHeader() *HttpService {
	hd := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded",
		"signature_method": "RSA",
		"timestamp":        hs.Timestamp,
		"echostr":          hs.EchoStr}
	hs.Headers = hd
	return hs
}

func (hs *HttpService) BuildSignBody(kwargs map[string]string) string {
	hs.InitTsAndStr()
	kwargs["api_key"] = hs.c.ApiKey
	kwargs["timestamp"] = hs.Timestamp
	if len(hs.c.SecretKey) > 0 {
		kwargs["signature_method"] = "RSA"
	} else {
		kwargs["signature_method"] = "HmacSHA256"
	}
	kwargs["echostr"] = hs.EchoStr

	paramsSortStr := pkg.FormatStringBySign(kwargs)
	var sign string
	if len(hs.c.SecretKey) > 0 {
		sign, _ = hs.BuildRsaSignV2(paramsSortStr, hs.c.SecretKey)
	} else {
		sign, _ = hs.BuildHmacSignV2(paramsSortStr, hs.c.SecretKey)
	}
	kwargs["sign"] = sign
	postData := url.Values{}
	for k, v := range kwargs {
		postData.Add(k, v)
	}
	postData.Del("echostr")
	postData.Del("timestamp")
	postData.Del("signature_method")
	hs.Body = postData.Encode()
	return postData.Encode()
}

// BuildRsaSignV2
func (hs *HttpService) BuildRsaSignV2(params, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("secret is empty")
	}
	b := []byte("-----BEGIN RSA PRIVATE KEY-----\n" + secret + "\n-----END RSA PRIVATE KEY-----")
	privateKey, err := pkg.ParsePKCS1PrivateKey(b)
	if err != nil {
		return "", err
	}
	return pkg.RSASign(params, privateKey), nil
}

func (hs *HttpService) BuildHmacSignV2(params, secret string) (string, error) {
	return pkg.HmacSha256Base64Signer(params, secret)
}
