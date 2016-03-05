package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"os"
	//	"webdb"
)

const (
	storagePath = "/tmp/test"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	t, _ := template.ParseFiles("goupload.gtpl")
	t.Execute(w, t)
}

func LoginShowHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("login.gtpl")
	t.Execute(w, nil)
}

/*
func LoginAccessHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form["username"]
	password := r.Form["password"]

	if len(username[0]) == 0 || len(password[0]) == 0 {
		fmt.Fprintf(w, "invalid username or password")
		return
	}

	account, err := webdb.QueryUser(username[0])
	if err != nil {
		fmt.Fprintf(w, "xxxx")
		return
	}

	fmt.Println(username[0], password[0])
	fmt.Println(account.User, account.Password)
	fmt.Println("-------------------")
	//密码匹配
	if password[0] == account.Password {
		fmt.Println("密码匹配")
		//fmt.Fprintf(w, "你好,%s", account.User)
		http.Redirect(w, r, "/proc/", http.StatusFound)

	} else {
		fmt.Println("密码不匹配")
		fmt.Fprintf(w, "密码错误,%s", account.User)
	}

}
*/
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
}

func DirShowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("show")
	fmt.Println(r.URL.Path)
	w.Header().Set("Cache-Control", "no-store, no-cache")
	http.ServeFile(w, r, "/proc")
	//http.ServeFile(w, r, "./welcome.html")
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
func TestHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("不见鸟~~"))
	//http.ServeFile(w, r, "./welcome.html")
	vars := mux.Vars(r)
	content := vars["filename"]
	stringa := "welcomxxxxx " + content
	fmt.Fprint(w, stringa)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("不见鸟~~"))
	/*
		if r.URL.Path == "/" {
			http.Redirect(w, r, "./login/index.html", http.StatusFound)
		}*/ //访问/路径的时候,自动跳转到登录页面
	http.ServeFile(w, r, "./error404.html")
}

func main() {
	//启动数据库
	/*
		opts := webdb.DbOpts{User: "root", Password: "123456", Ip: "192.168.2.119", Port: "3306"}
		err := webdb.Open(opts)
		if err != nil {
			panic("can not connect to database")
		}
		defer webdb.Close()
	*/

	//新增一个web url路由
	r := mux.NewRouter()
	//无效url的控制器函数
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", DirShowHandler).Name("article")
	url, err := r.Get("article").URL("category", "technology", "id", "43")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
	/*
		r.HandleFunc("/", HomeHandler).Methods("GET")
		r.HandleFunc("/", UploadHandler).Methods("POST")
		r.HandleFunc("/login", LoginShowHandler).Methods("GET")
		r.HandleFunc("/login", LoginAccessHandler).Methods("POST")
		r.HandleFunc("/proc/", DirShowHandler).Methods("GET")
		r.HandleFunc("/proc/{filename}", TestHandler).Methods("GET")
		r.HandleFunc("/proc/{filename}/", TestHandler).Methods("GET")
	*/
	//这里将/后面的内容当做filename参数传递给DownloadHandler
	//多个参数怎么传递?
	//	r.HandleFunc("/{filename}", DownloadHandler).Methods("GET")
	//使用自己写的路由
	err = http.ListenAndServe(":9090", r)
	if err != nil {
		fmt.Println(err)
	}
}
