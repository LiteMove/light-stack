package permission

import (
	"sync"
)

type PermissionCache struct {
	sync.RWMutex
	userPermissions map[uint64]map[string]bool
	userRoles       map[uint64]map[string]bool
}

var Cache = &PermissionCache{
	userPermissions: make(map[uint64]map[string]bool),
	userRoles:       make(map[uint64]map[string]bool),
}

func (p *PermissionCache) LoadUserPermissions(userID uint64, permissions []string) {
	p.Lock()
	defer p.Unlock()

	userPerms := make(map[string]bool, len(permissions))
	for _, perm := range permissions {
		userPerms[perm] = true
	}
	p.userPermissions[userID] = userPerms
}

func (p *PermissionCache) LoadUserRoles(userID uint64, roles []string) {
	p.Lock()
	defer p.Unlock()

	userRoles := make(map[string]bool, len(roles))
	for _, role := range roles {
		userRoles[role] = true
	}
	p.userRoles[userID] = userRoles
}

func (p *PermissionCache) GetUserPermissions(userID uint64) (map[string]bool, bool) {
	p.RLock()
	defer p.RUnlock()

	perms, exists := p.userPermissions[userID]
	return perms, exists
}

func (p *PermissionCache) GetUserRoles(userID uint64) (map[string]bool, bool) {
	p.RLock()
	defer p.RUnlock()

	roles, exists := p.userRoles[userID]
	return roles, exists
}

func (p *PermissionCache) ClearUserPermissions(userID uint64) {
	p.Lock()
	defer p.Unlock()

	delete(p.userPermissions, userID)
	delete(p.userRoles, userID)
}

func (p *PermissionCache) HasPermission(userID uint64, code string) bool {
	p.RLock()
	defer p.RUnlock()

	perms, exists := p.userPermissions[userID]
	if !exists {
		return false
	}

	return perms[code]
}

func (p *PermissionCache) HasAnyPermission(userID uint64, codes ...string) bool {
	p.RLock()
	defer p.RUnlock()

	perms, exists := p.userPermissions[userID]
	if !exists {
		return false
	}

	for _, code := range codes {
		if perms[code] {
			return true
		}
	}

	return false
}

func (p *PermissionCache) HasRole(userID uint64, roleCode string) bool {
	p.RLock()
	defer p.RUnlock()

	roles, exists := p.userRoles[userID]
	if !exists {
		return false
	}

	return roles[roleCode]
}

func (p *PermissionCache) HasAnyRole(userID uint64, roleCodes ...string) bool {
	p.RLock()
	defer p.RUnlock()

	roles, exists := p.userRoles[userID]
	if !exists {
		return false
	}

	for _, code := range roleCodes {
		if roles[code] {
			return true
		}
	}

	return false
}
