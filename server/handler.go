package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/server/log"
)

func (s *Server) Handler() (handler http.Handler) {
	sysRoutes(s.Router)

	handler = s
	if s.HandleTimeout > 0 {
		handler = http.TimeoutHandler(handler, s.HandleTimeout,
			fmt.Sprintf(`ServeHTTP timeout after %s.`, s.HandleTimeout),
		)
	}
	return
}

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	psData.Add(request.Method, request.URL.Path, startTime)
	defer psData.Remove(request.Method, request.URL.Path, startTime)

	req := xiaomei.NewRequest(request, s.Session)
	res := xiaomei.NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

	var notFound bool
	defer handleError(startTime, req, res, &notFound)

	// 如果返回true，继续交给路由处理
	if strings.HasPrefix(req.URL.Path, `/_`) || s.FilterFunc == nil || s.FilterFunc(req, res) {
		notFound = !s.Router.Handle(req, res)
	}
}

func handleError(t time.Time, req *xiaomei.Request, res *xiaomei.Response, notFound *bool) {
	if *notFound {
		handleNotFound(req, res)
	}

	err := recover()
	if err != nil {
		handleServerError(req, res)
	}
	if err == nil && strings.HasPrefix(req.URL.Path, `/_`) {
		return
	}
	log.Write(req, res, t, err)
}

func handleNotFound(req *xiaomei.Request, res *xiaomei.Response) {
	res.WriteHeader(404)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "404", "message": "Not Found."})
	}
}

func handleServerError(req *xiaomei.Request, res *xiaomei.Response) {
	res.WriteHeader(500)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "500", "message": "Application Server Error."})
	}
}
