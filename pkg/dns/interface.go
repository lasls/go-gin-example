package dns

import "fmt"

// DnsProvider DNS服务提供商接口
type DnsProvider interface {
	GetDomainList() ([]DnsDomain, error)
	GetRecordList(domain string, subDomain string) ([]DnsRecord, error)
	CreateRecord(domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error)
	UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error)
	DeleteRecord(recordID, domainID string) error
	SetRecordStatus(recordID, domainID, status string) error
}

// AliyunDnsProvider 阿里云DNS服务提供商接口
type AliyunDnsProvider interface {
	GetAliyunDomainList(pageNumber, pageSize int) ([]AliyunDnsRecord, error)
	GetAliyunRecordList(domainName string, rrKeyWord string) ([]AliyunDnsRecord, error)
	CreateAliyunRecord(domainName, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error)
	UpdateAliyunRecord(recordId, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error)
	DeleteAliyunRecord(recordId string) error
	SetAliyunRecordStatus(recordId, status string) error
}

// DnsManager 统一DNS管理器
type DnsManager struct {
	dnsPodClient    *DnsPodClient
	aliyunDnsClient *AliyunDnsClient
}

// NewDnsManager 创建DNS管理器
func NewDnsManager(dnsPodToken, aliyunAccessKeyId, aliyunAccessKeySecret, aliyunRegionId string) *DnsManager {
	manager := &DnsManager{}

	if dnsPodToken != "" {
		manager.dnsPodClient = NewDnsPodClient(dnsPodToken)
	}

	if aliyunAccessKeyId != "" && aliyunAccessKeySecret != "" {
		manager.aliyunDnsClient = NewAliyunDnsClient(aliyunAccessKeyId, aliyunAccessKeySecret, aliyunRegionId)
	}

	return manager
}

// 使用DNSPod
func (m *DnsManager) UseDnsPod() bool {
	return m.dnsPodClient != nil
}

// 使用阿里云DNS
func (m *DnsManager) UseAliyunDns() bool {
	return m.aliyunDnsClient != nil
}

// GetDnsPodDomainList 获取DNSPod域名列表
func (m *DnsManager) GetDnsPodDomainList() ([]DnsDomain, error) {
	if m.dnsPodClient == nil {
		return nil, fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.GetDomainList()
}

// GetDnsPodRecordList 获取DNSPod记录列表
func (m *DnsManager) GetDnsPodRecordList(domain, subDomain string) ([]DnsRecord, error) {
	if m.dnsPodClient == nil {
		return nil, fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.GetRecordList(domain, subDomain)
}

// CreateDnsPodRecord 创建DNSPod记录
func (m *DnsManager) CreateDnsPodRecord(domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error) {
	if m.dnsPodClient == nil {
		return nil, fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.CreateRecord(domainID, subDomain, recordType, value, recordLine)
}

// UpdateDnsPodRecord 更新DNSPod记录
func (m *DnsManager) UpdateDnsPodRecord(recordID, domainID, subDomain, recordType, value, recordLine string) (*DnsRecord, error) {
	if m.dnsPodClient == nil {
		return nil, fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine)
}

// DeleteDnsPodRecord 删除DNSPod记录
func (m *DnsManager) DeleteDnsPodRecord(recordID, domainID string) error {
	if m.dnsPodClient == nil {
		return fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.DeleteRecord(recordID, domainID)
}

// SetDnsPodRecordStatus 设置DNSPod记录状态
func (m *DnsManager) SetDnsPodRecordStatus(recordID, domainID, status string) error {
	if m.dnsPodClient == nil {
		return fmt.Errorf("DNSPod客户端未初始化")
	}
	return m.dnsPodClient.SetRecordStatus(recordID, domainID, status)
}

// GetAliyunDomainList 获取阿里云域名列表
func (m *DnsManager) GetAliyunDomainList(pageNumber, pageSize int) ([]AliyunDnsRecord, error) {
	if m.aliyunDnsClient == nil {
		return nil, fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.GetAliyunDomainList(pageNumber, pageSize)
}

// GetAliyunRecordList 获取阿里云记录列表
func (m *DnsManager) GetAliyunRecordList(domainName, rrKeyWord string) ([]AliyunDnsRecord, error) {
	if m.aliyunDnsClient == nil {
		return nil, fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.GetAliyunRecordList(domainName, rrKeyWord)
}

// CreateAliyunRecord 创建阿里云记录
func (m *DnsManager) CreateAliyunRecord(domainName, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error) {
	if m.aliyunDnsClient == nil {
		return nil, fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.CreateAliyunRecord(domainName, rr, recordType, value, ttl)
}

// UpdateAliyunRecord 更新阿里云记录
func (m *DnsManager) UpdateAliyunRecord(recordId, rr, recordType, value string, ttl int64) (*AliyunDnsRecord, error) {
	if m.aliyunDnsClient == nil {
		return nil, fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.UpdateAliyunRecord(recordId, rr, recordType, value, ttl)
}

// DeleteAliyunRecord 删除阿里云记录
func (m *DnsManager) DeleteAliyunRecord(recordId string) error {
	if m.aliyunDnsClient == nil {
		return fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.DeleteAliyunRecord(recordId)
}

// SetAliyunRecordStatus 设置阿里云记录状态
func (m *DnsManager) SetAliyunRecordStatus(recordId, status string) error {
	if m.aliyunDnsClient == nil {
		return fmt.Errorf("阿里云DNS客户端未初始化")
	}
	return m.aliyunDnsClient.SetAliyunRecordStatus(recordId, status)
}
