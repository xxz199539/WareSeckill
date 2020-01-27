package common

import "net/http"

// 声明一个新的数据类型
type FilterHandle func(w http.ResponseWriter, r *http.Request) error

// 拦截器结构体
type Filter struct {
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{filterMap:make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// 根据Uri获取相应的handle
func (f *Filter) GetFilterHandle(uri string)FilterHandle {
	return f.filterMap[uri]
}

// 声明新的函数类型
type WebHandle func(w http.ResponseWriter, r *http.Request)


// 执行拦截器，返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				// 执行拦截业务逻辑
				err := handle(w, r)
				if err != nil {
					_, _ = w.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		webHandle(w, r)
	}
}