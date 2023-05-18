package gocache

// 提供被其他节点访问的能力(基于http)

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_gocache/"

type HTTPPool struct {
	self     string // 记录自己的地址
	basePath string // 节点之间通讯地址的前缀
}

// 实例化 HTTPPool
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 日志输出 server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 处理所有的 HTTP 请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) { // 请求 URL 没有包含通用前缀
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	// parts 包含 group name 和 key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	// 通过 groupName 获取 group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	// 通过 key 获取该 group 的 value
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 写回数据
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.b)
}
