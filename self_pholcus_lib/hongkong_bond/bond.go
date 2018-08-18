package self_pholcus_lib

import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	// "github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// "github.com/tidwall/gjson"							//json解析

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/json"
	// "encoding/xml"

	// 字符串处理包
	"regexp"
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
	HongKongBond.Register();
}

var DATA string= "./data/香港证券/";
var CREATE_DATE=map[string]int{"year":1990,"month":12,"day":26};
var LIMIT_DATE=map[string]int{"year":0,"month":0,"day":0};
var PAGE_MAX = 25;

var history_HongKongBond = map[string]bool{

};

var rss_HongKongBondToken = map[string]string{
	"Token":"http://www.hkex.com.hk/Market-Data/Securities-Prices/Equities?sc_lang=zh-HK",
}

var rss_HongKongBond = map[string]string{
	"股本证券":"https://www1.hkex.com.hk/hkexwidget/data/getequityfilter?lang=chi&sort=5&order=0&all=1",
	"证券价格-股本证券":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1",
}

var rss_HongKongBondPrice = map[string]string{
	"证券价格-股本证券-1日-1分钟":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=0&int=0", //&ric=0700.HK
	"证券价格-股本证券-1日-5分钟":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=2&int=0",
	"证券价格-股本证券-1日-15分钟":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=3&int=0",
	"证券价格-股本证券-1日-小时":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=5&int=0",

	"证券价格-股本证券-1月-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=2",
	"证券价格-股本证券-3月-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=3",
	"证券价格-股本证券-3月-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=3",

	"证券价格-股本证券-6月-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=4",
	"证券价格-股本证券-6月-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=4",

	"证券价格-股本证券-本年-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=9",
	"证券价格-股本证券-本年-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=9",
	"证券价格-股本证券-本年-月线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=8&int=9",
	"证券价格-股本证券-本年-季线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=9&int=9",

	"证券价格-股本证券-1年-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=5",
	"证券价格-股本证券-1年-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=5",
	"证券价格-股本证券-1年-月线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=8&int=5",

	"证券价格-股本证券-2年-日线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=6&int=6",
	"证券价格-股本证券-2年-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=6",
	"证券价格-股本证券-2年-月线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=8&int=6",

	"证券价格-股本证券-5年-周线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=7&int=7",
	"证券价格-股本证券-5年-月线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=8&int=7",
	"证券价格-股本证券-5年-季线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=9&int=7",

	"证券价格-股本证券-10年-月线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=8&int=8",
	"证券价格-股本证券-10年-季线":"https://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=9&int=8",
};

var tokenRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);

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
		text := ctx.GetText();
		// Cache(dir + "/" + tm + ".html",text);

		reg := regexp.MustCompile(`return \".*?\"\;`);
		list := reg.FindAllString(text, -1);
		// fmt.Println(list);
		token := "";
		if(len(list) > 2) {
			text = list[1];
			text = TrimLeft(text,8);
			text = TrimRight(text,2);
			token = text;
		}
		// fmt.Println(token);

		ctx.Aid(map[string]interface{}{"time": tm,"token": token}, "loop")
	},
};

var bandRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);
		token := aid["token"].(string);

		v += "&token=";
		v += token;
		v += "&qid=";
		v += strconv.FormatInt(time.Now().UnixNano()/1000000,10);
		v += "&callback="
		v += "jQuery"
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
		// v := ctx.GetTemp("v","").(string);
		// loop := ctx.GetTemp("loop",false).(bool);
		// page := ctx.GetTemp("page",1).(int);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		text := ctx.GetText();
		text = TrimLeft(text,7);
		text = TrimRight(text,1);
		Cache(dir + "/" + tm + ".html",text);
	},
};

var HongKongBond = &Spider{
    Name: "香港证券交易所",
    Description: "证券数据",
    EnableCookie:false,
    Namespace: nil,
    SubNamespace:func(self *Spider,dataCell map[string]interface{})string{
        return dataCell["Data"].(map[string]interface{})["分类"].(string)
    },
    RuleTree:&RuleTree{
        Root:func(ctx * Context) {
			t := time.Now();
			t = OpeningQuotationDayNow(t);
			tm := GetDayString(t);

			for k,v:= range rss_HongKongBondToken {
				ctx.SetTimer(k, time.Minute*5, nil)
				loop := history_HongKongBond[k];
				ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop}, k)
			}
        },
        Trunk: map[string]*Rule{
			"Token":tokenRule,
			"loop":&Rule{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					t := time.Now();
					t = OpeningQuotationDayNow(t);
					tm := GetDayString(t);
					token := aid["token"].(string);

					fmt.Println("loop");
		            for k,v:= range rss_HongKongBond {
						ctx.SetTimer(k, time.Minute*5, nil)
						loop := history_HongKongBond[k];
						ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop,"token":token}, k)
					}
					return nil;
				},
			},
			"股本证券":bandRule,
			"证券价格-股本证券":&Rule{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					t := time.Now();
					t = OpeningQuotationDayNow(t);
					tm := GetDayString(t);
					token := aid["token"].(string);

					code := "0700";

					fmt.Println("loop");
		            for k,v:= range rss_HongKongBondPrice {
						ctx.SetTimer(k, time.Minute*5, nil)
						loop := history_HongKongBond[k];

						v += "&ric=";
						v += code;
						v += ".HK";

						ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v,"loop":loop,"token":token}, k)
					}
					return nil;
				},
			},
			"证券价格-股本证券-1日-1分钟":bandRule,
			"证券价格-股本证券-1日-5分钟":bandRule,
			"证券价格-股本证券-1日-15分钟":bandRule,
			"证券价格-股本证券-1日-小时":bandRule,
			"证券价格-股本证券-1月-日线":bandRule,
			"证券价格-股本证券-3月-日线":bandRule,
			"证券价格-股本证券-3月-周线":bandRule,
			"证券价格-股本证券-6月-日线":bandRule,
			"证券价格-股本证券-6月-周线":bandRule,
			"证券价格-股本证券-本年-日线":bandRule,
			"证券价格-股本证券-本年-周线":bandRule,
			"证券价格-股本证券-本年-月线":bandRule,
			"证券价格-股本证券-本年-季线":bandRule,
			"证券价格-股本证券-1年-日线":bandRule,
			"证券价格-股本证券-1年-周线":bandRule,
			"证券价格-股本证券-1年-月线":bandRule,
			"证券价格-股本证券-2年-日线":bandRule,
			"证券价格-股本证券-2年-周线":bandRule,
			"证券价格-股本证券-2年-月线":bandRule,
			"证券价格-股本证券-5年-周线":bandRule,
			"证券价格-股本证券-5年-月线":bandRule,
			"证券价格-股本证券-5年-季线":bandRule,
			"证券价格-股本证券-10年-月线":bandRule,
			"证券价格-股本证券-10年-季线":bandRule,
		},
	},
};
