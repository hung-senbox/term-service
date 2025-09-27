package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"term-service/internal/gateway/dto"
	"term-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type MessageLanguageGateway interface {
	UploadMessage(ctx context.Context, req dto.UploadMessageRequest) error
	UploadMessages(ctx context.Context, req dto.UploadMessageLanguagesRequest) error
	GetMessageLanguages(ctx context.Context, typeStr string, typeID string) ([]dto.MessageLanguageResponse, error)
	GetMessageLanguage(ctx context.Context, typeStr string, typeID string) (dto.MessageLanguageResponse, error)
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
func (g *messageLanguageGateway) UploadMessage(ctx context.Context, req dto.UploadMessageRequest) error {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	_, err = client.Call("POST", "/v1/gateway/messages", req)
	if err != nil {
		return err
	}

	return nil
}

func (g *messageLanguageGateway) UploadMessages(ctx context.Context, req dto.UploadMessageLanguagesRequest) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	_, err = client.Call("POST", "/v1/gateway/messages", req)
	if err != nil {
		return err
	}

	return nil
}

func (g *messageLanguageGateway) GetMessageLanguages(ctx context.Context, typeStr string, typeID string) ([]dto.MessageLanguageResponse, error) {
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

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages?type=%s&type_id=%s", typeStr, typeID)
	resp, err := client.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp dto.APIGateWayResponse[[]dto.MessageLanguageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get message languages fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *messageLanguageGateway) GetMessageLanguage(ctx context.Context, typeStr string, typeID string) (dto.MessageLanguageResponse, error) {
	// lấy token từ context
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return dto.MessageLanguageResponse{}, fmt.Errorf("token not found in context")
	}

	appLanguage, ok := ctx.Value(constants.AppLanguage).(uint)
	if !ok {
		return dto.MessageLanguageResponse{}, fmt.Errorf("app language not found in context")
	}

	// tạo client
	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return dto.MessageLanguageResponse{}, err
	}

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages/get-by-language?type=%s&type_id=%s&language_id=%d", typeStr, typeID, appLanguage)
	resp, err := client.Call("GET", url, nil)
	if err != nil {
		return dto.MessageLanguageResponse{}, err
	}

	// parse JSON
	var gwResp dto.APIGateWayResponse[dto.MessageLanguageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return dto.MessageLanguageResponse{}, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return dto.MessageLanguageResponse{}, fmt.Errorf("call gateway get message language fail: %s", gwResp.Message)
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

	// gọi API với query params
	url := fmt.Sprintf("/v1/gateway/messages?type=%s&type_id=%s", typeStr, typeID)
	_, err = client.Call("DELETE", url, nil)
	if err != nil {
		return err
	}

	return nil
}
