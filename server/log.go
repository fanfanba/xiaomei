package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
)

var accessLog = openFile(`log/app.log`)
var errLog = openFile(`log/app.err`)

func writeLog(req *Request, res *Response, t time.Time, err interface{}) []byte {
	line := getLogLine(req, res, t, err)
	if err != nil {
		errLog.Write(line)
	} else {
		accessLog.Write(line)
	}
	return line
}

func getLogLine(req *Request, res *Response, t time.Time, err interface{}) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ' '
	writer.Write(getLogFields(req, res, t, err))
	writer.Flush()
	return buf.Bytes()
}

/*
  $time_iso8601 $host $request_method $request_uri $content_length $server_protocol
  $status $body_bytes_sent $request_time
  $session $remote_addr $http_referer $http_user_agent, $error, $stack
*/
func getLogFields(req *Request, res *Response, t time.Time, err interface{}) []string {
	slice := []string{t.Format(config.ISO8601), req.Host,
		req.Method, req.URL.RequestURI(), strconv.FormatInt(req.ContentLength, 10), req.Proto,
		strconv.FormatInt(res.Status(), 10), strconv.FormatInt(res.Size(), 10), time.Since(t).String(),
		fmt.Sprint(req.Session), req.ClientAddr(), req.Referer(), req.UserAgent(),
	}
	if err != nil {
		slice = append(slice, fmt.Sprint(err), string(utils.Stack(6)))
	}
	for i, v := range slice {
		v = strings.TrimSpace(v)
		if v == `` {
			v = `-`
		}
		slice[i] = v
	}
	return slice
}

func openFile(p string) *os.File {
	if f, err := os.OpenFile(
		path.Join(config.App.Root(), p), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666,
	); err != nil {
		panic(err)
	} else {
		return f
	}
}