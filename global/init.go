package global

import (
	"fmt"

	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"github.com/beewit/beekit/utils/convert"
)

var (
	CFG        = conf.New("config.json")
	Log        = log.Logger
	DB         = mysql.DB
	RD         = redis.Cache
	IP         = CFG.Get("server.ip")
	Port       = CFG.Get("server.port")
	SortDoMain = convert.ToString(CFG.Get("sortUrl.domain"))
	Host       = fmt.Sprintf("http://%v:%v", IP, Port)
)
