package gateway

import (
	"encoding/json"
	"fmt"
	gatewaydto "term-service/internal/term/dto/gateway_dto"
	"term-service/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserGateway interface {
	GetAuthorInfo(ctx *gin.Context, userID string) (*User, error)
	GetCurrentUser(ctx *gin.Context) (*gatewaydto.CurrentUser, error)
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
func (g *userGatewayImpl) GetAuthorInfo(ctx *gin.Context, userID string) (*User, error) {
	token, ok := ctx.Value("token").(string) // hoặc dùng constants.TokenKey
	if !ok || token == "" {
		return nil, fmt.Errorf("token not exist context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/user/"+userID, nil)
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
func (g *userGatewayImpl) GetCurrentUser(ctx *gin.Context) (*gatewaydto.CurrentUser, error) {
	tokenValue, exists := ctx.Get(constants.Token)
	if !exists {
		return nil, fmt.Errorf("token not exist in gin context")
	}
	token := tokenValue.(string)

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/user/current-user", nil)
	if err != nil {
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp gatewaydto.APIGateWayResponse[gatewaydto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
