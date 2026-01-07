package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/dns"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// DnsService DNS服务结构
type DnsService struct {
	Manager *dns.DnsManager
}

// NewDnsService 创建DNS服务实例
func NewDnsService() *DnsService {
	manager := dns.NewDnsManager(
		setting.DnsPodToken,
		setting.AliyunAccessKeyId,
		setting.AliyunAccessKeySecret,
		setting.AliyunRegionId,
	)
	return &DnsService{
		Manager: manager,
	}
}

// GetDomainList 获取域名列表
func (s *DnsService) GetDomainList() ([]dns.DnsDomain, error) {
	// 如果配置了DNSPod，优先使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.GetDnsPodDomainList()
	}
	// 否则尝试使用阿里云DNS
	return nil, nil // 阿里云域名列表需要单独处理
}

// GetRecordList 获取DNS记录列表
func (s *DnsService) GetRecordList(domain, subDomain string, provider string) (interface{}, error) {
	if provider == "aliyun" && s.Manager.UseAliyunDns() {
		return s.Manager.GetAliyunRecordList(domain, subDomain)
	}
	// 默认使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.GetDnsPodRecordList(domain, subDomain)
	}
	return nil, nil
}

// CreateRecord 创建DNS记录
func (s *DnsService) CreateRecord(domainID, subDomain, recordType, value, recordLine string, provider string) (interface{}, error) {
	if provider == "aliyun" && s.Manager.UseAliyunDns() {
		// 对于阿里云，domainID是域名名称，subDomain是RR值
		return s.Manager.CreateAliyunRecord(domainID, subDomain, recordType, value, 600) // 默认TTL为600
	}
	// 默认使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.CreateDnsPodRecord(domainID, subDomain, recordType, value, recordLine)
	}
	return nil, nil
}

// UpdateRecord 更新DNS记录
func (s *DnsService) UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine string, provider string) (interface{}, error) {
	if provider == "aliyun" && s.Manager.UseAliyunDns() {
		return s.Manager.UpdateAliyunRecord(recordID, subDomain, recordType, value, 600) // 默认TTL为600
	}
	// 默认使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.UpdateDnsPodRecord(recordID, domainID, subDomain, recordType, value, recordLine)
	}
	return nil, nil
}

// DeleteRecord 删除DNS记录
func (s *DnsService) DeleteRecord(recordID, domainID string, provider string) error {
	if provider == "aliyun" && s.Manager.UseAliyunDns() {
		return s.Manager.DeleteAliyunRecord(recordID)
	}
	// 默认使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.DeleteDnsPodRecord(recordID, domainID)
	}
	return nil
}

// SetRecordStatus 设置记录状态
func (s *DnsService) SetRecordStatus(recordID, domainID, status string, provider string) error {
	if provider == "aliyun" && s.Manager.UseAliyunDns() {
		return s.Manager.SetAliyunRecordStatus(recordID, status)
	}
	// 默认使用DNSPod
	if s.Manager.UseDnsPod() {
		return s.Manager.SetDnsPodRecordStatus(recordID, domainID, status)
	}
	return nil
}

// GetAliyunRecordList 获取阿里云DNS记录列表
func (s *DnsService) GetAliyunRecordList(domainName, rrKeyWord string) ([]dns.AliyunDnsRecord, error) {
	if s.Manager.UseAliyunDns() {
		return s.Manager.GetAliyunRecordList(domainName, rrKeyWord)
	}
	return nil, nil
}

// CreateAliyunRecord 创建阿里云DNS记录
func (s *DnsService) CreateAliyunRecord(domainName, rr, recordType, value string, ttl int64) (*dns.AliyunDnsRecord, error) {
	if s.Manager.UseAliyunDns() {
		return s.Manager.CreateAliyunRecord(domainName, rr, recordType, value, ttl)
	}
	return nil, nil
}

// UpdateAliyunRecord 更新阿里云DNS记录
func (s *DnsService) UpdateAliyunRecord(recordId, rr, recordType, value string, ttl int64) (*dns.AliyunDnsRecord, error) {
	if s.Manager.UseAliyunDns() {
		return s.Manager.UpdateAliyunRecord(recordId, rr, recordType, value, ttl)
	}
	return nil, nil
}

// DeleteAliyunRecord 删除阿里云DNS记录
func (s *DnsService) DeleteAliyunRecord(recordId string) error {
	if s.Manager.UseAliyunDns() {
		return s.Manager.DeleteAliyunRecord(recordId)
	}
	return nil
}

// SetAliyunRecordStatus 设置阿里云DNS记录状态
func (s *DnsService) SetAliyunRecordStatus(recordId, status string) error {
	if s.Manager.UseAliyunDns() {
		return s.Manager.SetAliyunRecordStatus(recordId, status)
	}
	return nil
}
