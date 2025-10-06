package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"term-service/internal/gateway/dto"
	"term-service/logger"
	"term-service/pkg/constants"
	"term-service/pkg/helper"

	"github.com/hashicorp/consul/api"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserGateway interface {
	GetAuthorInfo(ctx context.Context, userID string) (*User, error)
	GetCurrentUser(ctx context.Context) (*dto.CurrentUser, error)
	GetStudentInfo(ctx context.Context, studentID string) (*dto.StudentResponse, error)
}

type userGatewayImpl struct {
	serviceName string
	consul      *api.Client
}

func NewUserGateway(serviceName string, consulClient *api.Client) UserGateway {
	return &userGatewayImpl{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

// GetAuthorInfo lấy thông tin user từ service user
func (g *userGatewayImpl) GetAuthorInfo(ctx context.Context, userID string) (*User, error) {
	token, ok := ctx.Value("token").(string) // hoặc dùng constants.TokenKey
	if !ok || token == "" {
		return nil, fmt.Errorf("token not exist context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("Call API user fail: %w", err)
	}

	var user User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, fmt.Errorf("encrypt response fail: %w", err)
	}

	return &user, nil
}

// GetCurrentUser
func (g *userGatewayImpl) GetCurrentUser(ctx context.Context) (*dto.CurrentUser, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		logger.WriteLogEx("warn", "token not found in context", nil)
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		logger.WriteLogEx("error", "init GatewayClient fail", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/user/current-user", nil, headers)
	if err != nil {
		logger.WriteLogEx("error", "call API user fail", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		logger.WriteLogEx("error", "unmarshal response fail", map[string]any{
			"error": string(resp),
		})
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		logger.WriteLogEx("warn", "gateway error", map[string]any{
			"status_code": gwResp.StatusCode,
			"message":     gwResp.Message,
		})
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetStudentInfo(ctx context.Context, studentID string) (*dto.StudentResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/students/"+studentID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API student fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.StudentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
