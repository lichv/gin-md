package main

import (
	"flag"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gin-gonic/gin"
	"github.com/lichv/go"
	"github.com/thinkerou/favicon"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)
type ReadFileForm struct {
	Filepath  string `json:"filepath" form:"filepath"`
}
type Index struct {
	Id   string
	From string
	Body string
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

	iconPath := path.Join("./","favicon.ico")
	if lichv.IsExist(iconPath) {
		fmt.Println("icon 存在")
		engine.Use(favicon.New(iconPath))
	}else{
		fmt.Println("icon 不存在")
	}
	engine.Static("/static", path.Join(PublicPath,"/static"))
	engine.LoadHTMLFiles(path.Join(PublicPath,"/index.html"))
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.GET("/index.html", func(c *gin.Context) {
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

	engine.GET("/api/buildIndex", func(c *gin.Context) {
		searchFile := "search"
		flag := BuildIndexFromLocal(searchFile,"./docs/")
		c.JSON(200,gin.H{
			"state":2000,
			"message":"success",
			"data":flag,
		})
	})

	engine.Any("/api/search", func(c *gin.Context) {
		var page,size = 1,50
		var query = ""
		searchFile := "search"
		input, _ := GetMapFromContext(c)
		inputQuery,ok := input["query"]
		if ok {
			query = lichv.Strval(inputQuery)
		}
		inputPage,ok := input["page"]
		if ok {
			page = lichv.IntVal(inputPage)
			if page == 0 {
				page = 1
			}
		}
		inputSize,ok := input["size"]
		if ok {
			size = lichv.IntVal(inputSize)
			if size == 0 {
				size = 50
			}
		}

		result:= SearchFromIndex(searchFile, query, page, size)
		c.JSON(200,gin.H{
			"state":2000,
			"message":"success",
			"data":result,
		})

	})

	outportStr := strconv.Itoa(outport)
	fmt.Println("从浏览器打开：http://localhost:"+outportStr)
	engine.Run(":"+outportStr)

}
func SearchFromIndex(searchFile string,query string,page,size int) *bleve.SearchResult {
	var from = (page - 1)* size
	index, _ := bleve.Open(searchFile)
	defer index.Close()
	stringQuery := bleve.NewQueryStringQuery(query)
	request := bleve.NewSearchRequest(stringQuery)
	request.From = from
	request.Size = size
	search, err := index.Search(request)
	if err != nil {
		fmt.Println(err.Error())
	}

	return search
}

func BuildIndexFromLocal(searchFile string,dirname string) bool {
	filepath.Walk(dirname, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path,".md"){
			text, err := lichv.Read(path)
			if err != nil{
				return err
			}
			filename := strings.Replace(path[5:],"\\","/",-1)
			BuildIndex(searchFile,[]*Index{&Index{filename,filename,*text}})
		}
		return nil
	})

	return true
}

func BuildIndex(searchFile string,indexes []*Index)  {
	var index bleve.Index
	_,err :=os.Stat(searchFile)
	if err != nil {
		if os.IsNotExist(err) {
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New(searchFile, mapping)
			if err != nil {
				fmt.Println(err.Error())
			}

		}else{
			return
		}
	}else{
		index, _ = bleve.Open(searchFile)
	}

	for _,v := range indexes{
		ids := bleve.NewDocIDQuery([]string{v.Id})
		request := bleve.NewSearchRequest(ids)
		search, err := index.Search(request)
		if err != nil {
			fmt.Println(err.Error())
		}
		if search.Total == 0 {
			index.Index(v.Id,v)
		}
	}

	index.Close()
}

func GetMapFromContext(context *gin.Context) (map[string]interface{},error) {
	result := make(map[string]interface{})
	post := make(map[string]interface{})
	query := context.Request.URL.Query()
	contentType := strings.ToLower(context.Request.Header.Get("content-type"))
	if strings.Contains(contentType,"multipart/form-data"){
		err := context.Request.ParseMultipartForm(128)
		if err != nil {
			return map[string]interface{}{},err
		}
		form := context.Request.Form
		for k,v :=range form{
			if len(v) == 1 {
				post[k] = v[0]
			}else{
				post[k] = strings.Join(v,";")
			}
		}
	}else if  strings.Contains(contentType,"x-www-form-urlencoded") {
		err := context.Request.ParseForm()
		if err != nil {
			return map[string]interface{}{},err
		}
		form := context.Request.Form
		for k,v :=range form{
			if len(v) == 1 {
				post[k] = v[0]
			}else{
				post[k] = strings.Join(v,";")
			}
		}
	}else{
		_ = context.ShouldBind(&post)
	}
	for v,k := range query{
		result[v] = k
	}
	for v,k := range post{
		result[v] = k
	}
	return result,nil
}