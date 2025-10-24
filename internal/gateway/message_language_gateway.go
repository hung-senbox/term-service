package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"term-service/internal/gateway/dto/request"
	"term-service/internal/gateway/dto/response"
	"term-service/pkg/constants"
	"term-service/pkg/helper"

	"github.com/hashicorp/consul/api"
)

type MessageLanguageGateway interface {
	UploadMessage(ctx context.Context, req request.UploadMessageRequest) error
	UploadMessages(ctx context.Context, req request.UploadMessageLanguagesRequest) error
	GetMessageLanguages(ctx context.Context, typeStr string, typeID string) ([]response.MessageLanguageResponse, error)
	GetMessageLanguage(ctx context.Context, typeStr string, typeID string) (response.MessageLanguageResponse, error)
	DeleleByTypeAndTypeID(ctx context.Context, typeStr string, typeID string) error
}

type messageLanguageGateway struct {
	serviceName string
	consul      *api.Client
}

func NewMessageLanguageGateway(serviceName string, consulClient *api.Client) MessageLanguageGateway {
	return &messageLanguageGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}
func (g *messageLanguageGateway) UploadMessage(ctx context.Context, req request.UploadMessageRequest) error {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	headers := helper.GetHeaders(ctx)

	_, err = client.Call("POST", "/v1/gateway/messages", req, headers)
	if err != nil {
		return err
	}

	return nil
}

func (g *messageLanguageGateway) UploadMessages(ctx context.Context, req request.UploadMessageLanguagesRequest) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	headers := helper.GetHeaders(ctx)

	_, err = client.Call("POST", "/v1/gateway/messages", req, headers)
	if err != nil {
		return err
	}

	return nil
}

func (g *messageLanguageGateway) GetMessageLanguages(ctx context.Context, typeStr string, typeID string) ([]response.MessageLanguageResponse, error) {
	// lấy token từ context
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	// tạo client
	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	headers := helper.GetHeaders(ctx)

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages?type=%s&type_id=%s", typeStr, typeID)
	resp, err := client.Call("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp response.APIGateWayResponse[[]response.MessageLanguageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get message languages fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *messageLanguageGateway) GetMessageLanguage(ctx context.Context, typeStr string, typeID string) (response.MessageLanguageResponse, error) {
	// lấy token từ context
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return response.MessageLanguageResponse{}, fmt.Errorf("token not found in context")
	}

	appLanguage, ok := ctx.Value(constants.AppLanguage).(uint)
	if !ok {
		return response.MessageLanguageResponse{}, fmt.Errorf("app language not found in context")
	}

	// tạo client
	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return response.MessageLanguageResponse{}, err
	}

	headers := helper.GetHeaders(ctx)

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages/get-by-language?type=%s&type_id=%s&language_id=%d", typeStr, typeID, appLanguage)
	resp, err := client.Call("GET", url, nil, headers)
	if err != nil {
		return response.MessageLanguageResponse{}, err
	}

	// parse JSON
	var gwResp response.APIGateWayResponse[response.MessageLanguageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return response.MessageLanguageResponse{}, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return response.MessageLanguageResponse{}, fmt.Errorf("call gateway get message language fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *messageLanguageGateway) DeleleByTypeAndTypeID(ctx context.Context, typeStr string, typeID string) error {
	// lấy token từ context
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	// tạo client
	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	headers := helper.GetHeaders(ctx)

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages?type=%s&type_id=%s", typeStr, typeID)
	_, err = client.Call("DELETE", url, nil, headers)
	if err != nil {
		return err
	}

	return nil
}
