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
	"go.mongodb.org/mongo-driver/bson"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := bson.M{}
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

	// 获取分页参数
	pageNum := util.GetPage(c)

	tags, err := models.GetTags(pageNum, setting.PageSize, maps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "获取标签列表失败: " + err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	total, err := models.GetTagTotal(maps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "获取标签总数失败: " + err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	data["lists"] = tags
	data["total"] = total
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
		exists, err := models.ExistTagByName(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  "检查标签是否存在失败: " + err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}
		if exists {
			code = e.ERROR_EXIST_TAG
		} else {
			err := models.AddTag(name, state, createdBy)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.ERROR,
					"msg":  "添加标签失败: " + err.Error(),
					"data": make(map[string]interface{}),
				})
				return
			}
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
	id := c.Param("id")
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	modifiedBy := c.Query("modified_by")
	vaild := validation.Validation{}
	vaild.Required(id, "id").Message("ID不能为空")
	vaild.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	vaild.MaxSize(modifiedBy, 100, "modified_by").Message("修改人不能超过100字符")
	vaild.MaxSize(name, 100, "name").Message("标签名最长为100字符")
	vaild.Range(state, 0, 1, "state").Message("状态只能为0或1")
	code := e.INVALID_PARAMS
	if !vaild.HasErrors() {
		exists, err := models.ExistTagByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  "检查标签是否存在失败: " + err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}
		if exists {
			data := bson.M{}
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state >= 0 {
				data["state"] = state
			}

			err := models.EditTag(id, data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.ERROR,
					"msg":  "更新标签失败: " + err.Error(),
					"data": make(map[string]interface{}),
				})
				return
			}
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
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
	id := c.Param("id")
	vaild := validation.Validation{}
	vaild.Required(id, "id").Message("ID不能为空")
	code := e.INVALID_PARAMS
	if !vaild.HasErrors() {
		exists, err := models.ExistTagByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  "检查标签是否存在失败: " + err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}
		if exists {

			err := models.DeleteTag(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.ERROR,
					"msg":  "删除标签失败: " + err.Error(),
					"data": make(map[string]interface{}),
				})
				return
			}
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": make(map[string]string),
		"msg":  e.GetMsg(code),
	})
}
