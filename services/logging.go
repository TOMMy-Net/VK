package services

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Logging struct {
	Writer  io.Writer
	Handler http.Handler
}

type LogParams struct {
	Request *http.Request
	URL     url.URL
	Time    time.Time
}

func LoggingHandler(out io.Writer, h http.Handler) http.Handler {
	return Logging{Writer: out, Handler: h}
}

func (l Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	l.Handler.ServeHTTP(w, r)
	params := LogParams{
		Request: r,
		Time:    t,
	}

	makeAndWriteLog(l.Writer, params)
}

func makeAndWriteLog(out io.Writer, params LogParams) {
	host, _, err := net.SplitHostPort(params.Request.RemoteAddr)
	if err != nil {
		host = params.Request.RemoteAddr
	}
	buf := make([]byte, 0, len(host))
	buf = append(buf, host...)
	buf = append(buf, " -- "...)
	buf = append(buf, params.Time.Format("15:04:05")...)
	buf = append(buf, " -- "...)
	buf = append(buf, params.Request.URL.Path...)
	buf = append(buf, " - "...)
	buf = append(buf, params.Request.Method...)
	buf = append(buf, " -- "...)
	buf = append(buf, params.Request.Proto...)
	buf = append(buf, '\n')
	out.Write(buf)
}
