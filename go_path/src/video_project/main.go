package main

/*
import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	server = http.Server{
		Addr:           ":8080",
		Handler:        &HandlerStruct{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	handlerMap = make(map[string]HandlerFunc)
)

type HandlerStruct struct {
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (*HandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := handlerMap[r.URL.String()]; ok {
		handler(w, r)
	}
	urlStr := r.URL.String()
	index := strings.Index(urlStr, "/register?")
	fmt.Println(index)
	fmt.Println("ServeHTTP::url = " + urlStr)
	if index == 0 {
		r.ParseForm()
		fmt.Fprintln(w, "user_name = "+r.Form.Get("user_name"))
		fmt.Fprintln(w, "pass_word = "+r.Form.Get("pass_word"))
	}

}

func register(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Println(r.Header["Accept-Encoding"])
	fmt.Println(r.URL.Query())

}

func login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "login")
}

func main() {
	handlerMap["/register"] = register
	handlerMap["/login"] = login

	server.ListenAndServe()
}
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
	data "video_project/data"
	"video_project/file"
	"video_project/user"
)

var baseProjectDir = "C:\\zxl\\zither_video_server\\go_path\\src\\video_project\\"
var ffmpegBinDir = "C:\\Users\\zxl\\Downloads\\ffmpeg-20190426-4b7166c-win64-static\\bin\\ffmpeg.exe"

//返回一个Router实例
func NewRouter() *Router {
	return new(Router)
}

//路由结构体，包含一个记录方法、路径的map
type Router struct {
	Route map[string]map[string]http.HandlerFunc
}

//实现Handler接口，匹配方法以及路径
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := handlerMap[req.URL.EscapedPath()]; ok {
		handler(w, req)
	}
}

//根据方法、路径将方法注册到路由
func (r *Router) HandleFunc(method, path string, f http.HandlerFunc) {

}

func register(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query()["user_name"][0]
	passWord := r.URL.Query()["pass_word"][0]

	userinfo := data.UserInfo{"", userName, passWord, 0}
	registerResult := user.Register(userinfo)

	fmt.Println("main::register::", registerResult, userName, passWord)

	var response data.ResponseBaseBean
	if registerResult == 0 {
		response = data.ResponseBaseBean{registerResult, "success"}
	} else {
		response = data.ResponseBaseBean{registerResult, "fail"}
	}

	result, error := json.Marshal(response)
	if error != nil {

	}
	fmt.Fprint(w, string(result))
}

func login(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query()["user_name"][0]
	passWord := r.URL.Query()["pass_word"][0]

	loginUserInfo := data.UserInfo{"", userName, passWord, 0}
	loginResult, userInfo := user.Login(loginUserInfo)

	var loginResponseBean data.LoginResponseBean
	if loginResult == 0 {
		loginResponseBean.BaseBean.Code = loginResult
		loginResponseBean.BaseBean.Desc = "success"
		loginResponseBean.UserInfo = userInfo
	} else {
		loginResponseBean.BaseBean.Code = loginResult
		loginResponseBean.BaseBean.Desc = "fail"
	}

	result, error := json.Marshal(loginResponseBean)
	if error != nil {

	}
	fmt.Fprint(w, string(result))
}

func getVideoFileList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getVideoFileList")

	var videoFileInfoResponseBean data.VideoFileInfoResponseBean
	result, videoFileInfoList := video_file.QueryAllVideoFile()
	if result == 0 && videoFileInfoList != nil {
		videoFileInfoResponseBean.BaseBean.Code = 0
		videoFileInfoResponseBean.BaseBean.Desc = "success"
		videoFileInfoResponseBean.VideoFileInfoList = videoFileInfoList
	} else {
		videoFileInfoResponseBean.BaseBean.Code = -1
		videoFileInfoResponseBean.BaseBean.Desc = "fail"
	}

	videoFileInfoResponseBeanResult, _ := json.Marshal(videoFileInfoResponseBean)
	fmt.Fprint(w, string(videoFileInfoResponseBeanResult))
	fmt.Printf("response", string(videoFileInfoResponseBeanResult))
}

func uploadVideoFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("uploadVideoFile start")
	//把上传的文件存储在内存和临时文件中
	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("Scheme", r.URL.Scheme)

	videoName := r.URL.Query()["video_name"][0]
	videoDesc := r.URL.Query()["video_desc"][0]
	userId := r.URL.Query()["user_id"][0]

	fmt.Println(time.Now().UnixNano())
	videoRealName := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println(videoRealName)

	//获取文件句柄，然后对文件进行存储等处理
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("form file err: ", err)
		return
	}
	defer file.Close()
	//fmt.Fprintf(w, "%v", handler.Header)
	//创建上传的目的文件
	f, err := os.OpenFile("./video_file/"+videoRealName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err: ", err)
		return
	}
	defer f.Close()
	//拷贝文件
	io.Copy(f, file)

	videoFileInfo := data.VideoFileInfo{videoRealName, videoName, videoDesc, userId, "0"}
	result := video_file.AddVideoFile(videoFileInfo)

	fmt.Println("main::add_video_file::result = ", result)

	var response data.ResponseBaseBean
	if result == 0 {
		response = data.ResponseBaseBean{result, "success"}
	} else {
		response = data.ResponseBaseBean{result, "fail"}
	}

	responseResult, error := json.Marshal(response)
	if error != nil {

	}
	fmt.Println("main::add_video_file::success = ", string(responseResult))
	fmt.Fprint(w, string(responseResult))

	go mp4toflv(videoRealName)
}

func playVideoFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query()["file_name"][0]
	fmt.Println("playVideoFile::fileName = ", fileName)

	video, err := os.Open(string(baseProjectDir + "video_file\\" + fileName))
	if err != nil {

	}
	defer video.Close()

	http.ServeContent(w, r, "test.mp4", time.Now(), video)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

var handlerMap = make(map[string]HandlerFunc)

func startHttpServer() {
	r := NewRouter()

	handlerMap["/register"] = register
	handlerMap["/login"] = login
	handlerMap["/get_video_file_list"] = getVideoFileList
	handlerMap["/upload_video_file"] = uploadVideoFile
	handlerMap["/play_video_file"] = playVideoFile

	http.ListenAndServe(":8080", r)
}

func startFileServer() {
	http.ListenAndServe(":8081", http.FileServer(http.Dir(baseProjectDir+"video_file")))
}

func mp4toflv(videoId string) {
	cmdStr := ffmpegBinDir + " -i " + baseProjectDir + "video_file\\" + videoId + " -c:v libx264 -crf 24 " + baseProjectDir + "video_file\\" + videoId + ".flv"
	fmt.Println("cmdStr::", cmdStr)
	c := exec.Command("cmd", "/C", cmdStr)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	} else {
		video_file.UpdateVideoFileConvert(videoId)
		fmt.Println("cmdResult::success")
	}
}

func main() {

	go startHttpServer()

	go startFileServer()

	http.ListenAndServe(":8082", nil)
}