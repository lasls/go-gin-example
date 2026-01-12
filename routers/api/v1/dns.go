package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/dns"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 获取域名列表
func GetDomains(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 检查使用哪个提供商
	provider := c.Query("provider")

	// 使用通用域名列表方法
	domains, err := dnsService.GetDomainList(provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "success",
		"data": domains,
	})
}

// 获取DNS记录列表
func GetDnsRecords(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	provider := c.Query("provider")
	domain := c.Query("domain")
	if domain == "" {
		domain = setting.DomainName // 使用配置中的默认域名
	}
	subDomain := c.Query("sub_domain")

	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "域名不能为空",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 使用通用记录列表方法
	records, err := dnsService.GetRecordList(domain, subDomain, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "success",
		"data": records,
	})
}

// 创建DNS记录
func CreateDnsRecord(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	provider := c.Query("provider")
	domainID := c.Query("domain_id")
	subDomain := c.Query("sub_domain")
	recordType := c.Query("record_type")
	value := c.Query("value")
	recordLine := c.DefaultQuery("record_line", "默认")
	ttlStr := c.DefaultQuery("ttl", "600")
	ttl, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil || ttl <= 0 {
		ttl = 600 // 默认TTL
	}
	// 确保TTL在合理范围内 (1秒到2147483647秒)
	if ttl < 1 {
		ttl = 1
	} else if ttl > 2147483647 {
		ttl = 2147483647
	}

	if domainID == "" || subDomain == "" || recordType == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 使用通用创建记录方法
	record, err := dnsService.CreateRecord(domainID, subDomain, recordType, value, recordLine, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "DNS记录创建成功",
		"data": record,
	})
}

// 更新DNS记录
func UpdateDnsRecord(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	provider := c.Query("provider")
	recordID := c.Param("id")
	domainID := c.Query("domain_id")
	subDomain := c.Query("sub_domain")
	recordType := c.Query("record_type")
	value := c.Query("value")
	recordLine := c.DefaultQuery("record_line", "默认")
	ttlStr := c.DefaultQuery("ttl", "600")
	ttl, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil || ttl <= 0 {
		ttl = 600 // 默认TTL
	}
	// 确保TTL在合理范围内 (1秒到2147483647秒)
	if ttl < 1 {
		ttl = 1
	} else if ttl > 2147483647 {
		ttl = 2147483647
	}

	if recordID == "" || domainID == "" || subDomain == "" || recordType == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 使用通用更新记录方法
	record, err := dnsService.UpdateRecord(recordID, domainID, subDomain, recordType, value, recordLine, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "DNS记录更新成功",
		"data": record,
	})
}

// 删除DNS记录
func DeleteDnsRecord(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	provider := c.Query("provider")
	recordID := c.Param("id")
	domainID := c.Query("domain_id")

	if recordID == "" || domainID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 使用通用删除记录方法
	err := dnsService.DeleteRecord(recordID, domainID, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "DNS记录删除成功",
		"data": make(map[string]interface{}),
	})
}

// 设置DNS记录状态
func SetDnsRecordStatus(c *gin.Context) {
	// 检查配置中的Token是否设置
	if !isAnyDnsProviderConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token、阿里云AccessKey或火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	provider := c.Query("provider")
	recordID := c.Param("id")
	domainID := c.Query("domain_id")
	status := c.Query("status")

	if recordID == "" || domainID == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 验证状态参数
	if status != "enable" && status != "disable" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "状态参数必须是enable或disable",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 使用通用设置记录状态方法
	err := dnsService.SetRecordStatus(recordID, domainID, status, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "DNS记录状态设置成功",
		"data": make(map[string]interface{}),
	})
}

// DomainVerificationRequest 域名验证请求
type DomainVerificationRequest struct {
	Domain string `json:"domain" binding:"required"`
	Type   string `json:"type" binding:"required"` // "dns", "file", "cname"
}

// VerifyDomainChallenge 验证域名挑战
type VerifyDomainChallenge struct {
	Domain    string `json:"domain" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Challenge string `json:"challenge" binding:"required"`
}

// GenerateDomainVerificationChallenge 生成域名验证挑战
func GenerateDomainVerificationChallenge(c *gin.Context) {
	var json DomainVerificationRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	// 创建DNS管理器
	dnsManager := dns.NewDnsManager(
		setting.DnsPodToken,
		setting.AliyunAccessKeyId,
		setting.AliyunAccessKeySecret,
		setting.AliyunRegionId,
		setting.VolcengineAccessKeyId,
		setting.VolcengineAccessKeySecret,
		setting.VolcengineRegionId,
	)

	// 创建域名验证服务
	verificationService := dns.NewDomainVerificationService(dnsManager)

	// 生成验证挑战
	response, err := verificationService.GenerateChallenge(json.Domain, json.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "生成验证挑战成功",
		"data": response,
	})
}

// VerifyDomain 通过挑战验证域名
func VerifyDomain(c *gin.Context) {
	var json VerifyDomainChallenge
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	// 创建DNS管理器
	dnsManager := dns.NewDnsManager(
		setting.DnsPodToken,
		setting.AliyunAccessKeyId,
		setting.AliyunAccessKeySecret,
		setting.AliyunRegionId,
		setting.VolcengineAccessKeyId,
		setting.VolcengineAccessKeySecret,
		setting.VolcengineRegionId,
	)

	// 创建域名验证服务
	verificationService := dns.NewDomainVerificationService(dnsManager)

	// 执行域名验证
	response, err := verificationService.VerifyDomain(json.Domain, json.Type, json.Challenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  response.Message,
		"data": response,
	})
}

// isAnyDnsProviderConfigured 检查是否有任何DNS提供商已配置
func isAnyDnsProviderConfigured() bool {
	return setting.DnsPodToken != "" || setting.AliyunAccessKeyId != "" || setting.VolcengineAccessKeyId != ""
}
