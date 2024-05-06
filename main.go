package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"net/http"
	"regexp"
	"scrap/db"
	"scrap/utils"
	"strconv"
	"strings"
	"time"
)

const ParagraphSpace = "    "
const Separator = "\r\n"

func main() {
	keyword := flag.String("key", "", "关键字")
	novelNum := flag.Int("novel", 0, "输入文章编码")
	flag.Parse()
	if *novelNum != 0 {
		fmt.Println("开始获取文章", *novelNum)
		getNovel(*novelNum)
	}
	if !strings.EqualFold(*keyword, "") {
		fmt.Println("关键字查询", *keyword)
		getNovelList(*keyword)
	}
}

func getNovelList(keyword string) {

	novelCollector := colly.NewCollector()
	novelCollector.OnHTML("div.result li > a", func(e *colly.HTMLElement) {
		imgUrl := e.ChildAttr("img", "data-original")
		resp, _ := http.Get(imgUrl)
		bytes, _ := io.ReadAll(resp.Body)
		novel := &db.Novel{
			Name:   e.Attr("title"),
			ImgUrl: imgUrl,
			Img:    bytes,
			Path:   e.Attr("href"),
		}
		connect := db.Connect()
		connect.Save(novel)
		conn, err := connect.DB()
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = conn.Close()
	})
	novelCollector.Visit("http://www.medabc.com.cn/search.php?key=" + keyword)
}

type Chapter struct {
	Title      string
	Content    string
	RequestURL string
}

const chapterCtx = "chapter"
const begin = "begin"

func getNovel(novelNum int) {
	chapterCollector := colly.NewCollector()
	err := chapterCollector.Limit(&colly.LimitRule{
		DomainRegexp: `www.medabc.com.cn`,
		Parallelism:  1,
	})
	novelCollector := colly.NewCollector()
	err = novelCollector.Limit(&colly.LimitRule{
		DomainRegexp: `www.medabc.com.cn`,
		RandomDelay:  30 * time.Second,
		Parallelism:  1,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// 配置回调
	chapterCollector.OnHTML("li#chapter a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(e.Text, link)
		novelCollector.Visit(e.Request.AbsoluteURL(link))
	})

	novelCollector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put(chapterCtx, &Chapter{
			Title:      "",
			Content:    "",
			RequestURL: "",
		})
		r.Ctx.Put(begin, time.Now())
	})
	novelCollector.OnHTML("div#txt p,title", func(e *colly.HTMLElement) {
		if chapter, ok := e.Request.Ctx.GetAny(chapterCtx).(*Chapter); ok {
			if e.Name == "title" {
				chapter.Title = e.Text
			} else {
				chapter.Content += ParagraphSpace + e.Text + Separator
			}
		} else {
			fmt.Println("")
		}

	})
	novelCollector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL)
	})
	novelCollector.OnScraped(func(r *colly.Response) {
		begin := r.Ctx.GetAny(begin).(time.Time)
		fmt.Println(time.Now(), " => ", "delay: ", time.Now().Sub(begin).Seconds())
		connect := db.Connect()
		if chapter, ok := r.Ctx.GetAny(chapterCtx).(*Chapter); ok {
			pattern := `第(.*?)章`
			re := regexp.MustCompile(pattern)
			matches := re.FindAllStringSubmatch(chapter.Title, -1)
			var chapterNum int
			if len(matches) > 0 {
				chapterNumStr := matches[0][1]
				if utils.IsNumeric(chapterNumStr) {
					chapterNum, _ = strconv.Atoi(chapterNumStr)
				} else {
					chapterNum = utils.ChineseToArabic(chapterNumStr)
				}
			}
			dbChapter := &db.Chapter{
				Title:      chapter.Title,
				Content:    chapter.Content,
				RequestURL: r.Request.URL.String(),
				Number:     chapterNum,
				BookId:     novelNum,
			}
			connect.Save(dbChapter)
			conn, _ := connect.DB()
			_ = conn.Close()
		}
	})

	url := fmt.Sprintf("http://www.medabc.com.cn/novel/%d/", novelNum)
	// 开启线程
	err = chapterCollector.Visit(url)
	if err != nil {
		return
	}
	chapterCollector.Wait()
	novelCollector.Wait()
	fmt.Println("end")
}

func (c *Chapter) clear() {
	c.Title = ""
	c.Content = ""
	c.RequestURL = ""
}
