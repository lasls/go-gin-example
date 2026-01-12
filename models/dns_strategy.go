package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/dns"
)

// DnsStrategy 策略接口
type DnsStrategy interface {
	GetDomainList(manager *dns.DnsManager) ([]interface{}, error)
	GetRecordList(manager *dns.DnsManager, domain, subDomain string) (interface{}, error)
	CreateRecord(manager *dns.DnsManager, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error)
	UpdateRecord(manager *dns.DnsManager, recordID, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error)
	DeleteRecord(manager *dns.DnsManager, recordID, domainID string) error
	SetRecordStatus(manager *dns.DnsManager, recordID, domainID, status string) error
}

// DnsPodStrategy DNSPod策略实现
type DnsPodStrategy struct{}

func (s *DnsPodStrategy) GetDomainList(manager *dns.DnsManager) ([]interface{}, error) {
	if !manager.UseDnsPod() {
		return []interface{}{}, nil
	}

	domains, err := manager.GetDnsPodDomainList()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(domains))
	for i, domain := range domains {
		result[i] = map[string]interface{}{
			"Domain": domain,
			"Status": domain.Status,
		}
	}
	return result, nil
}

func (s *DnsPodStrategy) GetRecordList(manager *dns.DnsManager, domain, subDomain string) (interface{}, error) {
	if manager.UseDnsPod() {
		return manager.GetDnsPodRecordList(domain, subDomain)
	}
	return nil, nil
}

func (s *DnsPodStrategy) CreateRecord(manager *dns.DnsManager, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseDnsPod() {
		return manager.CreateDnsPodRecord(domainID, subDomain, recordType, value, recordLine)
	}
	return nil, nil
}

func (s *DnsPodStrategy) UpdateRecord(manager *dns.DnsManager, recordID, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseDnsPod() {
		return manager.UpdateDnsPodRecord(recordID, domainID, subDomain, recordType, value, recordLine)
	}
	return nil, nil
}

func (s *DnsPodStrategy) DeleteRecord(manager *dns.DnsManager, recordID, domainID string) error {
	if manager.UseDnsPod() {
		return manager.DeleteDnsPodRecord(recordID, domainID)
	}
	return nil
}

func (s *DnsPodStrategy) SetRecordStatus(manager *dns.DnsManager, recordID, domainID, status string) error {
	if manager.UseDnsPod() {
		return manager.SetDnsPodRecordStatus(recordID, domainID, status)
	}
	return nil
}

// AliyunStrategy 阿里云策略实现
type AliyunStrategy struct{}

func (s *AliyunStrategy) GetDomainList(manager *dns.DnsManager) ([]interface{}, error) {
	if !manager.UseAliyunDns() {
		return []interface{}{}, nil
	}

	domains, err := manager.GetAliyunDomainList(1, 20)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(domains))
	for i, domain := range domains {
		result[i] = map[string]interface{}{
			"Domain": domain,
			"Status": domain.Status,
		}
	}
	return result, nil
}

func (s *AliyunStrategy) GetRecordList(manager *dns.DnsManager, domain, subDomain string) (interface{}, error) {
	if manager.UseAliyunDns() {
		return manager.GetAliyunRecordList(domain, subDomain)
	}
	return nil, nil
}

func (s *AliyunStrategy) CreateRecord(manager *dns.DnsManager, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseAliyunDns() {
		return manager.CreateAliyunRecord(domainID, subDomain, recordType, value, ttl)
	}
	return nil, nil
}

func (s *AliyunStrategy) UpdateRecord(manager *dns.DnsManager, recordID, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseAliyunDns() {
		return manager.UpdateAliyunRecord(recordID, subDomain, recordType, value, ttl)
	}
	return nil, nil
}

func (s *AliyunStrategy) DeleteRecord(manager *dns.DnsManager, recordID, domainID string) error {
	if manager.UseAliyunDns() {
		return manager.DeleteAliyunRecord(recordID)
	}
	return nil
}

func (s *AliyunStrategy) SetRecordStatus(manager *dns.DnsManager, recordID, domainID, status string) error {
	if manager.UseAliyunDns() {
		return manager.SetAliyunRecordStatus(recordID, status)
	}
	return nil
}

// VolcengineStrategy 火山引擎策略实现
type VolcengineStrategy struct{}

func (s *VolcengineStrategy) GetDomainList(manager *dns.DnsManager) ([]interface{}, error) {
	if !manager.UseVolcengineDns() {
		return []interface{}{}, nil
	}

	domains, err := manager.GetVolcengineDomainList("")
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(domains))
	for i, domain := range domains {
		result[i] = map[string]interface{}{
			"Domain": domain,
			"Status": domain.Status,
		}
	}
	return result, nil
}

func (s *VolcengineStrategy) GetRecordList(manager *dns.DnsManager, domain, subDomain string) (interface{}, error) {
	if manager.UseVolcengineDns() {
		return manager.GetVolcengineRecordList(domain, subDomain)
	}
	return nil, nil
}

func (s *VolcengineStrategy) CreateRecord(manager *dns.DnsManager, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseVolcengineDns() {
		return manager.CreateVolcengineRecord(domainID, subDomain, recordType, value, ttl)
	}
	return nil, nil
}

func (s *VolcengineStrategy) UpdateRecord(manager *dns.DnsManager, recordID, domainID, subDomain, recordType, value, recordLine string, ttl int64) (interface{}, error) {
	if manager.UseVolcengineDns() {
		return manager.UpdateVolcengineRecord(recordID, subDomain, recordType, value, ttl)
	}
	return nil, nil
}

func (s *VolcengineStrategy) DeleteRecord(manager *dns.DnsManager, recordID, domainID string) error {
	if manager.UseVolcengineDns() {
		return manager.DeleteVolcengineRecord(recordID)
	}
	return nil
}

func (s *VolcengineStrategy) SetRecordStatus(manager *dns.DnsManager, recordID, domainID, status string) error {
	if manager.UseVolcengineDns() {
		return manager.SetVolcengineRecordStatus(recordID, status)
	}
	return nil
}

// StrategyMap 策略映射
var StrategyMap = map[string]DnsStrategy{
	"dns_pod":    &DnsPodStrategy{},
	"aliyun":     &AliyunStrategy{},
	"volcengine": &VolcengineStrategy{},
}

// GetDefaultStrategy 获取默认策略
func GetDefaultStrategy() DnsStrategy {
	return &DnsPodStrategy{}
}

// GetStrategy 根据provider获取策略
func GetStrategy(provider string) DnsStrategy {
	if strategy, exists := StrategyMap[provider]; exists {
		return strategy
	}
	return GetDefaultStrategy()
}
