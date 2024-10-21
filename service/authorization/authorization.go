package authorization

import "github.com/casbin/casbin/v2"

type Authorization interface {
	// Returns true if the user has the permission
	IsAuthorized(username string, domain string, resource string, action string) (bool, error)
}

var _ Authorization = &CasbinAuthorizationService{}

type CasbinAuthorizationService struct {
	enforcer *casbin.Enforcer
}

func NewCasbinAuthorizationService(casbinEnforcer *casbin.Enforcer) *CasbinAuthorizationService {
	return &CasbinAuthorizationService{enforcer: casbinEnforcer}
}

func (a *CasbinAuthorizationService) IsAuthorized(username string, domain string, resource string, action string) (bool, error) {
	return a.enforcer.Enforce(username, domain, resource, action)
}
