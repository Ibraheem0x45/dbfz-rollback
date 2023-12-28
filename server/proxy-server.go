package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrInvalidPayloadSize = errors.New("invalid payload size")

var ENV_PAYLOAD = []byte{0x92, 0x98, 0xAD, 0x36, 0x35, 0x37, 0x64, 0x65, 0x38, 0x34, 0x30, 0x61, 0x33, 0x31, 0x36, 0x32, 0xB3, 0x32, 0x30, 0x32, 0x33, 0x2F, 0x31, 0x32, 0x2F, 0x31, 0x37, 0x20, 0x30, 0x33, 0x3A, 0x31, 0x31, 0x3A, 0x31, 0x32, 0xA5, 0x30, 0x2E, 0x30, 0x2E, 0x33, 0xA5, 0x30, 0x2E, 0x30, 0x2E, 0x33, 0xA5, 0x30, 0x2E, 0x30, 0x2E, 0x33, 0xA0, 0xA0, 0x92, 0xBE, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F, 0x73, 0x74, 0x2F, 0x61, 0x70, 0x69, 0x2F}

const DBFZ_API_URL = "https://dbf.channel.or.jp/api/"

type Server struct {
	Engine *gin.Engine
	Client *http.Client
}

func New() *Server {
	gin.SetMode(gin.ReleaseMode)

	return &Server{
		Engine: gin.Default(),
		Client: &http.Client{},
	}
}

func (server *Server) Run(address string) error {
	server.Engine.POST("/api/sys/get_env", server.GetEnvHandler)
	server.Engine.POST("/api/replay/data_save", func(ctx *gin.Context) { ctx.Status(http.StatusNotFound) })
	server.Engine.NoRoute(server.RequestHandler)

	return server.Engine.Run(address)
}

func (server *Server) GetEnvHandler(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=UTF-8", ENV_PAYLOAD)
}

func (server *Server) RequestBodyHandler(ctx *gin.Context, body []byte) ([]byte, error) {
	uri := ctx.Request.RequestURI

	if strings.Contains(uri, "user/login") {
		if len(body) < 32 {
			return nil, ErrInvalidPayloadSize
		}

		data := append([]byte{}, body[:27]...)
		data = append(data, 0x30, 0x33)
		data = append(data, body[31:]...)
		return data, nil
	}

	if len(body) < 94 {
		return nil, ErrInvalidPayloadSize
	}

	data := append([]byte{}, body[:89]...)
	data = append(data, 0x30, 0x33)
	data = append(data, body[93:]...)
	return data, nil
}

func (server *Server) RequestHandler(ctx *gin.Context) {
	request := ctx.Request

	request_body, err := io.ReadAll(request.Body)

	if err != nil {
		log.Printf("Failed to read request body: %s\n", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	modified_body, err := server.RequestBodyHandler(ctx, request_body)

	if err != nil {
		log.Printf("Failed to modify request body: %s\n", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	response, err := server.MakeRequest(ctx, modified_body)

	if err != nil {
		log.Printf("Failed to make request: %s\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	ctx.Status(response.StatusCode)

	for name, headers := range response.Header {
		ctx.Writer.Header()[name] = headers
	}

	_, err = io.Copy(ctx.Writer, response.Body)

	if err != nil {
		log.Printf("Failed to copy response body: %s\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (server *Server) MakeRequest(ctx *gin.Context, payload []byte) (*http.Response, error) {
	request := ctx.Request

	api, err := url.Parse(DBFZ_API_URL)

	if err != nil {
		return nil, err
	}

	api.Path = request.URL.Path

	request.URL = api
	request.Host = ""
	request.RequestURI = ""

	request.Body = io.NopCloser(bytes.NewReader(payload))
	request.ContentLength = int64(len(payload))
	return server.Client.Do(request)
}
