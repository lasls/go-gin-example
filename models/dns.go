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
		setting.VolcengineAccessKeyId,
		setting.VolcengineAccessKeySecret,
		setting.VolcengineRegionId,
	)
	return &DnsService{
		Manager: manager,
	}
}

// GetDomainList 获取域名列表
func (s *DnsService) GetDomainList(provider string) ([]interface{}, error) {
	strategy := GetStrategy(provider)
	return strategy.GetDomainList(s.Manager)
}

// GetRecordList 获取DNS记录列表
func (s *DnsService) GetRecordList(domain, subDomain string, provider string) (interface{}, error) {
	strategy := GetStrategy(provider)
	return strategy.GetRecordList(s.Manager, domain, subDomain)
}

// CreateRecord 创建DNS记录
func (s *DnsService) CreateRecord(domainID, subDomain, recordType, value, recordLine string, provider string) (interface{}, error) {
	// 默认TTL为600
	ttl := int64(600)
	strategy := GetStrategy(provider)
	return strategy.CreateRecord(s.Manager, domainID, subDomain, recordType, value, recordLine, ttl)
}

// UpdateRecord 更新DNS记录
func (s *DnsService) UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine string, provider string) (interface{}, error) {
	// 默认TTL为600
	ttl := int64(600)
	strategy := GetStrategy(provider)
	return strategy.UpdateRecord(s.Manager, recordID, domainID, subDomain, recordType, value, recordLine, ttl)
}

// DeleteRecord 删除DNS记录
func (s *DnsService) DeleteRecord(recordID, domainID string, provider string) error {
	strategy := GetStrategy(provider)
	return strategy.DeleteRecord(s.Manager, recordID, domainID)
}

// SetRecordStatus 设置记录状态
func (s *DnsService) SetRecordStatus(recordID, domainID, status string, provider string) error {
	strategy := GetStrategy(provider)
	return strategy.SetRecordStatus(s.Manager, recordID, domainID, status)
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
