package sve

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/LBank-exchange/lbank-connector-go/pkg"
	"github.com/tidwall/gjson"
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

func (hs *HttpService) Json() gjson.Result {
	return gjson.Parse(hs.Text)
}

func (hs *HttpService) InitTsAndStr() {
	hs.Timestamp = fmt.Sprintf("%d", time.Now().UnixMilli()) // 毫秒级时间戳
	hs.EchoStr = randomString(30, 40)                        // 随机长度 30-40
}

func randomString(minLen, maxLen int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	mrand.Seed(time.Now().UnixNano()) // 初始化随机种子
	length := mrand.Intn(maxLen-minLen+1) + minLen
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[mrand.Intn(len(charset))]
	}
	return string(b)
}

func (hs *HttpService) detectSignatureMethod(secret string) string {
	// 如果 secret 是空，默认 HMAC-SHA256
	if len(secret) == 0 {
		return "HmacSHA256"
	}

	// 尝试解析为 PEM 格式的 RSA 私钥
	pemData := []byte(secret)
	if !strings.Contains(secret, "-----BEGIN") {
		pemData = []byte("-----BEGIN RSA PRIVATE KEY-----\n" + secret + "\n-----END RSA PRIVATE KEY-----")
	}
	block, _ := pem.Decode(pemData)
	if block != nil {
		if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return "RSA"
		}
		if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if _, ok := key.(*rsa.PrivateKey); ok {
				return "RSA"
			}
		}
	}
	// 如果无法解析为 RSA，则使用 HMAC-SHA256
	return "HmacSHA256"
}

func (hs *HttpService) BuildHeader() *HttpService {
	hd := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded",
		"signature_method": hs.SignatureMethod, // 使用检测到的签名方法
		"timestamp":        hs.Timestamp,
		"echostr":          hs.EchoStr,
	}
	hs.Headers = hd
	return hs
}

func (hs *HttpService) BuildSignBody(kwargs map[string]string) string {
	hs.InitTsAndStr()
	kwargs["api_key"] = hs.c.ApiKey
	kwargs["timestamp"] = hs.Timestamp

	// 自动检测签名方法
	hs.SignatureMethod = hs.detectSignatureMethod(hs.c.SecretKey)
	kwargs["signature_method"] = hs.SignatureMethod

	kwargs["echostr"] = hs.EchoStr

	// 按键名排序并拼接为 key=value&key=value 格式
	paramsSortStr := formatParams(kwargs)
	// 计算 MD5 并转为大写，与 Python 一致
	md5Digest := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(paramsSortStr))))

	var sign string
	var err error
	if hs.SignatureMethod == "RSA" {
		sign, err = hs.BuildRsaSignV2(md5Digest, hs.c.SecretKey)
	} else {
		sign, err = hs.BuildHmacSignV2(md5Digest, hs.c.SecretKey)
	}
	if err != nil {
		hs.c.Logger.Error(fmt.Sprintf("[SignErr] %s", err.Error()))
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

// 辅助函数：参数排序和拼接
func formatParams(params map[string]string) string {
	var keys []string
	for k := range params {
		if k != "sign" { // 忽略 sign 参数，与 Python 一致
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(pairs, "&")
}

func (hs *HttpService) BuildRsaSignV2(params, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("secret is empty")
	}

	// 调试输出原始密钥
	hs.c.Logger.Debug(fmt.Sprintf("Raw secret: %s", secret))

	// 如果 secret 不包含 PEM 头，添加 PKCS#1 格式头尾
	pemData := []byte(secret)
	if !strings.Contains(secret, "-----BEGIN") {
		pemData = []byte("-----BEGIN RSA PRIVATE KEY-----\n" + secret + "\n-----END RSA PRIVATE KEY-----")
	}

	// 解码 PEM 数据
	block, _ := pem.Decode(pemData)
	if block == nil {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	// 尝试解析为 PKCS#1 私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 如果 PKCS#1 解析失败，尝试解析为 PKCS#8 私钥
		key, errPKCS8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if errPKCS8 != nil {
			return "", fmt.Errorf("failed to parse private key - PKCS#1: %v, PKCS#8: %v", err, errPKCS8)
		}
		// 确保是 RSA 私钥
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("parsed PKCS#8 key is not an RSA private key")
		}
	}

	// 计算 SHA256 哈希
	h := sha256.New()
	h.Write([]byte(params))
	hashed := h.Sum(nil)

	// 使用 PKCS1v15 签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign with RSA: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (hs *HttpService) BuildHmacSignV2(params, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("secret is empty")
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(params))
	return fmt.Sprintf("%x", h.Sum(nil)), nil // 返回十六进制字符串
}
