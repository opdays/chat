package controllers

import (
	"github.com/astaxie/beego"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"encoding/base64"
	"github.com/astaxie/beego/httplib"
	"time"
	"sort"
	"strconv"
	"strings"
)

const (
	Method    = "GET"
	MonitUrl  = "monitor.api.qcloud.com/v2/index.php"
	Nonce     = 123
	Region    = "hk"
	namespace = "qce/cvm"
	VMID      = "ins-ao308lqc"
)

type ServerControllers struct {
	beego.Controller
}

var (
	SecretId  = beego.AppConfig.String("SecretId")
	SecretKey = beego.AppConfig.String("SecretKey")
)

type DateTime string
type PublicParm struct {
	Action    string
	SecretId  string
	Timestamp int64
	Nonce     int
	Region    string
	Signature string
}

type DescribeMetricsParm struct {
	PublicParm
	namespace  string
	metricName string
	mointUrl   string
}

type GetMonitorDataParm struct {
	//https://www.qcloud.com/document/api/248/4667
	PublicParm
	namespace  string
	metricName string
	period     int
	startTime  DateTime
	endTime    DateTime
	dimensions map[string]string
	mointUrl   string
}

func base64AndSha1(src, key string) (secure string) {
	//http://studygolang.com/articles/2667
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(src))
	secure = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return
}

func (c *DescribeMetricsParm) getUrl2() (url string) {
	now := time.Now().Unix()
	sortMap := map[string]interface{}{
		"Action":    c.Action,
		"Nonce":     c.Nonce,
		"Region":    c.Region,
		"SecretId":  c.SecretId,
		"Timestamp": now,
		"namespace": c.namespace,
	}
	var klist []string
	for k, _ := range sortMap {
		klist = append(klist, k)
	}
	sort.Strings(klist) //排序
	url = c.mointUrl + "?"
	for _, v := range klist {
		url += v
		url += "="
		if p, ok := sortMap[v].(string); ok {
			url += p + "&"
		}
		if p, ok := sortMap[v].(int64); ok {
			url += strconv.FormatInt(p, 10) + "&"
		}
		if p, ok := sortMap[v].(int); ok {
			url += strconv.Itoa(p) + "&"
		}
	}
	beego.Debug(url)
	methodAndUrl := Method + url
	beego.Debug(methodAndUrl)
	c.Signature = base64AndSha1(strings.TrimRight(methodAndUrl, "&"), SecretKey)
	url = "https://" + url + "Signature=" + c.Signature
	beego.Debug(url)
	return url
}

func (c *DescribeMetricsParm) getUrl() (url string) {
	now := time.Now().Unix()
	src := fmt.Sprintf("%s?Action=%s&Nonce=%d&Region=%s&SecretId=%s&Timestamp=%d&namespace=%s",
		c.mointUrl, c.Action, c.Nonce, c.Region, c.SecretId, now, c.namespace)
	methodAndUrl := Method + src
	beego.Debug(methodAndUrl)
	c.Signature = base64AndSha1(methodAndUrl, SecretKey)
	url = "https://" + src + "&Signature=" + c.Signature
	beego.Debug(url)
	return
}

//func createUrl(Action string) (str string) {
//	now := time.Now().Unix()
//	key := []byte(SecretKey)
//	mac := hmac.New(sha1.New, key)
//	url := fmt.Sprintf("%s?Action=%s&Nonce=%d&Region=%s&SecretId=%s&Timestamp=%d&namespace=%s",
//		MonitUrl, Action, Nonce, Region, SecretId, now, namespace)
//	methodAndUrl := Method + url
//	mac.Write([]byte(methodAndUrl))
//	encodeString := base64.StdEncoding.EncodeToString(mac.Sum(nil))
//
//	//
//	str = "https://" + url + "&Signature=" + encodeString
//	beego.Debug(str)
//	return
//}
func (c *ServerControllers) Get() {
	var result interface{}
	DescribeMetrics := new(DescribeMetricsParm)
	DescribeMetrics.mointUrl = MonitUrl
	DescribeMetrics.Action = "DescribeMetrics"
	DescribeMetrics.Nonce = Nonce
	DescribeMetrics.Region = Region
	DescribeMetrics.SecretId = SecretId
	DescribeMetrics.namespace = namespace

	url := DescribeMetrics.getUrl2()
	err := httplib.Get(url).ToJSON(&result)
	if err != nil {
		beego.Error(err)
	}
	beego.Info(result)
	c.Data["json"] = result
	c.ServeJSON()
}
