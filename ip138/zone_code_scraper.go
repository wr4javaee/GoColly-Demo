package ip138

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/wr4javaee/web-scraper/engines"
	"net/url"
	"strings"
	"time"

)

// 定义cityVO对象
type CityVO struct {
	Id int `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// 定义城市ZonePO对象
type CityZonePO struct {
	ProvinceName string
	CityName string
	ZoneName string
	PostCode string
	ZoneCode string
	CreateTime time.Time
}

var gbkDecoder = mahonia.NewDecoder("gbk")
var gbkEncoder = mahonia.NewEncoder("gbk")

func ScrapeZoneCode() {
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
		DomainGlob:  "*ip138.*",
		Parallelism: 2,
		RandomDelay: 200 * time.Millisecond,
	})

	// 增加队列
	cityArr := make([]CityVO, 0)
	engines.DefaultOrm.Table("city").Where("length(name) > ?", 3).Find(&cityArr)

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	for _, v := range cityArr {
		baseUrlStr := "http://www.ip138.com/post/search.asp?action=area2zone&area=" + url.QueryEscape(gbkEncoder.ConvertString(v.Name))
		fmt.Println(baseUrlStr)
		q.AddURL(baseUrlStr)
	}

	// Called right after OnResponse if the received content is HTML
	// ip138此处返回的数据编码格式为GB2312
	c.OnHTML("table:nth-child(3)", func(e *colly.HTMLElement) {
		// 使用goquery解析, case如下
		// <table width="349" border="1" style="border-collapse: collapse" bordercolor="#3366cc" cellpadding="4">
		//	<tbody><tr>
		//		<td class="tdc1" align="center" height="24" bgcolor="#6699cc">++* 查询结果 *++</td>
		//	</tr>
		//	<tr><td style="padding-left:5%" noswap="" class="tdc2">◎&nbsp;河北&nbsp;石家庄市&nbsp;桥西区&nbsp;邮编：050000&nbsp;区号：0311</td></tr><tr><td style="padding-left:5%" noswap="" class="tdc2">◎&nbsp;河北&nbsp;张家口市&nbsp;桥西区&nbsp;邮编：075000&nbsp;区号：0313</td></tr><tr><td style="padding-left:5%" noswap="" class="tdc2">◎&nbsp;河北&nbsp;邢台市&nbsp;桥西区&nbsp;邮编：054000&nbsp;区号：0319</td></tr><tr><td align="center" class="tdc2"><a href="http://alexa.ip138.com/post/" target="_blank">更详细的...</a></td></tr>
		//</tbody></table>
		e.DOM.Find(".tdc2").Each(func(i int, selection *goquery.Selection) {
			findText := selection.Text()
			// 按照&nbsp;分割
			findTextArr := strings.Split(findText, "\u00a0")
			resLen := len(findTextArr)
			cityZone := new(CityZonePO)
			cityZone.CreateTime = time.Now()
			if resLen > 1 && resLen == 4 {
				// <td style="padding-left:5%" noswap="" class="tdc2">
				// ◎&nbsp;北京&nbsp;邮编：100000&nbsp;区号：010</td>
				cityZone.ProvinceName = gbkDecoder.ConvertString(findTextArr[1])
				cityZone.CityName = gbkDecoder.ConvertString("-1")
				cityZone.ZoneName = gbkDecoder.ConvertString("-1")
				cityZone.PostCode = gbkDecoder.ConvertString(findTextArr[2][6:])
				cityZone.ZoneCode = gbkDecoder.ConvertString(findTextArr[3][6:])
				engines.DefaultOrm.Table("city_zone").InsertOne(cityZone)
				fmt.Println(cityZone)
			} else if resLen > 1 && resLen == 5 {
				// <td style="padding-left:5%" noswap="" class="tdc2">
				// ◎&nbsp;北京&nbsp;北京市&nbsp;邮编：100000&nbsp;区号：010</td>
				cityZone.ProvinceName = gbkDecoder.ConvertString(findTextArr[1])
				cityZone.CityName = gbkDecoder.ConvertString(findTextArr[2])
				cityZone.ZoneName = gbkDecoder.ConvertString("-1")
				cityZone.PostCode = gbkDecoder.ConvertString(findTextArr[3][6:])
				cityZone.ZoneCode = gbkDecoder.ConvertString(findTextArr[4][6:])
				engines.DefaultOrm.Table("city_zone").InsertOne(cityZone)
				fmt.Println(cityZone)
			} else if resLen > 1 && resLen == 6 {
				// <td style="padding-left:5%" noswap="" class="tdc2">
				// ◎&nbsp;河北&nbsp;石家庄市&nbsp;桥西区&nbsp;邮编：050000&nbsp;区号：0311</td>
				cityZone.ProvinceName = gbkDecoder.ConvertString(findTextArr[1])
				cityZone.CityName = gbkDecoder.ConvertString(findTextArr[2])
				cityZone.ZoneName = gbkDecoder.ConvertString(findTextArr[3])
				cityZone.PostCode = gbkDecoder.ConvertString(findTextArr[4][6:])
				cityZone.ZoneCode = gbkDecoder.ConvertString(findTextArr[5][6:])
				engines.DefaultOrm.Table("city_zone").InsertOne(cityZone)
				fmt.Println(cityZone)
			}
		})
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		// 出现错误时处理
		fmt.Println(" ERROR Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	q.Run(c)
	//c.Visit("http://www.ip138.com/post/search.asp?action=area2zone&area=%B1%B1%BE%A9")
}
