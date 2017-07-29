package Controller

import (
	"MyProject/Modules"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//主页面
func GetIndex(context *gin.Context) {
	context.HTML(200, "index.html", gin.H{})
}

//上传文件
func UpLoadNewFile(context *gin.Context) {
	var save_file Modules.FileList
	form, err := context.MultipartForm()
	if err != nil {
		log.Println(err)
		return
	}
	files := form.File["files[]"]
	if len(files) == 0 {
		context.HTML(200, "success_load.html", gin.H{"Msg": "文件列表不能为空"})
		return
	}
	for _, f := range files {
		rst, err := f.Open()
		Modules.Db = Modules.DbSession.DB("go_web")
		file_data, err := ioutil.ReadAll(rst)
		if err != nil {
			log.Println(err)
			rst.Close()
			return
		}
		rst.Close()
		md5_ := md5.Sum(file_data)
		c := Modules.Db.C("file_list")
		md5_s := fmt.Sprintf("%x", md5_)
		err = c.Find(bson.M{"md5": string(md5_s)}).One(&save_file)
		if err == nil {
			log.Println("Already Have File MD5=" + string(md5_s) + " In List")
			context.HTML(200, "success_load.html", gin.H{"Msg": "文件\"" + f.Filename + "\"已经存在"})
			return
		}
		context.SaveUploadedFile(f, "./recv"+"/"+f.Filename)
		save_file.FileName = f.Filename
		save_file.FileSize = strconv.Itoa(len(file_data))
		save_file.MD5 = string(md5_s)
		save_file.FilePath = "./recv" + "/" + f.Filename
		err = c.Insert(&save_file)
		if err != nil {
			os.Remove("recv/" + f.Filename)
			log.Println(err)
			context.String(200, "写入数据库失败")
			return
		}
	}
	context.Redirect(http.StatusFound, "/success-upload")
}

//上传成功页面
func SuccessUpload(context *gin.Context) {
	context.HTML(200, "success_load.html", gin.H{"Msg": "上传成功"})
}

//下载列表页面
func DownloadList(context *gin.Context) {
	var err error
	var files []Modules.FileList
	Modules.Db = Modules.DbSession.DB("go_web")
	c := Modules.Db.C("file_list")
	err = c.Find(nil).All(&files)
	if err != nil {
		log.Print("Download fail:", err, "\n")
		context.HTML(200, "download_fail.html", gin.H{"Status": "加载资源失败"})
		return
	}
	context.HTML(200, "list_files.html", gin.H{"Files": files})
}

//下载
func Download(context *gin.Context) {
	var err error
	var file Modules.FileList
	md5_ := context.Param("md5")
	Modules.Db = Modules.DbSession.DB("go_web")
	c := Modules.Db.C("file_list")
	err = c.Find(bson.M{"md5": md5_}).One(&file)
	if err != nil {
		log.Println(err)
		context.HTML(200, "download_fail.html", gin.H{"Status": "下载失败:没有此文件资源"})
		return
	}
	context.File(file.FilePath)
}
