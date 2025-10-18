package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"l4d2serverquery-go/dto"
	"l4d2serverquery-go/logger"
	"l4d2serverquery-go/pkg/steamquery"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rumblefrog/go-a2s"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"

	"l4d2serverquery-go/service"
)

func Router() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(cors.Default())

	setUI(r)

	// api 接口

	r.GET("health", func(c *gin.Context) {
		c.String(http.StatusOK, "health check")
	})

	{
		r.GET("tags", func(c *gin.Context) {
			tags, err := service.Tags()
			if err != nil {
				c.String(500, err.Error())
				return
			}
			c.JSON(http.StatusOK, tags)
		})

		r.POST("tags", func(c *gin.Context) {
			type CreateTagRequest struct {
				Name string `json:"name" binding:"required"`
				Rank int    `json:"rank"`
			}

			var req CreateTagRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			tag, err := service.CreateTag(req.Name, req.Rank)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, dto.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			})
		})

		r.DELETE("tags/:id", func(c *gin.Context) {
			id := cast.ToInt(c.Param("id"))

			err := service.DeleteTag(id)
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
				return
			}

			c.Status(http.StatusOK)
		})

		// 绑定服务器标签
		r.POST("server/:serverId/tags", func(c *gin.Context) {
			serverID := cast.ToInt(c.Param("serverId"))

			type BindTagsRequest struct {
				TagIDs []int `json:"tagIds" binding:"required"`
			}

			var req BindTagsRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			err := service.BindServerTags(serverID, req.TagIDs)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.Status(http.StatusOK)
		})

		// 获取服务器标签
		r.GET("server/:serverId/tags", func(c *gin.Context) {
			serverID := cast.ToInt(c.Param("serverId"))

			tags, err := service.GetServerTags(serverID)
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
				return
			}

			c.JSON(http.StatusOK, tags)
		})
	}

	// 获取数据库中的所有服务器信息并排序
	r.POST("serverList/v2", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		tagIDRes := gjson.GetBytes(body, "tags").Array()
		name := gjson.GetBytes(body, "name").String()
		minPlayers := int(gjson.GetBytes(body, "minPlayers").Int())
		maxPlayers := int(gjson.GetBytes(body, "maxPlayers").Int())

		tagIDs := slice.Map(
			tagIDRes,
			func(_ int, item gjson.Result) int {
				return int(item.Int())
			},
		)
		servers, err := service.QueryServers(tagIDs, name, minPlayers, maxPlayers)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(http.StatusOK, servers)
	})

	r.GET("/singleServer/:id", func(c *gin.Context) {
		id := cast.ToInt(c.Param("id"))

		resp, err := service.QuerySingleServer(id)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	// 更新服务器的最后连接时间
	r.GET("lastCopyTimeUpdate/:id", func(c *gin.Context) {
		id := cast.ToInt(c.Param("id"))
		if err := service.UpdateLastCopyTime(id); err != nil {
			c.String(500, err.Error())
			return
		}

		c.Status(200)
	})

	r.POST("/serverPatch/:id", func(c *gin.Context) {
		id := cast.ToInt(c.Param("id"))

		var item dto.Server
		if err := c.ShouldBindJSON(&item); err != nil {
			c.String(404, err.Error())
			return
		}

		err := service.PatchServer(id, item)
		if err != nil {
			c.Status(404)
			return
		}
		c.Status(200)
	})

	r.DELETE("/serverDelete/:id", func(c *gin.Context) {
		id := cast.ToInt(c.Param("id"))
		err := service.DeleteServerById(id)
		if err != nil {
			c.Status(404)
			return
		}
		c.Status(200)
	})

	r.GET("/groupByTag", func(c *gin.Context) {
		logger.Log.Info("开始给服务器分类")
		if err := service.GroupServers(); err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.Status(200)
	})

	r.GET("/debug/cleanServers", func(c *gin.Context) {
		if err := service.CleanServers(); err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.Status(200)
	})

	// 通过 https://www.steamserverbrowser.com/games/left-4-dead-2/asia
	// 来查询可用的服务器信息
	r.GET("/master_query", func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", "10")
		serverName := c.DefaultQuery("serverName", "药抗")

		servers := steamquery.QuerySteamServerBrowswerApi(
			serverName,
			cast.ToInt(page),
			cast.ToInt(pageSize),
		)

		masterQueryResult := map[string]any{
			"查询时间":       time.Now().Format(time.DateTime),
			"本次查询的服务器名称": serverName,
			"本次查询的数量":    len(servers),
		}

		pretty, _ := formatter.Pretty(masterQueryResult)
		fmt.Println(pretty)

		servers = service.RemoveDuplicateServers(servers)
		c.JSON(http.StatusOK, servers)
	})

	r.GET("/serverAddByMasterQuery", func(c *gin.Context) {
		addr := c.Query("addr")

		fmt.Println("添加服务器", addr)

		err := service.AddServer(addr)
		if err != nil {
			log.Println("添加服务器失败", err)
			c.Status(400)
			return
		}

		c.Status(200)
	})

	// 获取玩家名称, 暂时只有玩家名称
	r.POST("/players/query", func(c *gin.Context) {
		type Player struct {
			Addr string `json:"addr"`
		}

		var p Player
		if err := c.BindJSON(&p); err != nil {
			return
		}

		address := "192.168.127.12:27015"
		address = p.Addr

		client, err := a2s.NewClient(address)
		if err != nil {
			c.String(400, fmt.Errorf("连接服务器失败: %s", err).Error())
			fmt.Println("查询玩家: 连接服务器失败", address, err)
			return
		}

		player, err := client.QueryPlayer()
		if err != nil {
			c.String(400, fmt.Errorf("查询玩家失败: %s", err).Error())
			fmt.Println("查询玩家: 查询玩家失败", address, err)
			return
		}

		type PlayerDto struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Score   int    `json:"score"`
			Seconds int    `json:"seconds"`
		}

		playerNames := slice.Map(
			player.Players,
			func(index int, item *a2s.Player) PlayerDto {
				return PlayerDto{Name: item.Name}
			},
		)

		c.JSON(200, playerNames)
	})

	// 一个用来测试的ws接口
	r.GET("/ws", func(c *gin.Context) {

	})

	return r
}
