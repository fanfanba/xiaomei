package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/renderer"
	"github.com/lovego/xiaomei/router"
	"github.com/lovego/xiaomei/server/funcs"
	"github.com/lovego/xiaomei/session"
)

func init() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}
}

type Server struct {
	*http.Server
	HandleTimeout  time.Duration
	FilterFunc     func(req *xiaomei.Request, res *xiaomei.Response) bool
	Router         *router.Router
	Session        session.Session
	Renderer       *renderer.Renderer
	LayoutDataFunc func(
		layout string, data interface{}, req *xiaomei.Request, res *xiaomei.Response,
	) interface{}
}

func NewSession() session.Session {
	return session.NewCookieSession(http.Cookie{
		Name: config.Name(),
		Path: `/`,
	}, config.Secret())
}

func NewRenderer() *renderer.Renderer {
	return renderer.New(
		path.Join(config.Root(), `views`), `layout/default`, !config.DevMode(), funcs.Index(),
	)
}

func (s *Server) ListenAndServe() {
	if s.Server == nil {
		s.Server = &http.Server{}
	}
	s.Server.Handler = s.Handler()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

	go func() {
		err := s.Server.Serve(getListener())
		if err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	<-ch
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(7*time.Second))
	defer cancel()
	if err := s.Server.Shutdown(ctx); err == nil {
		log.Println(`shutdown`)
	} else {
		log.Printf("shutdown error: %v", err)
	}
}

func getListener() net.Listener {
	port := os.Getenv(`GOPORT`)
	if port == `` {
		port = `3000`
	}
	addr := `:` + port
	ln, err := net.Listen(`tcp`, addr)
	if err != nil {
		log.Panic(err)
	}
	log.Println(color.GreenString(`started.(` + addr + `)`))
	return tcpKeepAliveListener{ln.(*net.TCPListener)}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
