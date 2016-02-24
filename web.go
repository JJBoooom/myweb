package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	storagePath = "/tmp/test"
)

/*简易路由器*/
type MyMux struct {
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	/*
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
	*/
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	//返回状态码
	//w.WriteHeader(http.StatusOK)
	fmt.Println(r.Proto)
	fmt.Println("-------------------------")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("HTTP status code returned"))
	//	fmt.Fprintf(w, "Hello world")

}

func login(w http.ResponseWriter, r *http.Request) {
	/*
		if r.URL.Path != "/login" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
	*/

	fmt.Println("method", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		//添加r.ParseForm()后服务端才能对表单数据进行操作
		r.ParseForm()
		//打印请求的头部
		fmt.Println(r.Header)

		//验证用户输入
		//r.Form包含所有的请求参数，里面可能存取多个值
		if len(r.Form["username"][0]) == 0 || len(r.Form["password"][0]) == 0 {
			fmt.Fprintf(w, "invalid username or password")
			return
		}

		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}

func mkdir_r(path string) (err error) {
	if len(path) == 0 {
		err := errors.New("path is empty")
		return err
	}

	if !strings.HasPrefix(path, "/") {
		err := errors.New("Not absolute path")
		return err
	}

	if err := os.MkdirAll(path, 0755); os.IsExist(err) {
		fmt.Printf("%s has existed!\n", path)
		return nil
	} else {
		return err
	}

}

func upload(w http.ResponseWriter, r *http.Request) {
	/*
		if r.URL.Path != "/upload" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
	*/

	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		crutime := time.Now().Unix()
		fmt.Print(crutime)
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)

	} else {
		//r.ParseMultipartForm解析request body作为multipart/form-data,
		//设置内存容量为32MB,request body被解析时将保存32MB的数据在内存中，其他部分保存在硬盘的临时文件中
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		err = mkdir_r(storagePath)
		if err != nil {
			fmt.Printf("Can't create dir %s: %s", storagePath, err)
			return
		}

		path := storagePath + "/" + handler.Filename

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

//自己写的路由
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/":
		sayhelloName(w, r)
	case "/login":
		login(w, r)
	case "/upload":
		upload(w, r)

	default:
		http.NotFound(w, r)
	}
	return

}
func main() {
	//	http.HandleFunc("/", sayhelloName)
	//	http.HandleFunc("/login", login)
	//	http.HandleFunc("/upload", upload)

	//	err := http.ListenAndServe(":9090", nil)
	//使用自己写的简单路由
	mux := &MyMux{}
	err := http.ListenAndServe(":9090", mux)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
