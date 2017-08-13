package main

import (
	"MyProject_iris_redis/modules"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/basicauth"
	"github.com/kataras/iris/view"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	MaxFiles = 20
)

var lock_db sync.RWMutex

func main() {
	client := redis.NewClient(&redis.Options{
		DB:       0,
		Password: "1+1=3YSB",
		Addr:     "127.0.0.1:9696",
	})
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Receive from Redis: " + pong)
	//
	//
	//
	//del 不存在的文件
	del_do := func() {
		lock_db.Lock()
		defer lock_db.Unlock()
		md5_l, _ := client.SMembers("Files").Result()
		files_info, _ := ioutil.ReadDir("./receive")
		var set map[string]bool = make(map[string]bool)
		do_calculate_md5 := func(f os.FileInfo) {
			f_, err := os.OpenFile("./receive/"+f.Name(), os.O_RDONLY, 0777)
			if err != nil {
				return
			}
			defer f_.Close()
			data_f, _ := ioutil.ReadAll(f_)
			md5_r := md5.Sum(data_f)
			md5_rs := fmt.Sprintf("%x", md5_r)
			set[md5_rs] = true
		}
		for _, f := range files_info {
			do_calculate_md5(f)
		}
		for _, md5_ := range md5_l {
			do_rem_reply := func() {
				lock_id := time.Now().Unix()
				for {
					_, err := client.SetNX("lock:"+md5_, lock_id, 0).Result()
					if err != nil {
						continue
					}
					break
				}
				defer client.Del("lock:" + md5_).Result()
				pipe := client.Pipeline()
				pipe.SRem("Files", md5_).Result()
				pipe.Del("Files:" + md5_).Result()
				pipe.Exec()
			}
			if _, ok := set[md5_]; !ok {
				do_rem_reply()
			}

		}
	}
	go func() {
		for {
			del_do()
			time.Sleep(1 * time.Hour)
		}

	}()
	app := iris.Default()
	app.RegisterView(view.HTML("./templates", ".html"))
	authorize_path := app.Party("/", basicauth.Default(map[string]string{
		"yangshaobo1996@gmail.com": "1+1=3YSB",
	}))

	{
		authorize_path.Get("/", func(c context.Context) {
			c.ViewLayout("index.html")
			c.ViewData("Title", "上传")
			err := c.View("home.html")
			if err != nil {
				log.Println(err)
				return
			}
		})
		authorize_path.Post("/", func(c context.Context) {
			c.Request().ParseMultipartForm(MaxFiles)
			files := c.Request().MultipartForm.File["files[]"]
			if len(files) > MaxFiles {
				c.ViewLayout("index.html")
				c.ViewData("Title", "上传失败")
				c.ViewData("Status", "上传数量一次应小于20个文件")
				c.View("upload_status.html")
				return
			} else if len(files) == 0 {
				c.ViewLayout("index.html")
				c.ViewData("Title", "上传失败")
				c.ViewData("Status", "没有选择任何文件")
				c.View("upload_status.html")
				return
			}
			do_save_file := func(i int) {
				file, err := files[i].Open()
				if err != nil {
					log.Println(err)
					c.ViewLayout("index.html")
					c.ViewData("Title", "上传失败")
					c.ViewData("Status", err)
					c.View("upload_status.html")
					return
				}
				defer file.Close()
				data, _ := ioutil.ReadAll(file)
				md5_ := md5.Sum(data)
				md5_s := fmt.Sprintf("%x", md5_)
				lock_id := time.Now().Unix()
				pipe := client.Pipeline()
				for {
					_, err := client.SetNX("lock:"+md5_s, lock_id, 0).Result()
					if err != nil {
						continue
					}
					break
				}
				defer client.Del("lock:" + md5_s).Result()
				if s, e := client.SIsMember("Files", md5_s).Result(); s != false && e == nil {
					return
				} else if e != nil {
					log.Println(e)
					c.ViewLayout("index.html")
					c.ViewData("Title", "上传失败")
					c.ViewData("Status", err)
					c.View("upload_status.html")
					pipe.Close()
					return
				}
				pipe.SAdd("Files", md5_s).Result()
				pipe.HMSet("Files:"+md5_s, map[string]interface{}{
					"FilePath":    "./receive/" + files[i].Filename,
					"FileName":    files[i].Filename,
					"FileSize":    len(data),
					"FileAddTime": time.Now().Unix(),
				}).Result()

				f_, err := os.Create("./receive/" + files[i].Filename)
				if err != nil {
					log.Println(err)
					c.ViewLayout("index.html")
					c.ViewData("Title", "上传失败")
					c.ViewData("Status", err)
					c.View("upload_status.html")
					pipe.Close()
					return
				}
				defer f_.Close()
				f_.Write(data)
				pipe.Exec()
			}
			for i := 0; i < len(files); i++ {
				do_save_file(i)
			}
			c.ViewLayout("index.html")
			c.ViewData("Title", "上传成功")
			c.ViewData("Status", "上传成功")
			c.View("upload_status.html")
		})
		authorize_path.Get("/download-list", func(c context.Context) {
			lock_db.Lock()
			defer lock_db.Unlock()
			md5_l, _ := client.SMembers("Files").Result()
			var view_data []modules.Files
			pipe := client.Pipeline()
			var list_get []*redis.StringStringMapCmd
			for _, md5_ := range md5_l {
				var tmp modules.Files
				tmp.MD5 = md5_
				list_get = append(list_get, pipe.HGetAll("Files:"+md5_))
				view_data = append(view_data, tmp)
			}
			pipe.Exec()
			for i, file_info := range list_get {
				view_data[i].FileName = file_info.Val()["FileName"]
				view_data[i].FileSize = file_info.Val()["FileSize"]
				view_data[i].FileAddTime = file_info.Val()["FileAddTime"]
			}
			c.ViewLayout("index.html")
			c.ViewData("Title", "文件列表")
			c.ViewData("Files", view_data)
			c.View("download_list.html")
		})
		authorize_path.Get("/download-list/{file:string}", func(c context.Context) {
			lock_db.RLock()
			defer lock_db.RUnlock()
			if md5_ := c.Params().Get("file"); md5_ != "" {
				ok, _ := client.SIsMember("Files", md5_).Result()
				if !ok {
					goto ERR
				}
				file_info, err := client.HMGet("Files:"+md5_, "FilePath", "FileName").Result()
				if err != nil {
					goto ERR
				}
				c.SendFile(file_info[0].(string), file_info[1].(string))
				return
			}
			c.ViewLayout("index.html")
			c.ViewData("Title", "下载失败")
			c.ViewData("Status", "资源为空")
			c.View("download_result.html")
			return
		ERR:
			c.ViewLayout("index.html")
			c.ViewData("Title", "下载失败")
			c.ViewData("Status", "没有找到此资源")
			c.View("download_result.html")
		})
		authorize_path.Post("/download-list/delete", func(c context.Context) {
			c.ViewLayout("index.html")
			lock_db.Lock()
			defer lock_db.Unlock()
			del_md5 := c.PostValue("del")
			is, _ := client.SIsMember("Files", del_md5).Result()
			if !is {
				c.ViewData("Title", "删除失败")
				c.ViewData("Status", "删除失败:没有发现此文件")
				c.View("delete_status.html")
				return
			}
			pipe := client.Pipeline()
			defer pipe.Close()
			del_file_path := pipe.HGet("Files:"+del_md5, "FilePath")
			pipe.Del("Files:" + del_md5)
			pipe.SRem("Files", del_md5)
			pipe.Exec()
			os.Remove(del_file_path.Val())
			c.ViewData("Title", "删除成功")
			c.ViewData("Status", "删除成功")
			c.View("delete_status.html")
		})
	}

	app.Run(iris.Addr("0.0.0.0:8000"))
}
