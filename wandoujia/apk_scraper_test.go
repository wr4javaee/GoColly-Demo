package wandoujia

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"testing"
	"time"
)

func TestScrapeApk01(t *testing.T) {
	// 实例化
	c := colly.NewCollector()

	// 构造Header
	userAgentHeader := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36"

	// 设置请求头
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgentHeader)
	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*ip138.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wandoujia.*",
		Parallelism: 1,
		RandomDelay: 200 * time.Millisecond,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	q.AddURL("https://www.wandoujia.com/wdjweb/api/category/more?catId=5023&subCatId=0&page=30&ctoken=-5MjJED6pIug-JXmrJXFBWC3")

	// scraped handler
	c.OnScraped(func(r *colly.Response) {
		log.Println("response received end", r.StatusCode, string(r.Body[:]))
	})

	// error handler
	c.OnError(func(r *colly.Response, err error) {
		// 出现错误时处理
		log.Println(" ERROR Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	q.Run(c)
}