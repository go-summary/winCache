package http

import (
	"../localWinCache"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/winCache/"

// 定义一个http池
type HPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HPool {
	return &HPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只处理前缀为指定字母的那个
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := localWinCache.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(view.Clone())
}