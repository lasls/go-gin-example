package dns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// AliyunDnsClient 阿里云DNS客户端
type AliyunDnsClient struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
}

// AliyunDnsRecord 阿里云DNS记录结构
type AliyunDnsRecord struct {
	DomainName string `json:"DomainName"`
	RecordId   string `json:"RecordId"`
	Rr         string `json:"Rr"`
	Type       string `json:"Type"`
	Value      string `json:"Value"`
	TTL        int64  `json:"TTL"`
	Status     string `json:"Status"`
	Locked     bool   `json:"Locked"`
	Remark     string `json:"Remark"`
	Line       string `json:"Line"`
	Priority   int64  `json:"Priority,omitempty"`
	Weight     int64  `json:"Weight,omitempty"`
}

// AliyunDnsRecordListResponse 阿里云DNS记录列表响应
type AliyunDnsRecordListResponse struct {
	RequestId     string            `json:"RequestId"`
	TotalCount    int               `json:"TotalCount"`
	PageNumber    int               `json:"PageNumber"`
	PageSize      int               `json:"PageSize"`
	DomainRecords []AliyunDnsRecord `json:"DomainRecords"`
}

// AliyunDnsRecordResponse 阿里云DNS记录操作响应
type AliyunDnsRecordResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
}

// AliyunDnsDomainListResponse 阿里云域名列表响应
type AliyunDnsDomainListResponse struct {
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
	Domains    struct {
		Domain []struct {
			DomainId   string `json:"DomainId"`
			DomainName string `json:"DomainName"`
			PunyCode   string `json:"PunyCode"`
			Remark     string `json:"Remark"`
		} `json:"Domain"`
	} `json:"Domains"`
}

// AliyunDnsRecordStatusResponse 阿里云DNS记录状态响应
type AliyunDnsRecordStatusResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
	Status    string `json:"Status"`
}

// NewAliyunDnsClient 创建阿里云DNS客户端
func NewAliyunDnsClient(accessKeyId, accessKeySecret, regionId string) *AliyunDnsClient {
	if regionId == "" {
		regionId = "cn-hangzhou" // 默认区域
	}
	return &AliyunDnsClient{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
	}
}

