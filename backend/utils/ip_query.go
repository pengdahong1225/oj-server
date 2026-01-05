package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// http://ip-api.com/json/110.242.68.66?fields=status,message,continent,country,region,regionName,city,district,zip,lat,lon,timezone,query&lang=zh-CN

const (
	baseUrl = "http://ip-api.com/json"
	fields  = "status,message,continent,country,region,regionName,city,district,zip,lat,lon,timezone,query"
	lang    = "zh-CN"
)

type IPInfoResp struct {
	Status     string  `json:"status"`     // success or fail
	Message    string  `json:"message"`    // included only when status is fail, Can be one of the following: private range, reserved range, invalid query
	Continent  string  `json:"continent"`  // 洲
	Country    string  `json:"country"`    // 国家
	Region     string  `json:"region"`     // 省份
	RegionName string  `json:"regionName"` // 省份名称
	City       string  `json:"city"`       // 城市
	District   string  `json:"district"`   // 区
	Zip        string  `json:"zip"`        // 邮编
	Lat        float64 `json:"lat"`        // 纬度
	Lon        float64 `json:"lon"`        // 经度
	TimeZone   string  `json:"timezone"`   // 时区
	Query      string  `json:"query"`      // 查询的ip
}

func QueryIpGeolocation(ip string) (*IPInfoResp, error) {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 构建请求URL
	url := fmt.Sprintf("%s/%s?fileds=%s&lang=%s", baseUrl, ip, fields, lang)

	// 发送GET请求
	resp, err := client.Get(url)
	if err != nil {
		return &IPInfoResp{
			RegionName: "未知地区",
		}, err
	}
	defer resp.Body.Close()

	// 解析响应体为JSON
	ipInfo := &IPInfoResp{}
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return &IPInfoResp{
			RegionName: "未知地区",
		}, err
	}
	if ipInfo.Status != "success" {
		return &IPInfoResp{
			RegionName: "未知地区",
		}, fmt.Errorf("ipInfo.Status != success")
	}

	return ipInfo, nil
}
