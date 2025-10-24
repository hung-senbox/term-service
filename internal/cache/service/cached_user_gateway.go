package service

import (
	"context"
	"term-service/internal/cache"
	"term-service/internal/gateway"
	gw_response "term-service/internal/gateway/dto/response"
	"time"
)

type CachedUserGateway struct {
	inner gateway.UserGateway
	cache cache.Cache
	ttl   time.Duration
}

func NewCachedUserGateway(inner gateway.UserGateway, cache cache.Cache, ttlSeconds int) gateway.UserGateway {
	return &CachedUserGateway{
		inner: inner,
		cache: cache,
		ttl:   time.Duration(ttlSeconds) * time.Second,
	}
}

// ==============================
// === GetStudentInfo ===
// ==============================
func (g *CachedUserGateway) GetStudentInfo(ctx context.Context, studentID string) (*gw_response.StudentResponse, error) {
	cacheKey := cache.StudentCacheKey(studentID)

	var cached gw_response.StudentResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	student, err := g.inner.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, student, int(g.ttl.Seconds()))
	return student, nil
}

func (g *CachedUserGateway) GetCurrentUser(ctx context.Context) (*gw_response.CurrentUser, error) {
	return g.inner.GetCurrentUser(ctx)
}
