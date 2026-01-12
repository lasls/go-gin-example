package dns

import (
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/dns"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

// VolcengineDnsClient 火山引擎DNS客户端
type VolcengineDnsClient struct {
	Client *dns.DNS
}

// VolcengineDnsRecord 火山引擎DNS记录结构
type VolcengineDnsRecord struct {
	RecordID      string `json:"RecordID"`
	SubDomainName string `json:"SubDomainName"`
	RecordType    string `json:"RecordType"`
	RecordValue   string `json:"RecordValue"`
	Status        string `json:"Status"`
	TTL           int64  `json:"TTL"`
	Remark        string `json:"Remark"`
	Line          string `json:"Line"`
	CreateTime    string `json:"CreateTime"`
	UpdateTime    string `json:"UpdateTime"`
}

// NewVolcengineDnsClient 创建火山引擎DNS客户端
func NewVolcengineDnsClient(accessKeyId, secretAccessKey, region string) *VolcengineDnsClient {
	// 创建凭证
	creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")

	// 创建会话
	sess := session.Must(session.NewSession(&volcengine.Config{
		Credentials: creds,
		Region:      volcengine.String(region),
	}))

	// 创建DNS客户端
	client := dns.New(sess)

	return &VolcengineDnsClient{
		Client: client,
	}
}

// GetVolcengineDomainList 获取火山引擎域名列表
func (c *VolcengineDnsClient) GetVolcengineDomainList(keyword string) ([]VolcengineDomain, error) {
	input := &dns.ListZonesInput{
		PageSize:   volcengine.Int32(100),
		PageNumber: volcengine.Int32(1),
	}
	if keyword != "" {
		input.Key = volcengine.String(keyword)
	}

	output, err := c.Client.ListZones(input)
	if err != nil {
		return nil, fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return nil, fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	// 转换响应数据
	var domains []VolcengineDomain
	for _, zone := range output.Zones {
		domains = append(domains, VolcengineDomain{
			Domain:       *zone.ZoneName,
			Status:       *zone.DnsSecurity, // 使用DnsSecurity字段作为状态
			VerifyStatus: *zone.DnsSecurity,
			CreatedAt:    *zone.CreatedAt,
			AccountId:    *zone.InstanceID,
		})
	}

	return domains, nil
}

// GetVolcengineRecordList 获取火山引擎DNS记录列表
func (c *VolcengineDnsClient) GetVolcengineRecordList(zoneId, rrKeyWord string) ([]VolcengineDnsRecord, error) {
	// zoneId在火山引擎DNS API中是ZID
	zid := int64(0)
	fmt.Sscanf(zoneId, "%d", &zid)

	input := &dns.ListRecordsInput{
		ZID:        volcengine.Int64(zid),
		PageSize:   volcengine.Int32(100),
		PageNumber: volcengine.Int32(1),
	}
	if rrKeyWord != "" {
		input.Host = volcengine.String(rrKeyWord)
	}

	output, err := c.Client.ListRecords(input)
	if err != nil {
		return nil, fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return nil, fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	// 转换响应数据
	var records []VolcengineDnsRecord
	for _, record := range output.Records {
		// 火山引擎的TTL是int32类型，转换为int64
		ttl := int64(*record.TTL)

		// 火山引擎的RecordID是string类型，使用它作为RecordID
		recordId := ""
		if record.RecordID != nil {
			recordId = *record.RecordID
		}

		// 状态从Enable字段获取
		status := "DISABLED"
		if record.Enable != nil && *record.Enable {
			status = "ENABLED"
		}

		records = append(records, VolcengineDnsRecord{
			RecordID:      recordId,
			SubDomainName: *record.Host,
			RecordType:    *record.Type,
			RecordValue:   *record.Value,
			Status:        status,
			TTL:           ttl,
			Remark:        *record.Remark,
			Line:          *record.Line,
			CreateTime:    *record.CreatedAt,
			UpdateTime:    *record.UpdatedAt,
		})
	}

	return records, nil
}

// CreateVolcengineRecord 创建火山引擎DNS记录
func (c *VolcengineDnsClient) CreateVolcengineRecord(zoneId, rr, recordType, value string, ttl int64) (*VolcengineDnsRecord, error) {
	// zoneId在火山引擎DNS API中是ZID
	zid := int64(0)
	fmt.Sscanf(zoneId, "%d", &zid)

	input := &dns.CreateRecordInput{
		ZID:    volcengine.Int64(zid),
		Host:   volcengine.String(rr),
		Type:   volcengine.String(recordType),
		Value:  volcengine.String(value),
		TTL:    volcengine.Int32(int32(ttl)),
		Line:   volcengine.String("default"), // 默认线路
		Remark: volcengine.String("Created by Go-Gin Application"),
	}

	output, err := c.Client.CreateRecord(input)
	if err != nil {
		return nil, fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return nil, fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	// 转换响应数据
	status := "DISABLED"
	if output.Enable != nil && *output.Enable {
		status = "ENABLED"
	}

	record := &VolcengineDnsRecord{
		RecordID:      *output.RecordID,
		SubDomainName: *output.Host,
		RecordType:    *output.Type,
		RecordValue:   *output.Value,
		Status:        status,
		TTL:           int64(*output.TTL),
		Remark:        *output.Remark,
		Line:          *output.Line,
		CreateTime:    *output.CreatedAt,
		UpdateTime:    *output.UpdatedAt,
	}

	return record, nil
}

// UpdateVolcengineRecord 更新火山引擎DNS记录
func (c *VolcengineDnsClient) UpdateVolcengineRecord(recordId, rr, recordType, value string, ttl int64) (*VolcengineDnsRecord, error) {
	// recordId在火山引擎DNS API中是RecordID
	input := &dns.UpdateRecordInput{
		RecordID: volcengine.String(recordId),
		Host:     volcengine.String(rr),
		Line:     volcengine.String("default"), // Line字段是必需的
		Type:     volcengine.String(recordType),
		Value:    volcengine.String(value),
		TTL:      volcengine.Int32(int32(ttl)),
		Remark:   volcengine.String("Updated by Go-Gin Application"),
	}

	output, err := c.Client.UpdateRecord(input)
	if err != nil {
		return nil, fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return nil, fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	// 转换响应数据
	status := "DISABLED"
	if output.Enable != nil && *output.Enable {
		status = "ENABLED"
	}

	record := &VolcengineDnsRecord{
		RecordID:      *output.RecordID,
		SubDomainName: *output.Host,
		RecordType:    *output.Type,
		RecordValue:   *output.Value,
		Status:        status,
		TTL:           int64(*output.TTL),
		Remark:        *output.Remark,
		Line:          *output.Line,
		CreateTime:    *output.CreatedAt,
		UpdateTime:    *output.UpdatedAt,
	}

	return record, nil
}

// DeleteVolcengineRecord 删除火山引擎DNS记录
func (c *VolcengineDnsClient) DeleteVolcengineRecord(recordId string) error {
	input := &dns.DeleteRecordInput{
		RecordID: volcengine.String(recordId),
	}

	output, err := c.Client.DeleteRecord(input)
	if err != nil {
		return fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	return nil
}

// SetVolcengineRecordStatus 设置火山引擎DNS记录状态
func (c *VolcengineDnsClient) SetVolcengineRecordStatus(recordId, status string) error {
	// 火山引擎API中状态是布尔值，"ENABLED"对应true，"DISABLED"对应false
	enable := status == "ENABLED"

	input := &dns.UpdateRecordStatusInput{
		RecordID: volcengine.String(recordId),
		Enable:   volcengine.Bool(enable),
	}

	output, err := c.Client.UpdateRecordStatus(input)
	if err != nil {
		return fmt.Errorf("API Error: %v", err)
	}

	// 检查响应中是否有错误
	if output.Metadata.Error != nil {
		return fmt.Errorf("API Error: %s - %s", output.Metadata.Error.Code, output.Metadata.Error.Message)
	}

	return nil
}

// VolcengineDomain 火山引擎域名结构
type VolcengineDomain struct {
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	VerifyStatus string `json:"verify_status"`
	CreatedAt    string `json:"created_at"`
	AccountId    string `json:"account_id"`
}
