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
	"github.com/beewit/beekit/utils/user_agent"
	"github.com/beewit/beekit/utils/convert"
)

func TestSortUrl(t *testing.T) {
	if !utils.IsUrl("https://www.baidu.com/?a=%E8%B4%A6%E5%8F%B7") {
		println("非正常域名")
	} else {
		println("正确域名")
	}

}

func TestStrLen(t *testing.T) {
	println(len("leng"))
	println(len("技术长度"))
	println(len("leng."))
	println(len("技术长度!"))
	println(len("leng。"))
	println(len("技术长度！"))
	println(len("leng技术长度"))
	println(ShowSubstr("leng技术长度", 6))
	println(ShowSubstr("盘龙奥园林芊儿童棉麻生活馆", 10))
	println(utils.SubStrByByteInChar("盘龙奥园林芊儿童棉麻生活馆", 32))
	println(SubstrByByte("盘龙奥园林芊儿童棉麻生活馆", 30)+"..")
	println(len("盘龙奥园林芊儿童棉麻.."))
}

func SubstrByByte(str string, length int) string {
	bs := []byte(str)[:length]
	bl := 0
	for i:=len(bs)-1; i>=0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++;
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i] & 252 == 252:
				cl = 6
			case bs[i] & 248 == 248:
				cl = 5
			case bs[i] & 240 == 240:
				cl = 4
			case bs[i] & 224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}

func show_substr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl + rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func ShowSubstr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl + rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
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
	rp, err := ApiPost("http://sort.9ee3.com/api/create?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.HV8iAqo9biUNtnhxoTMT4KCC4PL60NMpHSaq7PqtZCc&longUrl=http://www.baidu.com", nil)
	if err != nil {
		t.Error(err.Error())
	}
	str, err2 := json.Marshal(rp)
	if err2 != nil {
		t.Error(err2.Error())
	}
	println(string(str))
}




func TestUserAgent(t *testing.T) {
	ua := user_agent.New("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.11")

	println(convert.MustJsonPrettyString(ua))
}


