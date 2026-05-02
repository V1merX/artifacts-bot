package di

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type authHolder struct {
	mu          sync.RWMutex
	basicUser   string
	basicPass   string
	bearerToken string
}

func newAuthHolder(user, pass string) *authHolder {
	return &authHolder{basicUser: user, basicPass: pass}
}

func (a *authHolder) SetBearerToken(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.bearerToken = token
}

func (a *authHolder) Editor(_ context.Context, req *http.Request) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.bearerToken))
	} else {
		req.SetBasicAuth(a.basicUser, a.basicPass)
	}
	return nil
}
