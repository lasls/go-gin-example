package dns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DnsPodClient DNSPod API客户端
type DnsPodClient struct {
	Token string
}

// DnsRecord DNS记录结构
type DnsRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Status    string `json:"status"`
	Weight    string `json:"weight,omitempty"`
	Line      string `json:"line,omitempty"`
	LineID    string `json:"line_id,omitempty"`
	Enabled   string `json:"enabled"`
	Remark    string `json:"remark"`
	UpdatedOn string `json:"updated_on"`
}

// DnsRecordListResponse DNS记录列表响应
type DnsRecordListResponse struct {
	Status   DnsStatus   `json:"status"`
	Info     DnsInfo     `json:"info"`
	Domain   DnsDomain   `json:"domain"`
	Records  []DnsRecord `json:"records"`
	InfoType string      `json:"info_type"`
}

// DnsStatus API状态
type DnsStatus struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedOn string `json:"created_on"`
}

// DnsInfo API信息
type DnsInfo struct {
	SubDomainItems int `json:"sub_domain_items"`
	RecordTotal    int `json:"record_total"`
	PageLimit      int `json:"page_limit"`
	Page           int `json:"page"`
}

// DnsDomain 域名信息
type DnsDomain struct {
	ID               int           `json:"id"` // 修复：ID字段应为int类型，不是string
	Name             string        `json:"name"`
	PunyCode         string        `json:"punycode"`
	Grade            string        `json:"grade"`
	Owner            string        `json:"owner"`
	Status           string        `json:"status"`
	GradeLevel       int           `json:"grade_level"`
	GradeTitle       string        `json:"grade_title"`
	IsVip            string        `json:"is_vip"`
	OwnerEmail       string        `json:"owner"`
	Records          string        `json:"records"`
	CreatedOn        string        `json:"created_on"`
	UpdatedOn        string        `json:"updated_on"`
	Server           string        `json:"server"`
	ExtStatus        string        `json:"ext_status"`
	Remark           string        `json:"remark"`
	TTL              string        `json:"ttl"`
	CnameSpeedup     string        `json:"cname_speedup"`
	SearchenginePush string        `json:"searchengine_push"`
	IsMark           string        `json:"is_mark"`
	GroupId          string        `json:"group_id"`
	SrcFlag          string        `json:"src_flag"`
	GradeNs          []string      `json:"grade_ns"`
	TagList          []interface{} `json:"tag_list"` // 可能为空数组
}

// NewDnsPodClient 创建DNSPod客户端
func NewDnsPodClient(token string) *DnsPodClient {
	return &DnsPodClient{
		Token: token,
	}
}

// makeRequest 发起HTTP请求
func (c *DnsPodClient) makeRequest(method, url string, params map[string]string) ([]byte, error) {
	client := &http.Client{}

	// 添加认证参数
	if params == nil {
		params = make(map[string]string)
	}
	params["login_token"] = c.Token
	params["format"] = "json"

	// 构建请求体
	var req *http.Request
	var err error

	if method == "GET" {
		// 构建查询字符串
		query := ""
		for k, v := range params {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", k, v)
		}
		url = fmt.Sprintf("%s?%s", url, query)
		req, err = http.NewRequest(method, url, nil)
	} else {
		// POST请求
		// 将参数转换为表单格式
		formData := ""
		for k, v := range params {
			if formData != "" {
				formData += "&"
			}
			formData += fmt.Sprintf("%s=%s", k, v)
		}
		req, err = http.NewRequest(method, url, strings.NewReader(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
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

	return body, nil
}

// GetDomainList 获取域名列表
func (c *DnsPodClient) GetDomainList() ([]DnsDomain, error) {
	url := "https://dnsapi.cn/Domain.List"
	params := map[string]string{}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status  DnsStatus   `json:"status"`
		Info    DnsInfo     `json:"info"`
		Domains []DnsDomain `json:"domains"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return nil, fmt.Errorf("API Error: %s", result.Status.Message)
	}

	return result.Domains, nil
}

// GetRecordList 获取DNS记录列表
func (c *DnsPodClient) GetRecordList(domain string, subDomain string) ([]DnsRecord, error) {
	url := "https://dnsapi.cn/Record.List"
	params := map[string]string{
		"domain": domain,
	}

	if subDomain != "" {
		params["sub_domain"] = subDomain
	}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var result DnsRecordListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return nil, fmt.Errorf("API Error: %v", result.Status.Message)
	}

	return result.Records, nil
}

// CreateRecord 创建DNS记录
func (c *DnsPodClient) CreateRecord(domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error) {
	url := "https://dnsapi.cn/Record.Create"
	params := map[string]string{
		"domain_id":   domainID,
		"sub_domain":  subDomain,
		"record_type": recordType,
		"value":       value,
		"record_line": recordLine,
	}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status DnsStatus `json:"status"`
		Record DnsRecord `json:"record"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return nil, fmt.Errorf("API Error: %s", result.Status.Message)
	}

	return &result.Record, nil
}

// UpdateRecord 更新DNS记录
func (c *DnsPodClient) UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error) {
	url := "https://dnsapi.cn/Record.Modify"
	params := map[string]string{
		"record_id":   recordID,
		"domain_id":   domainID,
		"sub_domain":  subDomain,
		"record_type": recordType,
		"value":       value,
		"record_line": recordLine,
	}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status DnsStatus `json:"status"`
		Record DnsRecord `json:"record"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return nil, fmt.Errorf("API Error: %s", result.Status.Message)
	}

	return &result.Record, nil
}

// DeleteRecord 删除DNS记录
func (c *DnsPodClient) DeleteRecord(recordID, domainID string) error {
	url := "https://dnsapi.cn/Record.Remove"
	params := map[string]string{
		"record_id": recordID,
		"domain_id": domainID,
	}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return err
	}

	var result struct {
		Status DnsStatus `json:"status"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return fmt.Errorf("API Error: %s", result.Status.Message)
	}

	return nil
}

// SetRecordStatus 设置记录状态
func (c *DnsPodClient) SetRecordStatus(recordID, domainID, status string) error {
	url := "https://dnsapi.cn/Record.Status"
	params := map[string]string{
		"record_id": recordID,
		"domain_id": domainID,
		"status":    status,
	}

	resp, err := c.makeRequest("POST", url, params)
	if err != nil {
		return err
	}

	var result struct {
		Status DnsStatus `json:"status"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("解析API响应失败: %v, 响应内容: %s", err, string(resp))
	}

	if result.Status.Code != "1" {
		return fmt.Errorf("API Error: %s", result.Status.Message)
	}

	return nil
}
