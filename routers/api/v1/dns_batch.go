package v1

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
)

// 批量创建DNS记录
func BatchCreateDnsRecords(c *gin.Context) {
	provider := c.Query("provider")

	// 检查配置中的Token是否设置
	if provider == "dns_pod" && setting.DnsPodToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "aliyun" && (setting.AliyunAccessKeyId == "" || setting.AliyunAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "阿里云AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "volcengine" && (setting.VolcengineAccessKeyId == "" || setting.VolcengineAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 从请求体获取批量数据
	var records []struct {
		DomainID string `json:"domain_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Value    string `json:"value"`
		Line     string `json:"line"`
		TTL      int64  `json:"ttl"`
		Remark   string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "请求数据格式错误",
			"data": make(map[string]interface{}),
		})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "至少需要提供一个DNS记录",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()
	var results []map[string]interface{}

	for _, record := range records {
		if record.Name == "" || record.Type == "" || record.Value == "" {
			results = append(results, map[string]interface{}{
				"success": false,
				"error":   "参数不完整",
				"name":    record.Name,
			})
			continue
		}

		// 使用策略模式处理不同类型提供商
		ttl := record.TTL
		if ttl == 0 {
			ttl = 600 // 默认TTL
		}

		var result map[string]interface{}
		line := record.Line
		if line == "" {
			line = "默认"
		}

		recordResult, err := dnsService.CreateRecord(record.DomainID, record.Name, record.Type, record.Value, line, provider)
		if err != nil {
			result = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"name":    record.Name,
			}
		} else {
			result = map[string]interface{}{
				"success": true,
				"data":    recordResult,
				"name":    record.Name,
			}
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "批量创建完成",
		"data": map[string]interface{}{
			"results": results,
			"total":   len(results),
			"success": countSuccess(results),
		},
	})
}

// 批量更新DNS记录
func BatchUpdateDnsRecords(c *gin.Context) {
	provider := c.Query("provider")

	// 检查配置中的Token是否设置
	if provider == "dns_pod" && setting.DnsPodToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "aliyun" && (setting.AliyunAccessKeyId == "" || setting.AliyunAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "阿里云AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "volcengine" && (setting.VolcengineAccessKeyId == "" || setting.VolcengineAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 从请求体获取批量数据
	var updates []struct {
		ID       string `json:"id"`
		DomainID string `json:"domain_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Value    string `json:"value"`
		Line     string `json:"line"`
		TTL      int64  `json:"ttl"`
		Remark   string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "请求数据格式错误",
			"data": make(map[string]interface{}),
		})
		return
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "至少需要提供一个DNS记录",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()
	var results []map[string]interface{}

	for _, update := range updates {
		if update.ID == "" || update.Name == "" || update.Type == "" || update.Value == "" {
			results = append(results, map[string]interface{}{
				"success": false,
				"error":   "参数不完整",
				"id":      update.ID,
			})
			continue
		}

		// 使用策略模式处理不同类型提供商
		ttl := update.TTL
		if ttl == 0 {
			ttl = 600 // 默认TTL
		}

		var result map[string]interface{}
		line := update.Line
		if line == "" {
			line = "默认"
		}

		recordResult, err := dnsService.UpdateRecord(update.ID, update.DomainID, update.Name, update.Type, update.Value, line, provider)
		if err != nil {
			result = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"id":      update.ID,
			}
		} else {
			result = map[string]interface{}{
				"success": true,
				"data":    recordResult,
				"id":      update.ID,
			}
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "批量更新完成",
		"data": map[string]interface{}{
			"results": results,
			"total":   len(results),
			"success": countSuccess(results),
		},
	})
}

// 批量删除DNS记录
func BatchDeleteDnsRecords(c *gin.Context) {
	provider := c.Query("provider")

	// 检查配置中的Token是否设置
	if provider == "dns_pod" && setting.DnsPodToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "aliyun" && (setting.AliyunAccessKeyId == "" || setting.AliyunAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "阿里云AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "volcengine" && (setting.VolcengineAccessKeyId == "" || setting.VolcengineAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 从请求体获取批量数据
	var deletes []struct {
		ID       string `json:"id"`
		DomainID string `json:"domain_id"`
	}

	if err := c.ShouldBindJSON(&deletes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "请求数据格式错误",
			"data": make(map[string]interface{}),
		})
		return
	}

	if len(deletes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "至少需要提供一个DNS记录ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()
	var results []map[string]interface{}

	for _, delete := range deletes {
		if delete.ID == "" {
			results = append(results, map[string]interface{}{
				"success": false,
				"error":   "记录ID不能为空",
				"id":      delete.ID,
			})
			continue
		}

		// 使用策略模式处理不同类型提供商
		err := dnsService.DeleteRecord(delete.ID, delete.DomainID, provider)
		var result map[string]interface{}
		if err != nil {
			result = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"id":      delete.ID,
			}
		} else {
			result = map[string]interface{}{
				"success": true,
				"id":      delete.ID,
			}
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "批量删除完成",
		"data": map[string]interface{}{
			"results": results,
			"total":   len(results),
			"success": countSuccess(results),
		},
	})
}

// 批量更新DNS记录状态
func BatchUpdateDnsRecordStatus(c *gin.Context) {
	provider := c.Query("provider")
	status := c.Query("status")

	if status != "enable" && status != "disable" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "状态参数必须是enable或disable",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 检查配置中的Token是否设置
	if provider == "dns_pod" && setting.DnsPodToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "DNSPod Token未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "aliyun" && (setting.AliyunAccessKeyId == "" || setting.AliyunAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "阿里云AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	if provider == "volcengine" && (setting.VolcengineAccessKeyId == "" || setting.VolcengineAccessKeySecret == "") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "火山引擎AccessKey未配置",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 从请求体获取批量数据
	var statusUpdates []struct {
		ID       string `json:"id"`
		DomainID string `json:"domain_id"`
	}

	if err := c.ShouldBindJSON(&statusUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "请求数据格式错误",
			"data": make(map[string]interface{}),
		})
		return
	}

	if len(statusUpdates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "至少需要提供一个DNS记录ID",
			"data": make(map[string]interface{}),
		})
		return
	}

	dnsService := models.NewDnsService()
	var results []map[string]interface{}

	for _, update := range statusUpdates {
		if update.ID == "" {
			results = append(results, map[string]interface{}{
				"success": false,
				"error":   "记录ID不能为空",
				"id":      update.ID,
			})
			continue
		}

		// 使用策略模式处理不同类型提供商
		err := dnsService.SetRecordStatus(update.ID, update.DomainID, status, provider)
		var result map[string]interface{}
		if err != nil {
			result = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"id":      update.ID,
			}
		} else {
			result = map[string]interface{}{
				"success": true,
				"id":      update.ID,
			}
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  "批量更新状态完成",
		"data": map[string]interface{}{
			"results": results,
			"total":   len(results),
			"success": countSuccess(results),
		},
	})
}

// 辅助函数：计算成功数量
func countSuccess(results []map[string]interface{}) int {
	count := 0
	for _, result := range results {
		if success, ok := result["success"].(bool); ok && success {
			count++
		}
	}
	return count
}
