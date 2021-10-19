package application

import (
	"errors"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/commons"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/domain"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/ports"
	"time"
)

type AuthManager struct {
	userManager  *UserManager
	tokenHandler ports.TokenHandler
}

func NewAuthManager(userManager  *UserManager, 	tokenHandler ports.TokenHandler) *AuthManager {
	return &AuthManager{userManager : userManager, tokenHandler: tokenHandler}
}

// Login
func (p *AuthManager) Login(userName, password string) (tokenString string, err error) {
	account, err := p.userManager.FindAccountByName(userName)

	if err != nil {
		return
	}

	loginSuccess, err := account.CheckPassword(password)
	if err != nil {
		return
	}

	if !loginSuccess {
		err = &domain.UnauthorizedError{}
		return
	}

	token, err := p.CreateToken(account.Id, account.Name, account.TenantId)

	if err != nil {
		return
	} else {
		tokenString = token.TokenString
	}

	return
}

// Logout
func (p *AuthManager) Logout(tokenString string) (string, error) {
	token, err := p.tokenHandler.GetToken(tokenString)

	if err != nil {
		return "", &domain.UnauthorizedError{}
	} else if !token.IsValid() {
		return "", nil
	} else {
		accountName := token.AccountName
		token.Destroy()

		err := p.tokenHandler.Modify(&token)
		if err != nil {
			return "", err
		}

		return accountName, nil
	}
}

var SkipAuth = true

// Accessible
func (p *AuthManager) Accessible(pathType string, path string, tokenString string) (tenantId string, accountId, accountName string, err error) {
	if path == "" {
		err = errors.New("path cannot be blank")
		return
	}

	token, err := p.tokenHandler.GetToken(tokenString)

	if err != nil {
		return
	}

	accountId = token.AccountId
	accountName = token.AccountName
	tenantId = token.TenantId

	if SkipAuth {
		// todo checkAuth switch
		return
	}

	// 校验token有效
	if !token.IsValid() {
		err = &domain.UnauthorizedError{}
		return
	}

	// 根据token查用户
	account, err := p.userManager.findAccountAggregation(accountName)
	if err != nil {
		return
	}

	// 查权限
	permission, err := p.userManager.findPermissionAggregationByCode(tenantId, path)
	if err != nil {
		return
	}

	ok, err := p.checkAuth(account, permission)

	if err != nil {
		return
	}

	if !ok {
		err = &domain.ForbiddenError{}
	}

	return
}

// checkAuth
func (p *AuthManager) checkAuth(account *domain.AccountAggregation, permission *domain.PermissionAggregation) (bool, error) {

	accountRoles := account.Roles

	if accountRoles == nil || len(accountRoles) == 0 {
		return false, nil
	}

	accountRoleMap := make(map[string]bool)

	for _, r := range accountRoles {
		accountRoleMap[r.Id] = true
	}

	allowedRoles := permission.Roles

	if allowedRoles == nil || len(allowedRoles) == 0 {
		return false, nil
	}

	for _, r := range allowedRoles {
		if _, exist := accountRoleMap[r.Id]; exist {
			return true, nil
		}
	}

	return false, nil
}

func (p *AuthManager) CreateToken(accountId string, accountName string, tenantId string) (domain.TiEMToken, error) {
	token := domain.TiEMToken{
		AccountName: accountName,
		AccountId: accountId,
		TenantId: tenantId,
		ExpirationTime: time.Now().Add(commons.DefaultTokenValidPeriod),
	}

	tokenString, err := p.tokenHandler.Provide(&token)
	token.TokenString = tokenString
	return token, err
}

func (p *AuthManager) SetTokenHandler(tokenHandler ports.TokenHandler) {
	p.tokenHandler = tokenHandler
}

func (p *AuthManager) SetUserManager(userManager *UserManager) {
	p.userManager = userManager
}