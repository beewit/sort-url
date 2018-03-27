package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/beewit/beekit/utils/convert"
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/sort-url/global"
	"github.com/beewit/sort-url/handle"
)

func main() {
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.File("/", "app/page/index.html")
	e.POST("/api/create", handle.CreateSortUrl, handle.Filter)
	e.POST("/api/get/page", handle.GetSortUrlPage, handle.Filter)
	e.GET("*", handle.SortUrlJump)
	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))
}
