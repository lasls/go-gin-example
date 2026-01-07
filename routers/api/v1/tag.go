package v1

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}
	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  e.GetMsg(code),
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	vaild := validation.Validation{}
	vaild.Required(name, "name").Message("标签名不能为空")
	vaild.MaxSize(name, 100, "name").Message("标签名最长为100字符")
	vaild.Required(createdBy, "created_by").Message("创建人不能为空")
	vaild.MaxSize(createdBy, 100, "created_by").Message("创建人不能超过100字符")
	vaild.Range(state, 0, 1, "state").Message("状态只能为0或1")
	code := e.INVALID_PARAMS
	if !vaild.HasErrors() {
		if models.ExistTagByName(name) {
			code = e.ERROR_EXIST_TAG
		} else {
			models.AddTag(name, state, createdBy)
			code = e.SUCCESS
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": make(map[string]string),
		"msg":  e.GetMsg(code),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	modifiedBy := c.Query("modified_by")
	vaild := validation.Validation{}
	vaild.Min(id, 1, "id").Message("ID必须大于0")
	vaild.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	vaild.MaxSize(modifiedBy, 100, "modified_by").Message("修改人不能超过100字符")
	vaild.MaxSize(name, 100, "name").Message("标签名最长为100字符")
	vaild.Range(state, 0, 1, "state").Message("状态只能为0或1")
	code := e.INVALID_PARAMS
	if !vaild.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			data["name"] = name
			data["state"] = state
			models.EditTag(id, data)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": make(map[string]string),
		"msg":  e.GetMsg(code),
	})
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
}
