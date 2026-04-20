package system

import (
	"errors"
	"strconv"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

func (casbinService *CasbinService) UpdateCasbin(adminAuthorityID, AuthorityID uint, casbinInfos []request.CasbinInfo) error {
	err := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, AuthorityID)
	if err != nil {
		return err
	}

	if global.GVA_CONFIG.System.UseStrictAuth {
		apis, e := ApiServiceApp.GetAllApis(adminAuthorityID)
		if e != nil {
			return e
		}

		for i := range casbinInfos {
			hasAPI := false
			for j := range apis {
				if apis[j].Path == casbinInfos[i].Path && apis[j].Method == casbinInfos[i].Method {
					hasAPI = true
					break
				}
			}
			if !hasAPI {
				return errors.New("存在 api 不在权限列表中")
			}
		}
	}

	authorityID := strconv.Itoa(int(AuthorityID))
	casbinService.ClearCasbin(0, authorityID)

	rules := [][]string{}
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityID + v.Path + v.Method
		if _, ok := deduplicateMap[key]; ok {
			continue
		}
		deduplicateMap[key] = true
		rules = append(rules, []string{authorityID, v.Path, v.Method})
	}

	if len(rules) == 0 {
		return nil
	}

	e := utils.GetCasbin()
	if e == nil {
		return errors.New("权限系统未就绪，请稍后重试")
	}

	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同 api，添加失败，请联系管理员")
	}
	return nil
}

func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GVA_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	if err != nil {
		return err
	}

	e := utils.GetCasbin()
	if e == nil {
		return errors.New("权限系统未就绪，请稍后重试")
	}
	return e.LoadPolicy()
}

func (casbinService *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []request.CasbinInfo) {
	e := utils.GetCasbin()
	if e == nil {
		return pathMaps
	}

	authorityID := strconv.Itoa(int(AuthorityID))
	list, _ := e.GetFilteredPolicy(0, authorityID)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := utils.GetCasbin()
	if e == nil {
		return false
	}
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

func (casbinService *CasbinService) RemoveFilteredPolicy(db *gorm.DB, authorityID string) error {
	return db.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityID).Error
}

func (casbinService *CasbinService) SyncPolicy(db *gorm.DB, authorityID string, rules [][]string) error {
	err := casbinService.RemoveFilteredPolicy(db, authorityID)
	if err != nil {
		return err
	}
	return casbinService.AddPolicies(db, rules)
}

func (casbinService *CasbinService) AddPolicies(db *gorm.DB, rules [][]string) error {
	var casbinRules []gormadapter.CasbinRule
	for i := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rules[i][0],
			V1:    rules[i][1],
			V2:    rules[i][2],
		})
	}
	return db.Create(&casbinRules).Error
}

func (casbinService *CasbinService) FreshCasbin() error {
	e := utils.GetCasbin()
	if e == nil {
		return errors.New("权限系统未就绪，请稍后重试")
	}
	return e.LoadPolicy()
}
