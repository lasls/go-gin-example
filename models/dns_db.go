package models

import (
	"time"
)

// DnsDomain 域名信息模型
type DnsDomain struct {
	ID         int        `gorm:"primary_key" json:"id"`
	Name       string     `gorm:"column:name;size:255;not null" json:"name"`
	Provider   string     `gorm:"column:provider;size:50;not null" json:"provider"` // dns_pod, aliyun, volcengine
	DomainID   string     `gorm:"column:domain_id;size:100" json:"domain_id"`       // 云服务商的域名ID
	Status     string     `gorm:"column:status;size:20" json:"status"`              // active, inactive
	Grade      string     `gorm:"column:grade;size:50" json:"grade"`                // 域名等级
	Owner      string     `gorm:"column:owner;size:100" json:"owner"`               // 域名所有者
	Remark     string     `gorm:"column:remark;type:text" json:"remark"`            // 备注
	CreatedOn  time.Time  `json:"created_on"`
	ModifiedOn time.Time  `json:"modified_on"`
	DeletedOn  *time.Time `json:"deleted_on"`
}

// DnsRecord DNS解析记录模型
type DnsRecord struct {
	ID         int        `gorm:"primary_key" json:"id"`
	DomainID   int        `gorm:"column:domain_id;not null" json:"domain_id"`       // 关联域名ID
	Name       string     `gorm:"column:name;size:255;not null" json:"name"`        // 记录名称，如 www
	Type       string     `gorm:"column:type;size:10;not null" json:"type"`         // 记录类型，如 A, CNAME, MX
	Value      string     `gorm:"column:value;size:255;not null" json:"value"`      // 记录值，如IP地址
	Status     string     `gorm:"column:status;size:20" json:"status"`              // enable, disable
	Line       string     `gorm:"column:line;size:50" json:"line"`                  // 线路
	TTL        int        `gorm:"column:ttl;default:600" json:"ttl"`                // TTL值
	Remark     string     `gorm:"column:remark;type:text" json:"remark"`            // 备注
	Provider   string     `gorm:"column:provider;size:50;not null" json:"provider"` // dns_pod, aliyun, volcengine
	RemoteID   string     `gorm:"column:remote_id;size:100" json:"remote_id"`       // 云服务商的记录ID
	CreatedOn  time.Time  `json:"created_on"`
	ModifiedOn time.Time  `json:"modified_on"`
	DeletedOn  *time.Time `json:"deleted_on"`
}

// TableName 指定DnsDomain表名
func (DnsDomain) TableName() string {
	return "dns_domains"
}

// TableName 指定DnsRecord表名
func (DnsRecord) TableName() string {
	return "dns_records"
}

// AddDnsDomain 添加域名
func AddDnsDomain(domain *DnsDomain) error {
	if err := db.Create(domain).Error; err != nil {
		return err
	}
	return nil
}

// GetDnsDomainList 获取域名列表
func GetDnsDomainList(pageNum, pageSize int, maps interface{}) ([]DnsDomain, error) {
	var domains []DnsDomain
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

// GetDnsDomainTotal 获取域名总数
func GetDnsDomainTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&DnsDomain{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDnsDomain 更新域名
func UpdateDnsDomain(id int, data interface{}) error {
	if err := db.Model(&DnsDomain{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteDnsDomain 删除域名
func DeleteDnsDomain(id int) error {
	if err := db.Where("id = ?", id).Delete(&DnsDomain{}).Error; err != nil {
		return err
	}
	return nil
}

// ExistDnsDomainByID 检查域名是否存在
func ExistDnsDomainByID(id int) bool {
	var domain DnsDomain
	db.Select("id").Where("id = ?", id).First(&domain)
	return domain.ID > 0
}

// AddDnsRecord 添加DNS解析记录
func AddDnsRecord(record *DnsRecord) error {
	if err := db.Create(record).Error; err != nil {
		return err
	}
	return nil
}

// GetDnsRecordList 获取DNS解析记录列表
func GetDnsRecordList(pageNum, pageSize int, maps interface{}) ([]DnsRecord, error) {
	var records []DnsRecord
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// GetDnsRecordTotal 获取DNS解析记录总数
func GetDnsRecordTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&DnsRecord{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDnsRecord 更新DNS解析记录
func UpdateDnsRecord(id int, data interface{}) error {
	if err := db.Model(&DnsRecord{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteDnsRecord 删除DNS解析记录
func DeleteDnsRecord(id int) error {
	if err := db.Where("id = ?", id).Delete(&DnsRecord{}).Error; err != nil {
		return err
	}
	return nil
}

// ExistDnsRecordByID 检查DNS解析记录是否存在
func ExistDnsRecordByID(id int) bool {
	var record DnsRecord
	db.Select("id").Where("id = ?", id).First(&record)
	return record.ID > 0
}

// GetDnsRecordByDomainID 根据域名ID获取DNS解析记录
func GetDnsRecordByDomainID(domainID int) ([]DnsRecord, error) {
	var records []DnsRecord
	err := db.Where("domain_id = ?", domainID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// GetDnsDomainByName 根据域名名称获取域名信息
func GetDnsDomainByName(name string) (*DnsDomain, error) {
	var domain DnsDomain
	err := db.Where("name = ?", name).First(&domain).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}
