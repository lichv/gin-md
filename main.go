package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lichv/go"
	"github.com/thinkerou/favicon"
	"net/http"
	"os"
	"path"
	"strconv"
)
type ReadFileForm struct {
	Filepath  string `json:"filepath" form:"filepath"`
}

func main() {
	var DocsPath string
	var PublicPath string
	var outport int
	flag.StringVar(&DocsPath,"d","./docs","文档存在目录")
	flag.StringVar(&PublicPath,"s","./public/","前端目录")
	flag.IntVar(&outport,"o",8044,"输出端口")

	if flag.Parsed(){
		flag.Parse()
	}

	if !lichv.IsExist(DocsPath) {
		fmt.Println("文档目录不存在")
		return
	}
	if !lichv.IsExist(PublicPath) {
		fmt.Println("静态文件目录不存在")
		return
	}

	engine := gin.Default()

	iconPath := path.Join(".","favicon.icon")
	if lichv.IsExist(iconPath) {
		engine.Use(favicon.New(iconPath))
	}
	engine.Static("/_assets", path.Join(PublicPath,"/static"))
	engine.LoadHTMLFiles(path.Join(PublicPath,"/index.html"))
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.Any("/api/markdown/files", func(context *gin.Context) {
		result,_ := lichv.GetFileTree(DocsPath)
		context.JSON(200,gin.H{
			"state":2000,
			"message":"success",
			"data":result,
		})
	})

	engine.Any("/api/markdown/read", func(context *gin.Context) {
		var readfile ReadFileForm
		context.Bind(&readfile)
		filepath := readfile.Filepath
		if filepath == ""{
			filepath = context.DefaultPostForm("path","")
		}
		if filepath == "" {
			filepath = context.DefaultQuery("path","")
		}
		filepath = path.Join(DocsPath,filepath)
		if lichv.IsDir(filepath) {
			filepath = path.Join(filepath,"index.md")
		}
		if lichv.IsExist(filepath){
			bs,err := os.ReadFile(filepath)
			if err != nil {
				context.JSON(404,gin.H{
					"state":4000,
					"message":"failed",
				})
				return
			}
			context.JSON(200,gin.H{
				"state":2000,
				"message":"success",
				"data":string(bs),
			})
		}else{
			context.JSON(404,gin.H{
				"state":4000,
				"message":"failed",
			})
		}
	})

	outportStr := strconv.Itoa(outport)
	fmt.Println("从浏览器打开：http://localhost:"+outportStr)
	engine.Run(":"+outportStr)

}
