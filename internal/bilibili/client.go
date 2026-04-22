package bilibili

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	ErrSessionExpired = errors.New("bilibili session expired")
	ErrAuthFailed    = errors.New("bilibili authentication failed")
	ErrNetworkError  = errors.New("network error")
)

type Client struct {
	sessdata string
	biliJct  string
	buvid3   string
	client   *resty.Client
}

type Video struct {
	Title    string `json:"title"`
	Aid      int64  `json:"aid"`
	Bvid      string `json:"bvid"`
	Pic       string `json:"pic"`
	Duration  int    `json:"duration"`
	Author    string `json:"author"`
	ShortLink string `json:"short_link_v2"`
}

type ToViewResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []Video `json:"list"`
	} `json:"data"`
}

func NewClient(sessdata, biliJct, buvid3 string) *Client {
	return &Client{
		sessdata: sessdata,
		biliJct:  biliJct,
		buvid3:   buvid3,
		client: resty.New().
			SetRetryCount(3).
			SetRetryWaitTime(1 * time.Second).
			SetRetryMaxWaitTime(4 * time.Second),
	}
}

func (c *Client) GetToViewList() ([]Video, error) {
	resp, err := c.client.R().
		SetHeader("Cookie", c.buildCookie()).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get("https://api.bilibili.com/x/v2/history/toview")
	if err != nil {
		return nil, ErrNetworkError
	}

	var result ToViewResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	switch result.Code {
	case 0:
		return result.Data.List, nil
	case -101:
		return nil, ErrSessionExpired
	case -111:
		return nil, ErrAuthFailed
	default:
		return nil, errors.New(result.Message)
	}
}

func (c *Client) buildCookie() string {
	return "SESSDATA=" + c.sessdata + "; bili_jct=" + c.biliJct + "; buvid3=" + c.buvid3
}
