package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"xorm.io/xorm"
)

const RBAC_MODEL_CONF = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

const (
	NORMAL_USER = "normal_user"
)

func NewEnforcer(engine *xorm.Engine) (*casbin.Enforcer, error) {
	adapter, err := xormadapter.NewAdapterByEngineWithTableName(engine, "sys_casbin", "")

	cabinModel := model.NewModel()
	err = cabinModel.LoadModelFromText(RBAC_MODEL_CONF)
	if err != nil {
		return nil, err
	}
	enf, err := casbin.NewEnforcer(cabinModel, adapter)
	if err != nil {
		return nil, err
	}
	err = enf.LoadPolicy()
	if err != nil {
		return nil, err
	}

	// init casbin policy
	_, err = enf.AddPolicy(NORMAL_USER, "/api/v1/api/*", "*")
	if err != nil {
		return nil, err
	}

	_, err = enf.AddPolicy(NORMAL_USER, "/api/v1/disease/*", "*")
	if err != nil {
		return nil, err
	}

	return enf, nil
}
