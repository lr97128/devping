package main

import (
        "net/http"
        "strconv"
        "log"

        "github.com/gin-gonic/gin"
)

func getPing(c *gin.Context) {
        target := c.Query("target")
        countstr := c.Query("count")
        count, err := strconv.Atoi(countstr)
        if err != nil {
                log.Fatal(err)
        }
        result, err := GetPing(target, count)
        if err != nil {
                log.Fatal(err)
        }
	c.String(http.StatusOK, result)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
        r := gin.Default()
        r.GET("/metrics", getPing)
        r.Run(":9118")
}
