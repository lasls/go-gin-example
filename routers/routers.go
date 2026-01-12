package routers

import (
	v1 "github.com/EDDYCJY/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		// DNSPod API路由
		apiV1.GET("/domains", v1.GetDomains)
		apiV1.GET("/dns/records", v1.GetDnsRecords)
		apiV1.POST("/dns/records", v1.CreateDnsRecord)
		apiV1.PUT("/dns/records/:id", v1.UpdateDnsRecord)
		apiV1.DELETE("/dns/records/:id", v1.DeleteDnsRecord)
		apiV1.PUT("/dns/records/:id/status", v1.SetDnsRecordStatus)

		// 域名验证API路由
		apiV1.POST("/dns/verification/challenge", v1.GenerateDomainVerificationChallenge)
		apiV1.POST("/dns/verification/verify", v1.VerifyDomain)

		// DNS数据库API路由
		apiV1.GET("/dns/domains", v1.GetDnsDomains)
		apiV1.POST("/dns/domains", v1.AddDnsDomain)
		apiV1.PUT("/dns/domains/:id", v1.UpdateDnsDomain)
		apiV1.DELETE("/dns/domains/:id", v1.DeleteDnsDomain)
		apiV1.GET("/dns/records_db", v1.GetDnsRecordsDb)
		apiV1.POST("/dns/records_db", v1.AddDnsRecordDb)
		apiV1.PUT("/dns/records_db/:id", v1.UpdateDnsRecordDb)
		apiV1.DELETE("/dns/records_db/:id", v1.DeleteDnsRecordDb)

		// DNS批量操作API路由
		apiV1.POST("/dns/records/batch", v1.BatchCreateDnsRecords)
		apiV1.PUT("/dns/records/batch", v1.BatchUpdateDnsRecords)
		apiV1.DELETE("/dns/records/batch", v1.BatchDeleteDnsRecords)
		apiV1.PUT("/dns/records/batch/status", v1.BatchUpdateDnsRecordStatus)
	}
	return r
}
