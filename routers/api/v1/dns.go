package v1

import (
	"net/http"
	"strconv"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
)

// 获取域名列表
func GetDomains(c *gin.Context) {
	// 检查配置中的Token是否设置
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()

	// 检查使用哪个提供商
	provider := c.Query("provider")
	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 阿里云域名列表需要特殊处理
		pageNumber := 1
		pageSize := 20
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				pageNumber = p
			}
		}
		if sizeStr := c.Query("size"); sizeStr != "" {
			if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
				pageSize = s
			}
		}

		// 使用阿里云DNS获取域名列表
		_, err := dnsService.Manager.GetAliyunDomainList(pageNumber, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}

		// 暂时返回空列表，因为阿里云域名结构与DNSPod不同
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  "success",
			"data": make([]interface{}, 0),
		})
		return
	}

	// 使用DNSPod
	domains, err := dnsService.GetDomainList()
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
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
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

	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 使用阿里云DNS
		records, err := dnsService.GetAliyunRecordList(domain, subDomain)
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
		return
	}

	// 使用DNSPod
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
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
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
	if err != nil {
		ttl = 600 // 默认TTL
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

	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 使用阿里云DNS
		record, err := dnsService.CreateAliyunRecord(domainID, subDomain, recordType, value, ttl)
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
			"msg":  "阿里云DNS记录创建成功",
			"data": record,
		})
		return
	}

	// 使用DNSPod
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
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
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
	if err != nil {
		ttl = 600 // 默认TTL
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

	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 使用阿里云DNS
		record, err := dnsService.UpdateAliyunRecord(recordID, subDomain, recordType, value, ttl)
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
			"msg":  "阿里云DNS记录更新成功",
			"data": record,
		})
		return
	}

	// 使用DNSPod
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
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
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

	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 使用阿里云DNS
		err := dnsService.DeleteAliyunRecord(recordID)
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
			"msg":  "阿里云DNS记录删除成功",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 使用DNSPod
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
	if setting.DnsPodToken == "" && setting.AliyunAccessKeyId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token或阿里云AccessKey未配置",
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

	if provider == "aliyun" && dnsService.Manager.UseAliyunDns() {
		// 使用阿里云DNS，需要转换状态
		aliyunStatus := "ENABLE"
		if status == "disable" {
			aliyunStatus = "DISABLE"
		}

		err := dnsService.SetAliyunRecordStatus(recordID, aliyunStatus)
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
			"msg":  "阿里云DNS记录状态设置成功",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 使用DNSPod
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
