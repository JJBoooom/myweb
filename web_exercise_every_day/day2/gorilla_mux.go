package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"os"
)

const (
	storagePath = "/tmp/test"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	t, _ := template.ParseFiles("goupload.gtpl")
	t.Execute(w, t)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ready to upload")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = os.MkdirAll(storagePath, 0755)
	if err != nil {
		if !os.IsExist(err) {
			fmt.Println("can't mkdir")
			return
		}
	}

	path := storagePath + "/" + handler.Filename

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//显示
	http.ServeFile(w, r, storagePath)

}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//获取url传递的参数
	filename := vars["filename"]
	filepath := storagePath + "/" + filename

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("%s doesn't exist", filename)
		return
	}
	//返回指定路径文件或者目录的内容
	http.ServeFile(w, r, filepath)

}

func notfound(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("不见鸟~~"))
	http.ServeFile(w, r, "./error404.html")
}

func main() {
	//新增一个web url路由
	r := mux.NewRouter()
	//无效url的控制器函数
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/", UploadHandler).Methods("POST")

	err := http.ListenAndServe(":9090", r)
	if err != nil {
		fmt.Println(err)
	}
}
