package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "exercise")
}

func DirShowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("show")
	fmt.Println(r.URL.Path)
	//浏览器可能缓存了/proc获取的数据,减少对服务器的访问
	//导致页面显示的数据和实际的数据不同步，因此添加响应头部
	//告知浏览器不要缓存获取的数据
	w.Header().Set("Cache-Control", "no-store, no-cache")
	http.ServeFile(w, r, "/proc")
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	//提取url传递过来的实际参数
	vars := mux.Vars(r)
	content := vars["filename"]
	stringa := "welcome " + content
	fmt.Fprint(w, stringa)
}

//无效url的控制
func notfound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./error404.html")
}

func main() {
	//新增一个web url路由
	r := mux.NewRouter()
	//无效url的控制器函数
	r.NotFoundHandler = http.HandlerFunc(notfound)
	//指定url="/"时,HTTP请求为"GET"时的动作
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/proc/", DirShowHandler).Methods("GET")
	//注意/proc同时存在目录和文件的情况，两者的url不同
	r.HandleFunc("/proc/{filename}", WelcomeHandler).Methods("GET")
	r.HandleFunc("/proc/{filename}/", WelcomeHandler).Methods("GET")

	err := http.ListenAndServe(":9090", r)
	if err != nil {
		fmt.Println(err)
	}
}
