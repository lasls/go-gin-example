package v1

import (
	"net/http"
	"strconv"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
)

// 获取域名列表
func GetDnsDomains(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}

	page := util.GetPage(c)
	pageSize := setting.PageSize // 使用配置中的页面大小

	domains, err := models.GetDnsDomainList(page, pageSize, maps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	total, err := models.GetDnsDomainTotal(maps)
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
		"data": map[string]interface{}{
			"lists": domains,
			"total": total,
		},
	})
}

// 添加域名
func AddDnsDomain(c *gin.Context) {
	name := c.Query("name")
	provider := c.Query("provider")
	domainID := c.Query("domain_id")
	status := c.DefaultQuery("status", "active")
	grade := c.Query("grade")
	owner := c.Query("owner")
	remark := c.Query("remark")

	if name == "" || provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	domain := &models.DnsDomain{
		Name:     name,
		Provider: provider,
		DomainID: domainID,
		Status:   status,
		Grade:    grade,
		Owner:    owner,
		Remark:   remark,
	}

	err := models.AddDnsDomain(domain)
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
		"msg":  "域名添加成功",
		"data": domain,
	})
}

// 更新域名
func UpdateDnsDomain(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "无效的域名ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.ExistDnsDomainByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  "域名不存在",
			"data": make(map[string]interface{}),
		})
		return
	}

	name := c.Query("name")
	provider := c.Query("provider")
	domainID := c.Query("domain_id")
	status := c.Query("status")
	grade := c.Query("grade")
	owner := c.Query("owner")
	remark := c.Query("remark")

	updateData := make(map[string]interface{})
	if name != "" {
		updateData["name"] = name
	}
	if provider != "" {
		updateData["provider"] = provider
	}
	if domainID != "" {
		updateData["domain_id"] = domainID
	}
	if status != "" {
		updateData["status"] = status
	}
	if grade != "" {
		updateData["grade"] = grade
	}
	if owner != "" {
		updateData["owner"] = owner
	}
	if remark != "" {
		updateData["remark"] = remark
	}

	err = models.UpdateDnsDomain(id, updateData)
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
		"msg":  "域名更新成功",
		"data": make(map[string]interface{}),
	})
}

// 删除域名
func DeleteDnsDomain(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "无效的域名ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.ExistDnsDomainByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  "域名不存在",
			"data": make(map[string]interface{}),
		})
		return
	}

	err = models.DeleteDnsDomain(id)
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
		"msg":  "域名删除成功",
		"data": make(map[string]interface{}),
	})
}

// 获取DNS解析记录列表（数据库）
func GetDnsRecordsDb(c *gin.Context) {
	domainIDStr := c.Query("domain_id")
	maps := make(map[string]interface{})

	if domainIDStr != "" {
		domainID, err := strconv.Atoi(domainIDStr)
		if err == nil && domainID > 0 {
			maps["domain_id"] = domainID
		}
	}

	name := c.Query("name")
	if name != "" {
		maps["name"] = name
	}

	recordType := c.Query("type")
	if recordType != "" {
		maps["type"] = recordType
	}

	page := util.GetPage(c)
	pageSize := setting.PageSize // 使用配置中的页面大小

	records, err := models.GetDnsRecordList(page, pageSize, maps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	total, err := models.GetDnsRecordTotal(maps)
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
		"data": map[string]interface{}{
			"lists": records,
			"total": total,
		},
	})
}

// 添加DNS解析记录（数据库）
func AddDnsRecordDb(c *gin.Context) {
	domainID, err := strconv.Atoi(c.Query("domain_id"))
	if err != nil || domainID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "无效的域名ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	name := c.Query("name")
	recordType := c.Query("type")
	value := c.Query("value")
	provider := c.Query("provider")

	if name == "" || recordType == "" || value == "" || provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "参数不完整",
			"data": make(map[string]interface{}),
		})
		return
	}

	status := c.DefaultQuery("status", "enable")
	line := c.DefaultQuery("line", "默认")
	ttlStr := c.DefaultQuery("ttl", "600")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 600
	}
	remark := c.Query("remark")
	remoteID := c.Query("remote_id")

	record := &models.DnsRecord{
		DomainID: domainID,
		Name:     name,
		Type:     recordType,
		Value:    value,
		Status:   status,
		Line:     line,
		TTL:      ttl,
		Remark:   remark,
		Provider: provider,
		RemoteID: remoteID,
	}

	err = models.AddDnsRecord(record)
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
		"msg":  "DNS解析记录添加成功",
		"data": record,
	})
}

// 更新DNS解析记录（数据库）
func UpdateDnsRecordDb(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "无效的DNS记录ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.ExistDnsRecordByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  "DNS记录不存在",
			"data": make(map[string]interface{}),
		})
		return
	}

	domainIDStr := c.Query("domain_id")
	name := c.Query("name")
	recordType := c.Query("type")
	value := c.Query("value")
	status := c.Query("status")
	line := c.Query("line")
	ttlStr := c.Query("ttl")
	remark := c.Query("remark")
	provider := c.Query("provider")
	remoteID := c.Query("remote_id")

	updateData := make(map[string]interface{})
	if domainIDStr != "" {
		if domainID, err := strconv.Atoi(domainIDStr); err == nil && domainID > 0 {
			updateData["domain_id"] = domainID
		}
	}
	if name != "" {
		updateData["name"] = name
	}
	if recordType != "" {
		updateData["type"] = recordType
	}
	if value != "" {
		updateData["value"] = value
	}
	if status != "" {
		updateData["status"] = status
	}
	if line != "" {
		updateData["line"] = line
	}
	if ttlStr != "" {
		if ttl, err := strconv.Atoi(ttlStr); err == nil {
			updateData["ttl"] = ttl
		}
	}
	if remark != "" {
		updateData["remark"] = remark
	}
	if provider != "" {
		updateData["provider"] = provider
	}
	if remoteID != "" {
		updateData["remote_id"] = remoteID
	}

	err = models.UpdateDnsRecord(id, updateData)
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
		"data": make(map[string]interface{}),
	})
}

// 删除DNS解析记录（数据库）
func DeleteDnsRecordDb(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "无效的DNS记录ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.ExistDnsRecordByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  "DNS记录不存在",
			"data": make(map[string]interface{}),
		})
		return
	}

	err = models.DeleteDnsRecord(id)
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
