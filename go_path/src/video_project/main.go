package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	fmt.Println("response", string(videoFileInfoResponseBeanResult))
}

func uploadVideoFile(w http.ResponseWriter, r *http.Request) {
	/*
		r.ParseForm()
		fmt.Println("uploadVideoFile start")
		//把上传的文件存储在内存和临时文件中
		r.ParseMultipartForm(32 << 20)
		//fmt.Println(r.Form)
		fmt.Println("path", r.URL.Path)
		fmt.Println("Scheme", r.URL.Scheme)

		videoName := r.URL.Query()["video_name"][0]
		videoDesc := r.URL.Query()["video_desc"][0]
		userId := r.URL.Query()["user_id"][0]

		fmt.Println(time.Now().UnixNano())
		videoRealName := strconv.FormatInt(time.Now().UnixNano(), 10)
		fmt.Println(videoRealName)

		//获取文件句柄，然后对文件进行存储等处理
		file, fileHandler, err := r.FormFile("file")
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
	*/

	getFormData(w, r)
}

func getFormData(w http.ResponseWriter, r *http.Request) {
	/*
		//获取 multi-part/form body中的form value
		for k, v := range form.Value {
			//fmt.Println("value,k,v = ", k, ",", v)
			fmt.Println("getFormData::form.Value::k = ", k, ",", len(v))

			if len(v) > 0 {
				dstFile, err := os.Create("./video_file/" + k)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				defer dstFile.Close()
				dstFile.WriteString(v[0])
			}
		}
		fmt.Println()
	*/

	r.ParseForm()
	videoName := r.URL.Query()["video_name"][0]
	videoDesc := r.URL.Query()["video_desc"][0]
	userId := r.URL.Query()["user_id"][0]
	fmt.Println("videoName", videoName)
	fmt.Println("videoDesc", videoDesc)
	fmt.Println("userId", userId)

	/**
	底层通过调用multipartReader.ReadForm来解析
	如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	*/
	r.ParseMultipartForm(32 << 20)
	//fmt.Println("r.Form:         ", r.Form)
	//fmt.Println("r.PostForm:     ", r.PostForm)
	//fmt.Println("r.MultiPartForm:", r.MultipartForm)

	form := r.MultipartForm

	fmt.Println(time.Now().UnixNano())
	videoRealName := strconv.FormatInt(time.Now().UnixNano(), 10)
	imgRealName := strconv.FormatInt(time.Now().UnixNano(), 10)

	for k, v := range form.File {
		//fmt.Println("value,k,v = ", k, ",", v)
		fmt.Println("getFormData::form.File::k = ", k, ",", len(v))

		for _, value := range v {
			fmt.Println("getFormData::form.File::FileName = ", value.Filename)
			f, _ := value.Open()
			buf, _ := ioutil.ReadAll(f)

			if strings.HasSuffix(value.Filename, ".mp4") {
				fmt.Println("getFormData::videoRealName = " + videoRealName)
				dstFile, err := os.Create("./video_file/" + videoRealName)

				if err != nil {
					fmt.Println(err.Error())
					return
				}
				defer dstFile.Close()
				dstFile.WriteString(string(buf))
			} else {
				fmt.Println("getFormData::imgRealName = " + imgRealName)
				dstFile, err := os.Create("./img_file/" + imgRealName)

				if err != nil {
					fmt.Println(err.Error())
					return
				}
				defer dstFile.Close()
				dstFile.WriteString(string(buf))
			}
		}
	}
	fmt.Println("getFormData::form.File::videoRealName = " + videoRealName)
	fmt.Println("getFormData::form.File::imgRealName = " + imgRealName)

	videoFileInfo := data.VideoFileInfo{"", videoName, videoRealName, videoDesc, imgRealName, userId, "0"}
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
	http.ListenAndServe(":8081", http.FileServer(http.Dir(baseProjectDir+"")))
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
