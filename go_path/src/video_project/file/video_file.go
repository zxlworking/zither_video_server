package video_file

import (
	"database/sql"
	"fmt"
	data "video_project/data"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		//panic(err)
		fmt.Println("register::checkErr::" + err.Error())
	}
}

var (
	dbhostip   = "localhost"
	dbusername = "root"
	dbpassword = "root"
	dbname     = "zither_video"
)

func QueryAllVideoFile() (int, []data.VideoFileInfo) {
	fmt.Println("video_file::QueryAllVideoFile")

	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")

	checkErr(err)
	stmt, err := db.Prepare("SELECT * FROM video")
	checkErr(err)

	rows, err := stmt.Query()

	checkErr(err)

	var videoFileInfoList []data.VideoFileInfo

	for rows.Next() {
		var video_id string
		var video_name string
		var video_path string
		var student_video_name string
		var student_video_path string
		var video_desc string
		var img_name string
		var user_id string
		var convert_video string
		err := rows.Scan(&video_id, &video_name, &video_path, &student_video_name, &student_video_path, &video_desc, &img_name, &user_id, &convert_video)
		checkErr(err)
		fmt.Println("user::QueryAllVideoFile::", video_id, video_name, video_path, student_video_name, student_video_path, video_desc, img_name, user_id, convert_video)
		videoFileInfo := data.VideoFileInfo{video_id, video_name, video_path, student_video_name, student_video_path, video_desc, img_name, user_id, convert_video}
		videoFileInfoList = append(videoFileInfoList, videoFileInfo)
	}
	defer rows.Close()
	defer stmt.Close()
	defer db.Close()
	return 0, videoFileInfoList
}

func QueryVideoFileByVideoId(videoId string) (int, data.VideoFileInfo) {
	fmt.Println("video_file::queryVideoFileByVideoId::videoId = " + videoId)

	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")

	checkErr(err)
	stmt, err := db.Prepare("SELECT * FROM video where video_id =?")
	checkErr(err)

	rows, err := stmt.Query(videoId)

	checkErr(err)
	for rows.Next() {
		var video_id string
		var video_name string
		var video_path string
		var student_video_name string
		var student_video_path string
		var video_desc string
		var img_name string
		var user_id string
		var convert_video string
		err := rows.Scan(&video_id, &video_name, &video_path, &student_video_name, &student_video_path, &video_desc, &img_name, &user_id, &convert_video)
		checkErr(err)
		fmt.Println("user::queryVideoFileByVideoId::", video_id, video_name, video_path, student_video_name, student_video_path, video_desc, img_name, user_id, convert_video)
		return 0, data.VideoFileInfo{video_id, video_name, video_path, student_video_name, student_video_path, video_desc, img_name, user_id, convert_video}
	}
	defer rows.Close()
	defer stmt.Close()
	defer db.Close()
	return -1, data.VideoFileInfo{}
}

func AddVideoFile(videoFileInfo data.VideoFileInfo) int {
	fmt.Println("video_file::AddVideoFile")

	db, openErr := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")
	checkErr(openErr)

	stmt, prepareErr := db.Prepare("insert into video (video_name, video_path,student_video_name,student_video_path, video_desc,img_name, user_id) values (?,?,?,?,?,?,?)")
	checkErr(prepareErr)

	fmt.Println("video_file::AddVideoFile::", videoFileInfo.VideoName, videoFileInfo.VideoPath, videoFileInfo.StudentVideoName, videoFileInfo.StudentVideoPath, videoFileInfo.VideoDesc, videoFileInfo.ImgName, videoFileInfo.UserId)
	_, execErr := stmt.Exec(videoFileInfo.VideoName, videoFileInfo.VideoPath, videoFileInfo.StudentVideoName, videoFileInfo.StudentVideoPath, videoFileInfo.VideoDesc, videoFileInfo.ImgName, videoFileInfo.UserId)
	checkErr(execErr)

	defer stmt.Close()
	defer db.Close()

	return 0
}

func UpdateVideoFileConvert(videoPath string, studentVideoPath string) int {
	fmt.Println("video_file::UpdateVideoFileConvert::videoPath = " + videoPath + "::studentVideoPath = " + studentVideoPath)

	db, openErr := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")
	checkErr(openErr)

	stmt, prepareErr := db.Prepare("update video set convert_video = 1 where video_path = ? and student_video_path = ?")
	checkErr(prepareErr)

	_, execErr := stmt.Exec(videoPath, studentVideoPath)
	checkErr(execErr)

	defer stmt.Close()
	defer db.Close()

	return 0
}

func DeleteVideoFileByVideoIds(videoIds []string) int {
	fmt.Print("video_file::DeleteVideoFileByVideoIds::videoIds = ")
	fmt.Println(videoIds)

	db, openErr := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")
	checkErr(openErr)

	for _, videoId := range videoIds {
		stmt, prepareErr := db.Prepare("delete from video where video_id = ?")
		checkErr(prepareErr)
		defer stmt.Close()

		_, execErr := stmt.Exec(videoId)
		checkErr(execErr)
	}

	defer db.Close()

	return 0
}