// sign 计算签名
func (c *AliyunDnsClient) sign(params map[string]string, accessKeySecret string) string {
	// 对参数进行排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建查询字符串
	var queryStr string
	for _, k := range keys {
		if queryStr != "" {
			queryStr += "&"
		}
		queryStr += url.QueryEscape(k) + "=" + url.QueryEscape(params[k])
	}

	// 构建待签名字符串
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(queryStr)

	// 计算HMAC-SHA1签名
	key := []byte(accessKeySecret + "&")
	h := hmac.New(sha1.New, key)
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// makeRequest 发起阿里云DNS API请求
func (c *AliyunDnsClient) makeRequest(action string, params map[string]string) ([]byte, error) {
	client := &http.Client{}

	// 设置公共参数
	publicParams := map[string]string{
		"Action":           action,
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"AccessKeyId":      c.AccessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"RegionId":         c.RegionId,
	}

	// 合并参数
	for k, v := range params {
		publicParams[k] = v
	}

	// 计算签名
	signature := c.sign(publicParams, c.AccessKeySecret)
	publicParams["Signature"] = signature

	// 构建URL
	baseURL := "https://alidns.aliyuncs.com/"

	// 对参数进行排序以构建查询字符串
	var keys []string
	for k := range publicParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryStr string
	for _, k := range keys {
		if queryStr != "" {
			queryStr += "&"
		}
		queryStr += url.QueryEscape(k) + "=" + url.QueryEscape(publicParams[k])
	}

	urlStr := baseURL + "?" + queryStr

	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查响应是否为HTML（通常表示错误页面）
	if len(body) > 0 && body[0] == '<' {
		return nil, fmt.Errorf("API返回错误页面，可能认证失败或请求参数错误: %s", string(body))
	}

	// 检查响应是否为错误信息
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetAliyunDomainList 获取阿里云域名列表
func (c *AliyunDnsClient) GetAliyunDomainList(pageNumber, pageSize int) ([]AliyunDnsRecord, error) {
	params := map[string]string{
		"PageNumber": fmt.Sprintf("%d", pageNumber),
		"PageSize":   fmt.Sprintf("%d", pageSize),
	}

	resp, err := c.makeRequest("DescribeDomains", params)
	if err != nil {
		return nil, err
	}

	var result AliyunDnsDomainListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	// 将域名信息转换为DNS记录格式
	var records []AliyunDnsRecord
	for _, domain := range result.Domains.Domain {
		records = append(records, AliyunDnsRecord{
			DomainName: domain.DomainName,
			RecordId:   domain.DomainId,
			Remark:     domain.Remark,
		})
	}

	return records, nil
}

// GetAliyunRecordList 获取阿里云DNS记录列表
func (c *AliyunDnsClient) GetAliyunRecordList(domainName string, rrKeyWord string) ([]AliyunDnsRecord, error) {
	params := map[string]string{
		"DomainName": domainName,
	}

	if rrKeyWord != "" {
		params["RrKeyWord"] = rrKeyWord
	}

	resp, err := c.makeRequest("DescribeDomainRecords", params)
	if err != nil {
		return nil, err
	}

	var result AliyunDnsRecordListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	return result.DomainRecords, nil
}

// CreateAliyunRecord 创建阿里云DNS记录
func (c *AliyunDnsClient) CreateAliyunRecord(domainName, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error) {
	params := map[string]string{
		"DomainName": domainName,
		"RR":         rr,
		"Type":       recordType,
		"Value":      value,
	}

	if ttl > 0 {
		params["TTL"] = fmt.Sprintf("%d", ttl)
	}

	resp, err := c.makeRequest("AddDomainRecord", params)
	if err != nil {
		return nil, err
	}

	var result AliyunDnsRecordResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	// 返回新创建的记录信息
	newRecord := &AliyunDnsRecord{
		RecordId:   result.RecordId,
		Rr:         rr,
		Type:       recordType,
		Value:      value,
		DomainName: domainName,
		TTL:        ttl,
		Status:     "ENABLE", // 新创建的记录默认是启用的
	}

	return newRecord, nil
}

// UpdateAliyunRecord 更新阿里云DNS记录
func (c *AliyunDnsClient) UpdateAliyunRecord(recordId, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error) {
	params := map[string]string{
		"RecordId": recordId,
		"RR":       rr,
		"Type":     recordType,
		"Value":    value,
	}

	if ttl > 0 {
		params["TTL"] = fmt.Sprintf("%d", ttl)
	}

	resp, err := c.makeRequest("UpdateDomainRecord", params)
	if err != nil {
		return nil, err
	}

	var result AliyunDnsRecordResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	// 返回更新后的记录信息
	updatedRecord := &AliyunDnsRecord{
		RecordId: result.RecordId,
		Rr:       rr,
		Type:     recordType,
		Value:    value,
		TTL:      ttl,
	}

	return updatedRecord, nil
}

// DeleteAliyunRecord 删除阿里云DNS记录
func (c *AliyunDnsClient) DeleteAliyunRecord(recordId string) error {
	params := map[string]string{
		"RecordId": recordId,
	}

	resp, err := c.makeRequest("DeleteDomainRecord", params)
	if err != nil {
		return err
	}

	var result AliyunDnsRecordResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	return nil
}

// SetAliyunRecordStatus 设置阿里云DNS记录状态
func (c *AliyunDnsClient) SetAliyunRecordStatus(recordId, status string) error {
	params := map[string]string{
		"RecordId": recordId,
		"Status":   status,
	}

	resp, err := c.makeRequest("SetDomainRecordStatus", params)
	if err != nil {
		return err
	}

	var result AliyunDnsRecordStatusResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	return nil
}
