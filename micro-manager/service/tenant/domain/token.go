package domain

import (
	"github.com/pingcap/ticp/micro-manager/service/tenant/commons"
	"github.com/pingcap/ticp/micro-manager/service/tenant/port"
	"time"
)

type TiCPToken struct {
	TokenString 	string
	AccountName		string
	TenantId		uint
	ExpirationTime  time.Time
}

func (token *TiCPToken) destroy() error {
	token.ExpirationTime = time.Now()
	return port.TokenMNG.Modify(token)
}

func (token *TiCPToken) renew() error {
	token.ExpirationTime = time.Now().Add(commons.DefaultTokenValidPeriod)
	return port.TokenMNG.Modify(token)
}

func (token *TiCPToken) isValid() bool {
	now := time.Now()

	return now.Before(token.ExpirationTime)
}

func createToken(accountName string, tenantId uint) (TiCPToken, error) {
	token := TiCPToken{
		AccountName: accountName,
		TenantId: tenantId,
		ExpirationTime: time.Now().Add(commons.DefaultTokenValidPeriod),
	}

	tokenString, err := port.TokenMNG.Provide(&token)
	token.TokenString = tokenString
	return token, err
}