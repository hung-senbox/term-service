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

// ==============================
// === GetTeacherInfo ===
// ==============================
func (g *CachedUserGateway) GetTeacherInfo(ctx context.Context, teacherID string) (*gw_response.TeacherResponse, error) {
	cacheKey := cache.TeacherCacheKey(teacherID)

	var cached gw_response.TeacherResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	teacher, err := g.inner.GetTeacherInfo(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, teacher, int(g.ttl.Seconds()))
	return teacher, nil
}

// ==============================
// === GetTeacherByUserAndOrganization ===
// ==============================
func (g *CachedUserGateway) GetTeacherByUserAndOrganization(ctx context.Context, userID, organizationID string) (*gw_response.TeacherResponse, error) {
	cacheKey := cache.UserCacheKey(userID + ":" + organizationID)

	var cached gw_response.TeacherResponse
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	teacher, err := g.inner.GetTeacherByUserAndOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, teacher, int(g.ttl.Seconds()))
	return teacher, nil
}

// ==============================
// === GetUserByTeacher ===
// ==============================
func (g *CachedUserGateway) GetUserByTeacher(ctx context.Context, teacherID string) (*gw_response.CurrentUser, error) {
	cacheKey := cache.TeacherCacheKey(teacherID)

	var cached gw_response.CurrentUser
	if err := g.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
		return &cached, nil
	}

	user, err := g.inner.GetUserByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	_ = g.cache.Set(ctx, cacheKey, user, int(g.ttl.Seconds()))
	return user, nil
}

func (g *CachedUserGateway) GetCurrentUser(ctx context.Context) (*gw_response.CurrentUser, error) {
	return g.inner.GetCurrentUser(ctx)
}
