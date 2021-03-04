package helper

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"errors"
)

type Destination string

const (
	ProfileImages       Destination = "files/profile_images"
	DailyWorkImages     Destination = "files/daily_work_images"
	ImportInfo          Destination = "files/import-info"
	LearningResultFiles Destination = "files/learning_result_files"
	ReportDailyWork     Destination = "files/report_daily_work"
	ReportTimestamp     Destination = "files/report_timestamp"
	NewsCoverImage      Destination = "files/news/cover_images"
	NewsDocs            Destination = "files/news/news_images"
	NewsImage           Destination = "files/news/news_files"
)

type ExtensionFile string

const (
	IsPng  ExtensionFile = "png"
	IsJpg  ExtensionFile = "jpg"
	IsJpeg ExtensionFile = "jpeg"
	IsPpt  ExtensionFile = "ppt"
	IsPptx ExtensionFile = "pptx"
	IsDoc  ExtensionFile = "doc"
	IsDocx ExtensionFile = "docx"
	IsXls  ExtensionFile = "xls"
	IsXlsx ExtensionFile = "xlsx"
	IsPdf  ExtensionFile = "pdf"
)

type FileObj struct {
	FileName  string `json:"file_name"`
	FileUrl   string `json:"file_url"`
	Extension string `json:"extension"`
	MIMEType  string `json:"mime_type"`
}

type UploadService struct {
}

// NewUploadService() is constructor function for UploadService
func NewUploadService() *UploadService {
	return &UploadService{}
}

//UploadFileService
func (svc UploadService) UploadFileService(file *multipart.FileHeader, destination Destination, fileObj *FileObj) error {

	src, _ := file.Open()
	path := "/static/" + destination + "/"
	fileName := strings.ReplaceAll(file.Filename, " ", "_")

	// Create folder is not exist
	_, isNotExist := os.Stat("static/" + string(destination))
	if os.IsNotExist(isNotExist) {
		os.Mkdir("static/"+string(destination), 0755)
	}

	dst, err := os.Create("." + string(path) + fileName)
	if err != nil {
		log.Printf("Error UploadFileService line 78 \n%s\n", err.Error())
		return err
	}

	// copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Printf("Error UploadFileService line 84 \n%s\n", err.Error())
		return err
	}

	savePath := os.Getenv("PATH_DOWNLOAD") + string(path) + fileName
	filePath := savePath
	fileObj.FileName = fileName
	fileObj.FileUrl = filePath
	return nil
}

//CheckFileType
func (svc UploadService) CheckFileType(file *multipart.FileHeader, fileObj *FileObj) error {

	mimeType := file.Header["Content-Type"][0]
	fileObj.MIMEType = mimeType
	switch mimeType {
	case "image/png":
		fileObj.Extension = "png"
	case "image/jpg":
		fileObj.Extension = "jpg"
	case "image/jpeg":
		fileObj.Extension = "jpeg"
	case "application/vnd.ms-powerpoint":
		fileObj.Extension = "ppt"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		fileObj.Extension = "pptx"
	case "application/msword":
		fileObj.Extension = "doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		fileObj.Extension = "docx"
	case "application/vnd.ms-excel":
		fileObj.Extension = "xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		fileObj.Extension = "xlsx"
	case "application/pdf":
		fileObj.Extension = "pdf"
	default:
		return FileIsNotSupport()
	}
	return nil
}

//UploadCondition
func (svc UploadService) UploadCondition(fileObj *FileObj, extension []ExtensionFile) bool {
	isOk := false
	for _, ex := range extension {
		if string(ex) == fileObj.Extension {
			isOk = true
		}
	}
	return isOk
}

//IUploadFileService is interface
type IUploadFileService interface {
	CheckFileType(file *multipart.FileHeader, fileObj *FileObj) error
	UploadFileService(file *multipart.FileHeader, destination Destination, fileObj *FileObj) error
	UploadCondition(fileObj *FileObj, extension []ExtensionFile) bool
}

//UploadFile
func UploadFile(svc IUploadFileService, file *multipart.FileHeader, destination Destination, extensions []ExtensionFile, fileObj *FileObj) error {
	errType := svc.CheckFileType(file, fileObj)
	if errType != nil {
		return errType
	}

	ok := svc.UploadCondition(fileObj, extensions)
	if !ok {
		return FileIsNotSupport()
	}

	err := svc.UploadFileService(file, destination, fileObj)
	if err != nil {
		return err
	}
	return nil
}

//FileIsNotSupport
func FileIsNotSupport() error {
	return errors.New("file is not support")
}
