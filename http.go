/*
 * @Author: ccchieh
 * @Github: https://github.com/ccchieh
 * @Email: email@zzj.cool
 * @Date: 2021-01-18 16:10:17
 * @LastEditors: ccchieh
 * @LastEditTime: 2021-05-22 19:39:06
 */
package goutils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// URLAddQuery 提供一个URL，然后添加query参数
func URLAddQuery(addr *url.URL, key, value string) {
	query := addr.Query()
	query.Add(key, value)
	addr.RawQuery = query.Encode()
}

func HttpRequest(method string, url string, data interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	client := &http.Client{}
	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = transport
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	//res := make(map[string]interface{})
	//if err := json.Unmarshal(body, &res); err != nil {
	//	return nil, err
	//}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprint(body))
	}
	return body, err
}

// ReverseProxy 用于反向代理，endpoint是代理的地址
func ReverseProxy(endpoint string, rw http.ResponseWriter, req *http.Request) {
	target, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport
	proxy.ModifyResponse = func(r *http.Response) error {
		r.Header.Del("X-Frame-Options") // 重点：代理时移除 X-Frame-Options 头
		return nil
	}
	proxy.ErrorHandler = handlerError
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Header.Set("X-Forwarded-Proto", getScheme(req))
	proxy.ServeHTTP(rw, req)
}

func LegalPort(port uint) bool {
	return port <= 65535
}

func handlerError(rw http.ResponseWriter, req *http.Request, err error) {
	// global.Log.Error(err.Error())
	_, _ = rw.Write([]byte("404 page not found"))
	rw.WriteHeader(http.StatusNotFound)
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	} else {
		return "http"
	}
}
