package data

type ResponseBaseBean struct {
	Code int
	Desc string
}

type FileInfoResponseBean struct {
	BaseBean     ResponseBaseBean
	FileInfoList []FileInfo
}

type VideoFileInfoResponseBean struct {
	BaseBean          ResponseBaseBean
	VideoFileInfoList []VideoFileInfo
}

type LoginResponseBean struct {
	BaseBean ResponseBaseBean
	UserInfo UserInfo
}
