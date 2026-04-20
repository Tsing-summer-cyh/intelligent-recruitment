package utils

import (
	"errors"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	casbinDB             *gorm.DB
	casbinMu             sync.Mutex
)

// GetCasbin 获取 casbin 实例。数据库未就绪或初始化失败时返回 nil，
// 这样在初始化完成后可以重试构建，而不会被 sync.Once 锁死。
func GetCasbin() *casbin.SyncedCachedEnforcer {
	db := global.GVA_DB
	if db == nil {
		zap.L().Warn("casbin enforcer skipped because database is not ready")
		return nil
	}

	casbinMu.Lock()
	defer casbinMu.Unlock()

	if syncedCachedEnforcer != nil && casbinDB == db {
		return syncedCachedEnforcer
	}

	enforcer, err := newCasbinEnforcer(db)
	if err != nil {
		zap.L().Error("casbin enforcer initialization failed", zap.Error(err))
		return nil
	}

	syncedCachedEnforcer = enforcer
	casbinDB = db
	return syncedCachedEnforcer
}

func newCasbinEnforcer(db *gorm.DB) (*casbin.SyncedCachedEnforcer, error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	text := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
	`
	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	enforcer.SetExpireTime(60 * 60)
	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}
