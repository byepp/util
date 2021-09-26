package httputil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/go-http-utils/headers"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)
//
//// AutoRetryResult 自动重试结果
//type AutoRetryResult struct {
//	Object          interface{}
//	ExecuteDuration time.Duration
//	Error           error
//	ResponseBody    json.RawMessage
//}
//
//// AutoRetryPostForm 带自动重试的PostForm
//func AutoRetryPostForm(postURL string, postData url.Values, maxRetryCount int, sleep time.Duration, chatRet chan AutoRetryResult, object interface{}) {
//	var ret AutoRetryResult
//	ret.Object = object
//	beginTime := time.Now()
//	ret.Error = looputil.Retry(maxRetryCount, sleep, func() error {
//		var err error
//		ret.ResponseBody, err = PostForm(postURL, postData)
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	ret.ExecuteDuration = time.Now().Sub(beginTime)
//	chatRet <- ret
//}
//
//// AutoRetryGet 带自动重试的Get
//func AutoRetryGet(uri string, maxRetryCount int, sleep time.Duration, chatRet chan<- AutoRetryResult, object interface{}) {
//	var ret AutoRetryResult
//	ret.Object = object
//	beginTime := time.Now()
//	ret.Error = looputil.Retry(maxRetryCount, sleep, func() error {
//		var err error
//		ret.ResponseBody, err = Get(uri)
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	ret.ExecuteDuration = time.Now().Sub(beginTime)
//	chatRet <- ret
//}

// Get 拉取网页内容
func Get(uri string) ([]byte, error) {
	timeout := 3 * time.Minute
	client := &http.Client{
		Timeout: timeout,
	}
	defer client.CloseIdleConnections()
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return ioutil.ReadAll(resp.Body)
}

func GetWithJsonDecode(uri string, v interface{}) error {
	body, err := Get(uri)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// PostForm 提交FORM类型表单
func PostForm(uri string, postData url.Values) ([]byte, error) {
	client := &http.Client{
		Timeout: 3 * time.Minute,
	}
	defer client.CloseIdleConnections()
	resp, err := client.PostForm(uri, postData)
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func PostFormJsonDecode(uri string, postData url.Values, v interface{}) error {
	client := &http.Client{
		Timeout: 3 * time.Minute,
	}
	defer client.CloseIdleConnections()
	resp, err := client.PostForm(uri, postData)
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// EncodeURIComponent URLEncode
func EncodeURIComponent(str string) string {
	ret := url.QueryEscape(str)
	ret = strings.Replace(ret, "+", "%20", -1)
	return ret
}

// PostJsonWithJsonDecode 发送HTTP请求输入输出参数都是json对象
func PostJsonWithJsonDecode(uri string, data interface{}, out interface{}) (err error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return json.Unmarshal(resBody, out)
}

func DoRequest(client *http.Client, req *http.Request) (resp *http.Response, body string, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	contentEncoding := resp.Header.Get(headers.ContentEncoding)
	var bodyReader io.Reader
	bodyReader = resp.Body
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") { // 使用了压缩
		bodyReader, _ = gzip.NewReader(bodyReader)
	}
	contentType := resp.Header.Get(headers.ContentType)
	rx := regexp.MustCompile(`charset=([\(\):\.\w-]+)`)
	ss := rx.FindStringSubmatch(contentType)
	if len(ss) == 2 { // 有识别出来编码
		contentCharset := ss[1]
		if strings.ToLower(contentCharset) != "utf-8" { // 不是UTF8的才需要转换
			var enc encoding.Encoding
			enc, err = ianaindex.IANA.Encoding(contentCharset)
			if err != nil {
				return
			}
			bodyReader = transform.NewReader(bodyReader, enc.NewDecoder())
			//if strings.Contains(strings.ToLower(contentType), "gbk") {// 使用了GBK编码
			//	bodyReader = transform.NewReader(bodyReader, simplifiedchinese.GBK.NewDecoder())
			//}
		}
	}
	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return
	}
	body = string(bodyBytes)
	return
}

func DoRequestBytes(client *http.Client, req *http.Request) (resp *http.Response, body []byte, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	contentEncoding := resp.Header.Get(headers.ContentEncoding)
	var bodyReader io.Reader
	bodyReader = resp.Body
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") { // 使用了压缩
		bodyReader, _ = gzip.NewReader(bodyReader)
	}
	body, err = ioutil.ReadAll(bodyReader)
	return
}

func DoRequestJsonDecode(client *http.Client, req *http.Request, res interface{}) (resp *http.Response, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(res)
	return
}

type ContentType string

const (
	ContentTypeJson ContentType = "application/json; charset=utf-8"
	ContentTypeAll  ContentType = "*/*; charset=utf-8"
)

const (
	UserAgent = "AdsBot-Google (+http://www.google.com/adsbot.html)"
)

func MakeCommonRequestJson(method string, uri string, data interface{}) *http.Request {
	body, _ := json.Marshal(data)
	return MakeCommonRequest(method, uri, ContentTypeJson, body)
}

func MakeCommonRequest(method string, uri string, contentType ContentType, body []byte) *http.Request {
	return MakeCommonRequestWithUserAgent(method, uri, contentType, UserAgent, body)
}

func MakeCommonRequestWithUserAgent(method string, uri string, contentType ContentType, userAgent string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, uri, bytes.NewReader(body))
	req.Header.Set(headers.ContentType, string(contentType))
	req.Header.Set(headers.AcceptLanguage, "zh-cn")
	req.Header.Set(headers.AcceptEncoding, "gzip,deflate")
	req.Header.Set(headers.Accept, "*/*")
	req.Header.Set(headers.UserAgent, userAgent)
	return req
}
