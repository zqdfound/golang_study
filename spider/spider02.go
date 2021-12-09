package main

//并发爬取图片资源
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//异常处理
func HandleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}

func DownloadFile(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	HandleError(err, "get url error")
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "read source error")
	filename = "E:/myCode/goStudy/imgs" + filename
	//write data
	err = ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		return false
	} else {
		return true
	}
}

var (
	//channel for images links store
	chanImagUrl chan string
	waitGroup   sync.WaitGroup
	//monitor
	chanTask chan string
	reImg    = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

//初始化管道
//获取图片路径协程：将图片路径加入到管道chanImagUrl
//监听任务协程:chanTask,检测到所有任务完成则关闭
//下载协程:从chanImagUrl读取url并下载图片到本地
func main() {
	//init channel
	chanImagUrl = make(chan string, 1000000)
	// waitGroup = make(chan string, 26)

	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go getImgUrl("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/" + strconv.Itoa(i) + ".html")
	}
	//任务统计协程 判断任务是否已经全部完成，完成则关闭channel
	waitGroup.Add(1)
	go CheckOk()
	//下载协程，从管道读取并下载
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go DownloadImg()
	}
	waitGroup.Wait()
}

//获取页面图片链接并存入管道,url是整个页面链接
func getImgUrl(url string) {
	urls := getImg(url)
	//遍历所有url并存入channel
	for _, url := range urls {
		chanImagUrl <- url
	}
	chanTask <- url
	waitGroup.Done()
}

//获取当前页面的图片链接
func getImg(url string) (urls []string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url = result[0]
		urls = append(urls, url)
	}
	return
}

//根据url获取页面内容
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError(err, "http get url")
	defer resp.Body.Close()
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil readAll")
	pageStr = string(pageBytes)
	return pageStr
}

//任务统计协程
func CheckOk() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s 完成了爬取任务\n", url)
		count++
		if count == 9 {
			close(chanImagUrl)
			break
		}
	}
	waitGroup.Done()
}

func DownloadImg() {
	for url := range chanImagUrl {
		filename := GetFilenameFromUrl(url)
		ok := DownloadFile(url, filename)
		if ok {
			fmt.Printf("%s 下载成功\n", filename)
		} else {
			fmt.Printf("%s 下载失败\n", filename)
		}
	}
	waitGroup.Done()
}

// 截取url名字
func GetFilenameFromUrl(url string) (filename string) {
	// 返回最后一个/的位置
	lastIndex := strings.LastIndex(url, "/")
	// 切出来
	filename = url[lastIndex+1:]
	// 时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	return
}
