package cls

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

// Client wraps the Tencent Cloud CLS API client.
type Client struct {
	api *cls.Client
}

// NewClient creates a CLS client with the given credentials and region.
func NewClient(secretID, secretKey, region string) (*Client, error) {
	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	api, err := cls.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	return &Client{api: api}, nil
}

// API returns the underlying SDK client for direct use when needed.
func (c *Client) API() *cls.Client {
	return c.api
}
