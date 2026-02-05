package cls

import (
	"context"
	"strconv"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

// GetContextInput holds parameters for DescribeLogContext.
type GetContextInput struct {
	TopicId  string
	PkgId    string
	PkgLogId int64
}

// GetContext runs DescribeLogContext and returns the context logs.
func (c *Client) GetContext(ctx context.Context, in GetContextInput) ([]*cls.LogContextInfo, error) {
	req := cls.NewDescribeLogContextRequest()
	req.TopicId = &in.TopicId
	req.PkgId = &in.PkgId
	req.PkgLogId = &in.PkgLogId

	resp, err := c.api.DescribeLogContextWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Response == nil || resp.Response.LogContextInfos == nil {
		return nil, nil
	}
	return resp.Response.LogContextInfos, nil
}

// PkgLogIdFromString parses PkgLogId from string (e.g. "65536").
func PkgLogIdFromString(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
