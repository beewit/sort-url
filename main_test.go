package main

import (
	"testing"
	"github.com/beewit/beekit/utils"
	"fmt"
	"crypto/md5"
	"net/http/httptest"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo"
	"github.com/beewit/sort-url/handle"
	"net/http"
	"net/url"
	"github.com/beewit/beekit/utils/uhttp"
	"encoding/json"
)

func TestSortUrl(t *testing.T) {
	if !utils.IsUrl("https://www.baidu.com/?a=%E8%B4%A6%E5%8F%B7") {
		println("非正常域名")
	} else {
		println("正确域名")
	}

}

func TestHashCode(t *testing.T) {
	urlMd5 := fmt.Sprintf("%x", md5.Sum([]byte("https://www.baidu.com/s?wd=%E9%A1%BA%E4%B8%B0%E5%BF%AB%E9%80%92%E5%8D%95%E5%8F%B7%E6%9F%A5%E8%AF%A2&rsv_spt=1&rsv_iqid=0xac6cf48a000154fa&issp=1&f=3&rsv_bp=0&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_sug3=3&rsv_sug1=3&rsv_sug7=101&rsv_t=7852rbLbUY9f3gGNKQ9nmhc6y0OF6RhoIJ%2FHX8prwKjWSU5B7Ko7tvo9K7uMNEB44sZx&rsv_sug2=1&prefixsug=sf&rsp=1&inputT=1846&rsv_sug4=2087")))
	println(urlMd5)

}

func TestCreateSortUrl(t *testing.T) {
	e := echo.New()
	f := url.Values{}
	f.Set("longUrl", "http://sso.9ee3.com")
	f.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.XMvP2ISedrRFJH9N0-G-YIACkXsO49ZdcbNKtS8GdO8872")
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// 断言
	if assert.NoError(t, handle.CreateSortUrl(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t.Log(rec.Body.String())
	}
}

func ApiPost(url string, m map[string]string) (utils.ResultParam, error) {
	b, _ := json.Marshal(m)
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    url,
		Body:   b,
	})
	if err != nil {
		return utils.ResultParam{}, err
	}
	return utils.ToResultParam(body), nil
}

func TestCreateSortUrlApi(t *testing.T) {
	rp, err := ApiPost("http://127.0.0.1:8085/api/create?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.YE5cN3CIAq1Bk24hEU0euHz0tGcwW_NUhOaBy7diwF0", nil)
	if err != nil {
		t.Error(err.Error())
	}
	str, err2 := json.Marshal(rp)
	if err2 != nil {
		t.Error(err2.Error())
	}
	println(string(str))
}
