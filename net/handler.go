package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/code560/audigo/player"
	"github.com/code560/audigo/util"
	"github.com/gin-gonic/gin"
)

var (
	handle = newHandler()
	log    = util.GetLogger()
)

const (
	INIT_PLAYER_COUNT = 25
)

type handler struct {
	players map[string]player.Proxy
}

// SetHandler は、ginにapiハンドラーを設定します。
func SetHandler(r *gin.Engine) {
	setV1(r)
}

func setV1(r *gin.Engine) {
	v1 := r.Group("audio/v1")
	{
		v1.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
		v1.GET("/players", handle.ids)
		v1.POST("/init/:content_id", handle.create)
		v1.POST("/play/:content_id", handle.play)
		v1.POST("/stop/:content_id", handle.stop)
		v1.POST("/volume/:content_id", handle.volume)
		v1.POST("/pause/:content_id", handle.pause)
		v1.POST("/resume/:content_id", handle.resume)
	}
}

func newHandler() *handler {
	var inst *handler
	inst = &handler{
		players: make(map[string]player.Proxy, INIT_PLAYER_COUNT),
	}
	return inst
}

func (h *handler) create(c *gin.Context) {
	log.Info("call init rest-api audio module.\n", c)
	code := http.StatusNoContent
	h.getPlayer(c.Param("content_id"), true)
	code = http.StatusCreated
	c.JSON(code, nil)
}

func (h *handler) getPlayer(id string, create bool) (player.Proxy, error) {
	p, ok := h.players[id]
	if !ok {
		if create {
			p = player.NewProxy()
			h.players[id] = p
		} else {
			return nil, fmt.Errorf("not found id player: %s", id)
		}
	}
	return p, nil
}

func (h *handler) ids(c *gin.Context) {
	code := http.StatusBadRequest
	keys := make([]string, len(h.players))
	for k := range h.players {
		keys = append(keys, k)
	}
	res, err := json.Marshal(keys)
	if err != nil {
		code = http.StatusInternalServerError
		c.JSON(code, err)
		return
	}
	c.JSON(code, res)
}

func (h *handler) play(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Info("call play rest-api audio module.\n", c.Request.Header, "\n", string(body))
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	// create args
	code := http.StatusAccepted
	p, _ := h.getPlayer(c.Param("content_id"), true)
	var args player.PlayArgs
	if err := c.ShouldBindJSON(&args); err != nil {
		log.Error("Json binded error: ", err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	// send
	select {
	case p.GetChannel() <- &player.Action{Act: player.Play, Args: &args}:
		break
	default:
		log.Error("dont send player chan: play")
	}
	c.JSON(code, nil)
}

func (h *handler) stop(c *gin.Context) {
	log.Info("call stop rest-api audio module.\n", c.Request.Header)

	code := http.StatusAccepted
	p, err := h.getPlayer(c.Param("content_id"), false)
	if err != nil {
		return
	}
	// send
	select {
	case p.GetChannel() <- &player.Action{Act: player.Stop, Args: struct{}{}}:
		break
	default:
		log.Error("dont send player chan: stop")
	}
	c.JSON(code, nil)
}

func (h *handler) volume(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Info("call play rest-api audio module.\n", c.Request.Header, "\n", string(body))
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	code := http.StatusAccepted
	p, err := h.getPlayer(c.Param("content_id"), true)
	if err != nil {
		return
	}
	var args player.VolumeArgs
	if err := c.ShouldBindJSON(&args); err != nil {
		log.Error("Json binded error: ", err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	// send
	select {
	case p.GetChannel() <- &player.Action{Act: player.Volume, Args: &args}:
		break
	default:
		log.Error("dont send player chan: stop")
	}
	c.JSON(code, nil)
}

func (h *handler) pause(c *gin.Context) {
	log.Info("call pause rest-api audio module.\n", c.Request.Header)
	code := http.StatusAccepted
	p, err := h.getPlayer(c.Param("content_id"), true)
	if err != nil {
		return
	}
	select {
	case p.GetChannel() <- &player.Action{Act: player.Pause, Args: struct{}{}}:
		break
	default:
		log.Error("dont send player chan: stop")
	}
	c.JSON(code, nil)
}

func (h *handler) resume(c *gin.Context) {
	log.Info("call resume rest-api audio module.\n", c.Request.Header)
	code := http.StatusAccepted
	p, err := h.getPlayer(c.Param("content_id"), true)
	if err != nil {
		return
	}
	select {
	case p.GetChannel() <- &player.Action{Act: player.Resume, Args: struct{}{}}:
		break
	default:
		log.Error("dont send player chan: stop")
	}
	c.JSON(code, nil)
}
