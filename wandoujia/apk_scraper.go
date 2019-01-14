package wandoujia

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/wr4javaee/web-scraper/engines"
	"log"
	"strconv"
	"time"
)

// 豌豆荚抓取信息Response对象
type ApkListResponseVO struct {
	Data ApkListResponseDataVO `json:"data"`
	State ApkListResponseStateVO `json:"state"`
}

// 豌豆荚抓取信息Response对象
type ApkListResponseDataVO struct {
	Content string `json:"content"`
	CurrPage int `json:"currPage"`
}

// 豌豆荚抓取信息Response对象
type ApkListResponseStateVO struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	Tips string `json:"tips"`
}

// 豌豆荚抓取信息持久化对象
type ApkPO struct {
	Id int `xorm:"pk autoincr"`
	AppId string `json:"appId"`
	AppVid string `json:"appVid"`
	AppName string `json:"appName"`
	AppPname string `json:"appPname"`
	AppVname string `json:"appVname"`
	AppVcode string `json:"appVcode"`
	AppCategoryId string `json:"appCategoryId"`
	AppCategoryName string `json:"appCategoryName"`
	AppRtype string `json:"appType"`
	AppInstallCount string `json:"appInstallCount"`
	AppTags string `json:"appTags"`
	AppTagLink string `json:"appTagLink"`
	AppIconLink string `json:"appIconLink"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// 金融理财一级分类
var GroupMap = make(map[string]string, 0)

// 初始化app的一级分组
func init() {
	// app一级分类
	//GroupMap["5017"] = "网上购物"
	GroupMap["5023"] = "金融理财"
	//GroupMap["5029"] = "影音播放"
	//GroupMap["5018"] = "系统工具"
	//GroupMap["5014"] = "通讯社交"
	//GroupMap["5024"] = "手机美化"
	//GroupMap["5019"] = "新闻阅读"
	//GroupMap["5016"] = "摄影图像"
	//GroupMap["5026"] = "考试学习"
	//GroupMap["5020"] = "生活休闲"
	//GroupMap["5021"] = "旅游出行"
	//GroupMap["5028"] = "健康运动"
	//GroupMap["5022"] = "办公商务"
	//GroupMap["5027"] = "育儿亲子"
}

// 抓取豌豆荚app list
func ScrapeApk() {
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

	baseUrl := "https://www.wandoujia.com/wdjweb/api/category/more"
	ctoken := "-5MjJED6pIug-JXmrJXFBWC3"
	maxPageSize := 1999
	startPageNo := 1
	// 根据app分组、分页信息, 初始化抓取队列
	// 一级分组, 如金融理财
	for groupKey := range GroupMap {
		catId := groupKey
		// 不处理子分组, 因为子分组信息无法从响应报文中准确获取, 还需另行爬取
		subCatId := "0"
		// 处理豌豆荚app列表分页, 每一组抓取分页由小到大
		for pageNo := startPageNo; pageNo <= maxPageSize; pageNo ++ {
			baseUrlStr := baseUrl + "?catId=" + catId + "&subCatId=" + subCatId + "&page=" + strconv.Itoa(pageNo) + "&ctoken=" + ctoken
			q.AddURL(baseUrlStr)
		}
	}

	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received start", r.StatusCode)
		if 200 != r.StatusCode {
			log.Println("response received error", string(r.Body[:]))
			return
		}

		// parse html
		responseVO := ApkListResponseVO{}
		json.Unmarshal(r.Body, &responseVO)
		responseBodyStr := responseVO.Data.Content

		// init goquery
		document, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(responseBodyStr)))
		if err != nil {
			log.Println("init goquery error", err)
		}

		document.Find(".card").Each(func(i int, selection *goquery.Selection) {
			// parse apk detail
			dtlSel := selection.Find(".detail-check-btn")
			// 7252955
			appId, _ := dtlSel.Attr("data-app-id")
			// 400730728
			appVid, _ := dtlSel.Attr("data-app-vid")
			// 信而富投资
			appName, _ := dtlSel.Attr("data-app-name")
			// com.crfchina.financial
			appPname, _ := dtlSel.Attr("data-app-pname")
			// 81
			appVcode, _ := dtlSel.Attr("data-app-vcode")
			// 3.2.9
			appVname, _ := dtlSel.Attr("data-app-vname")
			// https://android-artworks.25pp.com/fs08/2019/01/10/9/110_5aeaeae48c58d0d92d585db0e8c9bc24_con_130x130.png
			appIconLink, _ := dtlSel.Attr("data-app-icon")
			// 5023
			appCategoryId, _ := dtlSel.Attr("data-app-categoryid")
			// 0
			appRtype, _ := dtlSel.Attr("data-app-rtype")

			// parse app tag
			appCategoryName := selection.Find(".tag-link").Text()
			appTagLink, _ := selection.Find(".tag-link").Attr("href")
			// parse app install count
			appInstallCount := selection.Find(".install-count").Text()

			// 校验html
			if appPname == "" {
				log.Println("parse end for empty app_pname")
				return
			}

			// save app detail
			apkVO := new(ApkPO)
			apkPO := new(ApkPO)
			apkPO.AppId = appId
			apkPO.AppVid = appVid
			apkPO.AppName = appName
			apkPO.AppPname = appPname
			apkPO.AppVcode = appVcode
			apkPO.AppVname = appVname
			apkPO.AppIconLink = appIconLink
			apkPO.AppCategoryId = appCategoryId
			apkPO.AppCategoryName = appCategoryName
			apkPO.AppRtype = appRtype
			apkPO.AppTagLink = appTagLink
			apkPO.AppInstallCount = appInstallCount

			// 查询db是否有记录
			has, _ := engines.DefaultOrm.Table("tb_wdj_apk_info").
				Where("app_pname = ?", appPname).Get(apkVO)
			// 若存在记录, 则判定是否需更新
			if has {
				apkPO.UpdateTime = time.Now()
				_, err := engines.DefaultOrm.Table("tb_wdj_apk_info").ID(apkVO.Id).Update(apkPO)
				if err != nil {
					log.Println("orm update error", err)
				}
				log.Println("update a new log", apkPO)
			} else {
				// 不存在记录, 新增
				apkPO.CreateTime = time.Now()
				apkPO.UpdateTime = apkPO.CreateTime
				_, err := engines.DefaultOrm.Table("tb_wdj_apk_info").InsertOne(apkPO)
				if err != nil {
					log.Println("orm inert error", err)
				}
				log.Println("insert a new log", apkPO)
			}
		})
	})

	// scraped handler
	c.OnScraped(func(r *colly.Response) {
		log.Println("response received end")
	})

	// error handler
	c.OnError(func(r *colly.Response, err error) {
		// 出现错误时处理
		log.Println(" ERROR Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	q.Run(c)
}
