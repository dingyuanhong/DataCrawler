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
	// "regexp"
	"strconv"
	// "strings"

	// 其他包
	"fmt"
	// "math"
	"time"

    //
    "os"
    "io/ioutil"
	// "math/rand"
)


func init(){
	ShangHaiBond.Register();
}

//缓存数据
func cache(path string,data string){
    var d []byte = []byte(data);
    ioutil.WriteFile(path,d,0644);
}

//前一天
func LastDay(t time.Time) time.Time{
	return t.AddDate(0,0,-1);
}

//当前开盘日
func OpeningQuotationDay(t time.Time) time.Time {
	//周六周日不开盘
	if(t.Weekday() == time.Saturday){
		t = t.AddDate(0,0,-1);
	}else if(t.Weekday() == time.Sunday){
		t = t.AddDate(0,0,-2);
	}
	return t;
}

func OpeningQuotationDayNow(t time.Time) time.Time {
	//9:30 ~ 15:00 开盘收盘
	if(t.Hour() <= 9 && t.Minute() < 30){
		t = t.AddDate(0,0,-1);
	}
	if(t.Hour() < 15){
		t = t.AddDate(0,0,-1);
	}
	return OpeningQuotationDay(t);
}

//日期字符串
func GetDayString(t time.Time) string{
	return fmt.Sprintf("%04d-%02d-%02d",t.Year(),int(t.Month()),t.Day());
}

var DATA string= "./data/上海证券/";
var CREATE_DATE=map[string]int{"year":1990,"month":12,"day":26};
var LIMIT_DATE=map[string]int{"year":0,"month":0,"day":0};

//创建交易所日期
func GreaterCreateDay(t time.Time)bool{
	if(t.Year() < CREATE_DATE["year"]){
		return false;
	}else if(t.Year() == CREATE_DATE["year"]){
		if(int(t.Month()) < CREATE_DATE["month"]){
			return false;
		}
	}
	return true;
}

func LimitDay(t time.Time) bool{
	limit := LIMIT_DATE;
	if(t.Year() < limit["year"]){
		return false;
	}else if(t.Year() == limit["year"]){
		if(int(t.Month()) < limit["month"]){
			return false;
		}
	}
	return true;
}

//日期是否有效
func ValidDay(t time.Time)bool{
	if(LimitDay(t) == false){
		return false;
	}
	return GreaterCreateDay(t);
}

var rss_ShangHaiBond = map[string]string{
	"股票-总貌-统计":"http://query.sse.com.cn/security/stock/queryDataStatistics.do?jsonCallBack=jsonpCallback88632",
};

var DateSectionRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);
		tm := aid["time"].(string);

		v += "&startDate=";
		v += tm;
		v += "&endDate=";
		v += tm;

		v += "&_=";
		v += strconv.FormatInt(time.Now().Unix(),10);

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
		cache(dir + "/" + tm + ".html",ctx.GetText());

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
				ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v}, k)
			}
        },
        Trunk: map[string]*Rule{
			"股票-总貌-统计":DateSectionRule,
		},
	},
};
