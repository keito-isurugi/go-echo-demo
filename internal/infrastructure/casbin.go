package infrastructure

import (
	"log"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/usecase"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewCasbinEnforcer Casbinエンフォーサーを作成
func NewCasbinEnforcer() (*casbin.Enforcer, error) {
	// モデル設定を読み込み
	m, err := model.NewModelFromFile("config/rbac_model.conf")
	if err != nil {
		log.Printf("Warning: モデル設定ファイルの読み込みに失敗しました: %v", err)
		// デフォルトのモデルを使用
		m, _ = model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	}

	// ファイルアダプターを作成
	adapter := fileadapter.NewAdapter("config/rbac_policy.csv")

	// エンフォーサーを作成
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 自動保存を有効化
	enforcer.EnableAutoSave(true)

	return enforcer, nil
}

// NewCasbinRBACRepository Casbin RBACリポジトリを作成
func NewCasbinRBACRepository() (domain.CasbinRBACRepository, error) {
	enforcer, err := NewCasbinEnforcer()
	if err != nil {
		return nil, err
	}

	return domain.NewCasbinEnforcer(enforcer), nil
}

// NewCasbinRBACUsecase Casbin RBACユースケースを作成
func NewCasbinRBACUsecase(casbinRepo domain.CasbinRBACRepository, rbacRepo domain.RBACRepository) domain.CasbinRBACUsecase {
	return usecase.NewCasbinRBACUsecase(casbinRepo, rbacRepo)
}
