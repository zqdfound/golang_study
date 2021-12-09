// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"regexp"
// )

// var (
// 	reQQEmail = `(\d+)@qq.com`
// )

// func GetEmail() {
// 	resp, err := http.Get("https://tieba.baidu.com/p/6051076813?red_tag=1573533731")
// 	HandleError(err, "get url error")
// 	defer resp.Body.Close()

// 	pageBytes, err := ioutil.ReadAll(resp.Body)
// 	HandleError(err, "read resource error")
// 	pageStr := string(pageBytes)

// 	re := regexp.MustCompile(reQQEmail)
// 	results := re.FindAllStringSubmatch(pageStr, -1)
// 	for _, result := range results {
// 		fmt.Println(result[0])
// 		fmt.Println(result[1])
// 	}
// }

// func HandleError(err error, why string) {
// 	if err != nil {
// 		fmt.Println(why, err)
// 	}
// }

// func main() {
// 	GetEmail()
// }
