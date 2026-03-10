package controller

import (
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/gin-gonic/gin"
)

// filterGroupsByUsername 根据用户名前缀过滤分组
// 只保留以用户名为前缀的分组和以 "default" 为前缀的分组
// root 用户（role >= RoleRootUser）不进行过滤，返回 true 表示通过
func filterGroupsByUsername(groupName string, username string, role int) bool {
	// root 用户不过滤
	if role >= common.RoleRootUser {
		return true
	}
	// 以用户名为前缀的分组
	if username != "" && strings.HasPrefix(groupName, username) {
		return true
	}
	// 以 "default" 为前缀的分组
	if strings.HasPrefix(groupName, "default") {
		return true
	}
	return false
}

func GetGroups(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetInt("role")

	groupNames := make([]string, 0)
	for groupName := range ratio_setting.GetGroupRatioCopy() {
		if filterGroupsByUsername(groupName, username, role) {
			groupNames = append(groupNames, groupName)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groupNames,
	})
}

func GetUserGroups(c *gin.Context) {
	usableGroups := make(map[string]map[string]interface{})
	userGroup := ""
	userId := c.GetInt("id")
	userGroup, _ = model.GetUserGroup(userId, false)
	userUsableGroups := service.GetUserUsableGroups(userGroup)

	// 获取用户名和角色用于前缀过滤
	username := c.GetString("username")
	role := c.GetInt("role")
	// 如果 context 中没有 username，尝试从用户缓存获取
	if username == "" {
		if user, err := model.GetUserCache(userId); err == nil {
			username = user.Username
		}
	}

	for groupName, groupRatioVal := range ratio_setting.GetGroupRatioCopy() {
		// 按用户名前缀过滤
		if !filterGroupsByUsername(groupName, username, role) {
			continue
		}
		// 如果分组在管理员配置的可用分组列表中，使用其描述
		if desc, ok := userUsableGroups[groupName]; ok {
			usableGroups[groupName] = map[string]interface{}{
				"ratio": service.GetUserGroupRatio(userGroup, groupName),
				"desc":  desc,
			}
		} else if username != "" && strings.HasPrefix(groupName, username) {
			// 用户名前缀的分组即使不在管理员可用分组列表中，也显示出来
			usableGroups[groupName] = map[string]interface{}{
				"ratio": groupRatioVal,
				"desc":  groupName,
			}
		}
	}
	if _, ok := userUsableGroups["auto"]; ok {
		// auto 分组作为特殊分组，不受前缀过滤影响
		usableGroups["auto"] = map[string]interface{}{
			"ratio": "自动",
			"desc":  setting.GetUsableGroupDescription("auto"),
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    usableGroups,
	})
}
