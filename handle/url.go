package handle

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/sort-url/global"
	"github.com/beewit/beekit/utils/convert"
	"fmt"
	"crypto/md5"
	"github.com/pkg/errors"
	"net/http"
)

//创建短链接
func CreateSortUrl(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	var longUrl = c.FormValue("longUrl")
	if longUrl == "" {
		return utils.ErrorNull(c, "转换失败，原地址为空")
	}
	if len(longUrl) >= 1024 {
		return utils.ErrorNull(c, "转换失败，原地址为超过1024字符")
	}
	//是否为正常url判断
	if !utils.IsUrl(longUrl) {
		return utils.ErrorNull(c, "转换失败，原地址不是有效的Url，格式示例：http://www.baidu.com")
	}
	//判断数据库是否存在
	m, err := getByLongUrl(longUrl)
	if err != nil {
		return utils.ErrorNull(c, "转换失败")
	}
	if m != nil {
		return utils.Success(c, "转换成功", getSortUrl(convert.ToString(m["sort_url"])))
	}
	//添加短链接
	sortUrl, err := addSortUrl(longUrl, c.RealIP(), acc.ID)
	if err != nil {
		global.Log.Error("CreateSortUrl error：%v", err.Error())
		return utils.ErrorNull(c, "转换创建失败")
	}
	if sortUrl == "" {
		return utils.ErrorNull(c, "转换创建失败")
	}
	return utils.Success(c, "转换短链接成功", getSortUrl(sortUrl))
}

//短链接跳转
func SortUrlJump(c echo.Context) error {
	sortUrl := c.ParamValues()[0]
	m, err := getBySortUrl(sortUrl)
	if err != nil {
		return utils.ResultString(c, "短链接无效")
	}
	if m == nil {
		return utils.ResultString(c, "短链接无效")
	}
	head := c.Request().Header
	go addSortUrlLog(m, head, c.RealIP())
	return c.Redirect(http.StatusFound, convert.ToString(m["origin_url"]))
}

//Referer
func addSortUrlLog(m map[string]interface{}, head http.Header, ip string) {
	global.DB.InsertMap("short_link_logs", map[string]interface{}{
		"id":           utils.ID(),
		"sort_link_id": m["id"],
		"ct_time":      utils.CurrentTime(),
		"browser":      head.Get("User-Agent"),
		"referer":      head.Get("Referer"),
		"ip":           ip,
		"remark":       convert.ToObjStr(head),
	})
}

func getSortUrl(sortUrl string) string {
	return global.SortDoMain + "/" + sortUrl
}

func addSortUrl(longUrl, ip string, accId int64) (string, error) {
	var flog = false
	var sortUrl string
	var e error
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		id, err := tx.InsertMap("short_link", map[string]interface{}{
			"hash_code":  getHashCode(longUrl),
			"origin_url": longUrl,
			"ip":         ip,
			"ct_time":    utils.CurrentTime(),
			"account_id": accId,
		})
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if id <= 0 {
			panic(errors.New("保存短链接失败"))
		}
		sortUrl = convert.TransTo62(id)
		x, err := tx.Update("UPDATE short_link SET sort_url=? WHERE id=?", sortUrl, id)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			panic(errors.New("保存短链接失败"))
		}
		flog = true
	}, func(err error) {
		if err != nil {
			global.Log.Error("addSortUrl，error：%s", err.Error())
			flog = false
			e = err
		}
	})
	if flog {
		return sortUrl, nil
	}
	return "", e
}

func getBySortUrl(sortUrl string) (map[string]interface{}, error) {
	//判断数据库是否存在此链接地址
	rows, err := global.DB.Query("SELECT * FROM short_link WHERE sort_url=? LIMIT 1", sortUrl)
	if err != nil {
		global.Log.Error("getBySortUrl sql error：：%s", err.Error())
		return nil, err
	}
	if len(rows) != 1 {
		return nil, nil
	}
	return rows[0], nil
}

func getByLongUrl(longUrl string) (map[string]interface{}, error) {
	//判断数据库是否存在此链接地址
	rows, err := global.DB.Query("SELECT * FROM short_link WHERE hash_code=? LIMIT 1", getHashCode(longUrl))
	if err != nil {
		global.Log.Error("getByLongUrl sql error：%s", err.Error())
		return nil, err
	}
	if len(rows) != 1 {
		return nil, nil
	}
	return rows[0], nil
}

func getHashCode(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
