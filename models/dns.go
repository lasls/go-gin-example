package models

import (
	"context"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/dns"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DnsService DNS服务结构
type DnsService struct {
	Manager *dns.DnsManager
}

// DnsDomain DNS域名结构
type DnsDomain struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Provider   string             `bson:"provider" json:"provider"`
	DomainID   string             `bson:"domain_id" json:"domain_id"`
	Status     string             `bson:"status" json:"status"`
	Grade      string             `bson:"grade" json:"grade"`
	Owner      string             `bson:"owner" json:"owner"`
	Remark     string             `bson:"remark" json:"remark"`
	CreatedOn  time.Time          `bson:"created_on" json:"created_on"`
	ModifiedOn time.Time          `bson:"modified_on" json:"modified_on"`
}

// DnsRecord DNS记录结构
type DnsRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DomainID   string             `bson:"domain_id" json:"domain_id"`
	Name       string             `bson:"name" json:"name"`
	Type       string             `bson:"type" json:"type"`
	Value      string             `bson:"value" json:"value"`
	Status     string             `bson:"status" json:"status"`
	Line       string             `bson:"line" json:"line"`
	TTL        int64              `bson:"ttl" json:"ttl"`
	Remark     string             `bson:"remark" json:"remark"`
	Provider   string             `bson:"provider" json:"provider"`
	RemoteID   string             `bson:"remote_id" json:"remote_id"`
	CreatedOn  time.Time          `bson:"created_on" json:"created_on"`
	ModifiedOn time.Time          `bson:"modified_on" json:"modified_on"`
}

const (
	CollectionDnsDomain = "dns_domains"
	CollectionDnsRecord = "dns_records"
)

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

// AddDnsDomain 添加域名
func AddDnsDomain(domain *DnsDomain) error {
	ctx := context.Background()
	domain.CreatedOn = time.Now()
	domain.ModifiedOn = time.Now()
	_, err := db.Collection(CollectionDnsDomain).InsertOne(ctx, domain)
	return err
}

// GetDnsDomainList 获取域名列表
func GetDnsDomainList(pageNum, pageSize int, maps bson.M) ([]DnsDomain, error) {
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetSkip(int64(pageNum))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := db.Collection(CollectionDnsDomain).Find(ctx, maps, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var domains []DnsDomain
	if err = cursor.All(ctx, &domains); err != nil {
		return nil, err
	}
	return domains, nil
}

// GetDnsDomainTotal 获取域名总数
func GetDnsDomainTotal(maps bson.M) (int64, error) {
	ctx := context.Background()
	count, err := db.Collection(CollectionDnsDomain).CountDocuments(ctx, maps)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDnsDomain 更新域名
func UpdateDnsDomain(id string, data bson.M) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	data["modified_on"] = time.Now()
	_, err = db.Collection(CollectionDnsDomain).UpdateByID(ctx, objectId, bson.M{"$set": data})
	return err
}

// DeleteDnsDomain 删除域名
func DeleteDnsDomain(id string) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Collection(CollectionDnsDomain).DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}

// ExistDnsDomainByID 检查域名是否存在
func ExistDnsDomainByID(id string) bool {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	count, err := db.Collection(CollectionDnsDomain).CountDocuments(ctx, bson.M{"_id": objectId})
	if err != nil {
		return false
	}
	return count > 0
}

// AddDnsRecord 添加DNS记录
func AddDnsRecord(record *DnsRecord) error {
	ctx := context.Background()
	record.CreatedOn = time.Now()
	record.ModifiedOn = time.Now()
	_, err := db.Collection(CollectionDnsRecord).InsertOne(ctx, record)
	return err
}

// GetDnsRecordList 获取DNS记录列表
func GetDnsRecordList(pageNum, pageSize int, maps bson.M) ([]DnsRecord, error) {
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetSkip(int64(pageNum))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := db.Collection(CollectionDnsRecord).Find(ctx, maps, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []DnsRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// GetDnsRecordTotal 获取DNS记录总数
func GetDnsRecordTotal(maps bson.M) (int64, error) {
	ctx := context.Background()
	count, err := db.Collection(CollectionDnsRecord).CountDocuments(ctx, maps)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDnsRecord 更新DNS记录
func UpdateDnsRecord(id string, data bson.M) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	data["modified_on"] = time.Now()
	_, err = db.Collection(CollectionDnsRecord).UpdateByID(ctx, objectId, bson.M{"$set": data})
	return err
}

// DeleteDnsRecord 删除DNS记录
func DeleteDnsRecord(id string) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Collection(CollectionDnsRecord).DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}

// ExistDnsRecordByID 检查DNS记录是否存在
func ExistDnsRecordByID(id string) bool {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	count, err := db.Collection(CollectionDnsRecord).CountDocuments(ctx, bson.M{"_id": objectId})
	if err != nil {
		return false
	}
	return count > 0
}

// GetDnsRecordByDomainID 根据域名ID获取DNS记录
func GetDnsRecordByDomainID(domainID string) ([]DnsRecord, error) {
	ctx := context.Background()
	cursor, err := db.Collection(CollectionDnsRecord).Find(ctx, bson.M{"domain_id": domainID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []DnsRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}
