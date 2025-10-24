package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"term-service/internal/gateway/dto/response"
	"term-service/pkg/constants"
	"term-service/pkg/helper"

	"github.com/hashicorp/consul/api"
)

type OrganizationGateway interface {
	GetOrganizationInfo(ctx context.Context, organizationID string) (*response.OrganizationInfo, error)
	GetAllOrg(ctx context.Context) ([]response.OrganizationInfo, error)
}

type organizationGatewayImpl struct {
	serviceName string
	consul      *api.Client
}

func NewOrganizationGateway(serviceName string, consulClient *api.Client) OrganizationGateway {
	return &organizationGatewayImpl{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *organizationGatewayImpl) GetOrganizationInfo(ctx context.Context, organizationID string) (*response.OrganizationInfo, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/organization/"+organizationID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API get info organization fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.OrganizationInfo]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *organizationGatewayImpl) GetAllOrg(ctx context.Context) ([]response.OrganizationInfo, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/organizations", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API get all organization fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[[]response.OrganizationInfo]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}
