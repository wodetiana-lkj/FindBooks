package medabc

import (
	"io"
	"net/http"
	"scrap/utils"
)

const homeUrl = "http://www.medabc.com.cn"
const indexPath = "/quanben/1.html"
const searchPath = "/search.php?key="

func Search(name string) string {
	url := utils.Concat(homeUrl, searchPath, name)
	return get(url)
}

func Novel(path string) string {
	url := utils.Concat(homeUrl, path)
	return get(url)
}

func Chapter(novelPath string, chapterPath string) string {
	url := utils.Concat(homeUrl, novelPath, chapterPath)
	return get(url)
}

func setHeader(request *http.Request) {
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0")
	request.Header.Set("Cookie", "clickbids=null%2C11218")
}

func getBody(request *http.Request) string {
	client := http.Client{}
	setHeader(request)
	response, err := client.Do(request)
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		print("error read body:", err)
		return ""
	}
	return string(bodyBytes)
}

func get(path string) string {
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		print("error create request:", err)
		return ""
	}
	body := getBody(request)
	return body
}
