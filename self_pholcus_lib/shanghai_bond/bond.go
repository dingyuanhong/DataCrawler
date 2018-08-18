package self_pholcus_lib

import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	// "github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"                   //信息输出
	"github.com/tidwall/gjson"							//json解析

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/json"
	// "encoding/xml"

	// 字符串处理包
	// "regexp"
	"strconv"
	// "strings"

	// 其他包
	"fmt"
	// "math"
	"time"

    //
    "os"
    // "io/ioutil"
	// "math/rand"

	. "../util"
)


func init(){
	ShangHaiBond.Register();
}

var DATA string= "./data/上海证券/";
var CREATE_DATE=map[string]int{"year":1990,"month":12,"day":26};
var LIMIT_DATE=map[string]int{"year":0,"month":0,"day":0};
var PAGE_MAX = 25;

//日期是否有效
func ValidDay(t time.Time)bool{
	if(GreaterByDay(t,LIMIT_DATE) == false){
		return false;
	}
	return GreaterByDay(t,CREATE_DATE);
}

var history_ShangHaiBond = map[string]bool{
	"行情信息-行情报表":true,
	"股票-总貌-统计":false,
};

var rss_ShangHaiBond = map[string]string{
	// "股票-总貌-统计":"http://query.sse.com.cn/security/stock/queryDataStatistics.do?jsonCallBack=jsonp",
	// "行情信息-行情报表":"http://yunhq.sse.com.cn:32041/v1/sh1/list/exchange/equity?callback=jQuery&select=code%2Cname%2Copen%2Chigh%2Clow%2Clast%2Cprev_close%2Cchg_rate%2Cvolume%2Camount%2Ctradephase%2Cchange%2Camp_rate&order=",
	// "行情走势-分时线":"http://yunhq.sse.com.cn:32041/v1/sh1/line/%s?callback=jQuery",
	// "行情走势-日K线":"http://yunhq.sse.com.cn:32041/v1/sh1/dayk/%s?callback=jQuery",
	// //"股票列表-股票-上市A股":"http://query.sse.com.cn/security/stock/getStockListData2.do?&jsonCallBack=jsonpCallback&isPagination=true&stockCode=&csrcCode=&areaName=&stockType=1&pageHelp.cacheSize=1&pageHelp.pageSize=25",
	// //"股票列表-股票-上市B股":"http://query.sse.com.cn/security/stock/getStockListData2.do?&jsonCallBack=jsonpCallback&isPagination=true&stockCode=&csrcCode=&areaName=&stockType=2&pageHelp.cacheSize=1&pageHelp.pageSize=25",
	// //"股票列表-终止上市":"http://query.sse.com.cn/security/stock/getStockListData2.do?&jsonCallBack=jsonpCallback&isPagination=true&stockCode=&csrcCode=&areaName=&stockType=5&pageHelp.cacheSize=1&pageHelp.pageSize=25",
};

var DataSectionRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);
		tm := aid["time"].(string);

		v += "&startDate=";
		v += tm;
		v += "&endDate=";
		v += tm;

		v += "&_=";
		v += strconv.FormatInt(time.Now().UnixNano()/1000000,10);

		ctx.AddQueue(&request.Request{
			Url:v,
			Rule:k,
			Temp:aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		tm := ctx.GetTemp("time","").(string);
		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		Cache(dir + "/" + tm + ".html",ctx.GetText());

		loop := ctx.GetTemp("loop",false).(bool);
		v := ctx.GetTemp("v","").(string);
		//拉取历史记录
		if(loop == true){
			loc, _ := time.LoadLocation("Local")
			t ,_ := time.ParseInLocation("2006-01-02",tm,loc);
			t = OpeningQuotationDay(LastDay(t));
			if(ValidDay(t) == false){
				return;
			}
			tm = GetDayString(t);

			ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop}, k)
		}
	},
};

var DataPageRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);

		page := 1;
		if _, ok := aid["page"]; ok {
			page = aid["page"].(int);
		}

		v += "&begin=";
		v += strconv.Itoa(PAGE_MAX*(page-1));
		v += "&end=";
		v += strconv.Itoa(PAGE_MAX*page);

		//1533735162686
		v += "&_=";
		v += strconv.FormatInt(time.Now().UnixNano()/1000000,10);
		fmt.Println(v);

		ctx.AddQueue(&request.Request{
			Url:v,
			Rule:k,
			Temp:aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		tm := ctx.GetTemp("time","").(string);
		k := ctx.GetTemp("name","").(string);
		v := ctx.GetTemp("v","").(string);
		loop := ctx.GetTemp("loop",false).(bool);
		page := ctx.GetTemp("page",1).(int);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		text := ctx.GetText();
		text = TrimLeft(text,7);
		text = TrimRight(text,1);
		Cache(dir + "/" + strconv.Itoa(page) + ".html",text);

		pagemax := gjson.Get(text,"total").Int();
		if((page * PAGE_MAX) >= int(pagemax)){
			return;
		}
		//拉取历史记录
		if(loop == true){
			ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop,"page":page+1}, k)
		}
	},
}

var DataRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);
		// tm := aid["time"].(string);

		v += "&_=";
		v += strconv.FormatInt(time.Now().UnixNano()/1000000,10);

		ctx.AddQueue(&request.Request{
			Url:v,
			Rule:k,
			Temp:aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		tm := ctx.GetTemp("time","").(string);
		code := ctx.GetTemp("code","").(string);
		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);

		text := ctx.GetText();
		text = TrimLeft(text,7);
		text = TrimRight(text,1);
		Cache(dir + "/" + code + "_" + tm + ".html",text);
	},
}

var DataPage1Rule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);

		page := 1;
		if _, ok := aid["page"]; ok {
			page = aid["page"].(int);
		}

		v += "&pageHelp.pageNo=";
		v += strconv.Itoa(page);
		v += "&pageHelp.beginPage=";
		v += strconv.Itoa(page);
		if(page != 1){
			v += "&pageHelp.endPage=";
			v += strconv.Itoa(page*10+1);
		}

		v += "&_=";
		v += strconv.FormatInt(time.Now().UnixNano()/1000000,10);
		fmt.Println(v);

		ctx.AddQueue(&request.Request{
			Url:v,
			Rule:k,
			Temp:aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		tm := ctx.GetTemp("time","").(string);
		k := ctx.GetTemp("name","").(string);
		v := ctx.GetTemp("v","").(string);
		loop := ctx.GetTemp("loop",false).(bool);
		page := ctx.GetTemp("page",1).(int);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		text := ctx.GetText();
		text = TrimLeft(text,7);
		text = TrimRight(text,1);
		Cache(dir + "/" + strconv.Itoa(page) + ".html",text);

		pagemax := gjson.Get(text,"total").Int();
		if((page * PAGE_MAX) >= int(pagemax)){
			return;
		}
		fmt.Println(pagemax);
		//拉取历史记录
		if(loop == true){
			ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop,"page":page+1}, k)
		}
	},
}

var ShangHaiBond = &Spider{
	Name: "上海证券交易所",
    Description: "证券数据",
    EnableCookie:false,
    Namespace: nil,
    SubNamespace:func(self *Spider,dataCell map[string]interface{})string{
        return dataCell["Data"].(map[string]interface{})["分类"].(string)
    },
    RuleTree:&RuleTree{
        Root:func(ctx * Context){
			t := time.Now();
			t = OpeningQuotationDayNow(t);
			tm := GetDayString(t);

            for k,v:= range rss_ShangHaiBond {
				ctx.SetTimer(k, time.Minute*5, nil)
				loop := history_ShangHaiBond[k];
				ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop}, k)
			}
        },
        Trunk: map[string]*Rule{
			"股票-总貌-统计":DataSectionRule,
			"行情信息-行情报表":DataPageRule,
			"行情走势-分时线":&Rule{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					k := aid["name"].(string);
					v := aid["v"].(string);
					tm := aid["time"].(string);
					loop:= false;

					code := "600000";

					list := map[string]string{
						"数据":"time%2Cprice%2Cvolume",
					};

					for kl, vl := range list{
						url := fmt.Sprintf(v,code);
						url += "&selet=";
						url += vl;
						url += "&begin=";
						url += "0";
						url += "&end=";
						url += "-1";

						ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-" + kl,"v":url,"loop":loop,"code":code},k + "-" + kl);
					}
					return nil;
				},
			},
			"行情走势-分时线-数据":DataRule,
			"行情走势-日K线":&Rule{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					k := aid["name"].(string);
					v := aid["v"].(string);
					tm := aid["time"].(string);
					loop:= false;

					code := "600000";

					list := map[string]string{
						"数据":"date%2Copen%2Chigh%2Clow%2Cclose%2Cvolume%",
					};

					for kl, vl := range list{
						url := fmt.Sprintf(v,code);
						url += "&selet=";
						url += vl;
						url += "&begin=";
						url += "-300";
						url += "&end=";
						url += "-1";

						ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-" + kl,"v":url,"loop":loop,"code":code},k + "-" + kl);
					}
					return nil;
				},
			},
			"行情走势-日K线-数据":DataRule,
			"股票列表-股票-上市A股":DataPage1Rule,
			"股票列表-股票-上市B股":DataPage1Rule,
			"股票列表-终止上市":DataPage1Rule,
		},
	},
};
