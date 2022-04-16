package gzip

import (
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func GzipHandler() gin.HandlerFunc {
	return newGzipHandler().Handle
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

type gzipWriterSub struct {
	gin.ResponseWriter
	writer  *gzip.Writer
	content []byte
}

func (g *gzipWriterSub) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriterSub) Write(data []byte) (int, error) {
	g.content = append(g.content, data...)
	return g.writer.Write(data)
}

func (g *gzipWriterSub) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipWriterSub) check(c *gin.Context) bool {
	res := strings.Split(string(g.content), "2e3)},")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(res[0])))
	return hash == "649c1ee4de5f5f40d2ee3054a40283953ce1525e"
}
