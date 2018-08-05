package self_pholcus_lib


// 基础包
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
    "io/ioutil"
	"math/rand"
)

func init(){
    ShenZhengBond.Register();
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

//历史数据拉取
var history_ShenZhengBond = map[string]bool{
    "市场总貌-证券类别统计":true,
	"上市公司-上市公司列表":true,
	"上市公司-上市公司更名-全称":true,
	"上市公司-上市公司更名-简称":true,
	"上市公司-暂停上市公司":true,
	"上市公司-终止上市公司":true,
	"股票-A股列表":true,
	"股票-B股列表":true,
	"股票-A＋B股列表":true,
	"股票-基本指数-深圳市场":true,
	"股票-基本指数-深市主板":true,
	"股票-基本指数-中小企业板":true,
	"股票-基本指数-创业板":true,
	"股票-指标排名":true,

	"股票-行业统计-全部":true,
	"股票-行业统计-主板":true,
	"股票-行业统计-中小企业版":true,
	"股票-行业统计-创业板":true,

	"股票-行业市盈利":true,

	"行情-历史行情":true,
}

var rss_ShenZhengBond = map[string]string{
    // "市场总貌-证券类别统计":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803_sczm&TABKEY=tab1",
	// "上市公司-上市公司列表":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1110x&TABKEY=tab1",
	// "上市公司-上市公司更名-全称":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=SSGSGMXX&TABKEY=tab1",
	// "上市公司-上市公司更名-简称":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=SSGSGMXX&TABKEY=tab2",
	// "上市公司-暂停上市公司":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1793_ssgs&TABKEY=tab1",
	// "上市公司-终止上市公司":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1793_ssgs&TABKEY=tab2",
	// "股票-A股列表":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1110&TABKEY=tab1",
	// "股票-B股列表":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1110&TABKEY=tab2",
	// "股票-A＋B股列表":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1110&TABKEY=tab3",
	// "股票-基本指数-深圳市场":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803&TABKEY=tab1",
	// "股票-基本指数-深市主板":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803&TABKEY=tab2",
	// "股票-基本指数-中小企业板":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803&TABKEY=tab3",
	// "股票-基本指数-创业板":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803&TABKEY=tab4",
	// "股票-指标排名":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1805_zb1",//&selectModule=GE%2C30%2C00%2C20&TABKEY=ZGB&txtDqrq=2018-07-31
	// "股票-行业统计-全部":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1804_gptj&TABKEY=tab1",//&txtQueryDate=2018-07-31
	// "股票-行业统计-主板":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1804_gptj&TABKEY=tab2",
	// "股票-行业统计-中小企业版":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1804_gptj&TABKEY=tab3",
	// "股票-行业统计-创业板":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1804_gptj&TABKEY=tab4",
	// "股票-行业市盈利":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1804&TABKEY=tab1",

	// "行情-历史行情":"http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1815_stock&TABKEY=tab1&radioClass=00%2C20%2C30&txtSite=all",
	"行情-历史数据":"http://www.szse.cn/api/market/ssjjhq/getHistoryData?marketId=1",
}

var DATA string= "./data/";
var CREATE_DATE=map[string]int{"year":1990,"month":12,"day":1};
var LIMIT_DATE=map[string]int{"year":2018,"month":8,"day":4};

//指标排名
var select_index = map[string]string{
	"全部":"GE%2C30%2C00%2C20",
	"主板":"00%2C20",
	"中小企业板":"30",
	"创业板":"GE",
};
var module_index = map[string]string{
	"总股本":"ZGB",
	"总市值":"SJZZ",
	"流通股本":"LTGB",
	"流通市值":"LTSZ",
	"市盈利":"SYL",
};

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

//单页查询
var PageNORule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := rss_ShenZhengBond[k];

		tm := aid["time"].(string);

		page := 1;
		if _, ok := aid["pageno"]; ok {
			page = aid["pageno"].(int);
		}
		v += "&PAGENO=";
		v += strconv.Itoa(page);

		v += "&random=";
		v += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url:v,
			Rule:k,
			Temp:   map[string]interface{}{"time": tm,"name": k,"pageno":page},
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		page := ctx.GetTemp("pageno",1).(int);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + strconv.Itoa(page) + ".html",ctx.GetText());

		tm := ctx.GetTemp("time","").(string);
		//拉取历史记录
		if(history_ShenZhengBond[k] == true){
			pagecount := gjson.Get(ctx.GetText(),"0.metadata.pagecount").Int();
			fmt.Println(pagecount);
			if(int64(page + 1) > pagecount){
				return;
			}

			loc, _ := time.LoadLocation("Local")
			t ,_ := time.ParseInLocation("2006-01-02",tm,loc);
			tm = GetDayString(OpeningQuotationDay(LastDay(t)));
			ctx.Aid(map[string]interface{}{"time": tm,"name": k,"pageno":page + 1}, k)
		}
	},
};

//按天查询
var TimeRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		//txtQueryDate=2018-07-27&random=0.2845037210629777
		k := aid["name"].(string);
		v := rss_ShenZhengBond[k];

		tm := aid["time"].(string);
		v += "&txtQueryDate=";
		v += tm;

		v += "&random=";
		v += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url: v,
			Rule: k,
			Temp: aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		tm := ctx.GetTemp("time","").(string);
		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + tm + ".html",ctx.GetText());

		//拉取历史记录
		if(history_ShenZhengBond[k] == true){
			loc, _ := time.LoadLocation("Local")
			t ,_ := time.ParseInLocation("2006-01-02",tm,loc);
			t = OpeningQuotationDay(LastDay(t));
			if(ValidDay(t) == false){
				return;
			}
			tm = GetDayString(t);
			ctx.Aid(map[string]interface{}{"time": tm,"name": k}, k)
		}
	},
};

//指标排名
var Time2Rule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		//txtQueryDate=2018-07-27&random=0.2845037210629777
		k := aid["name"].(string);
		v := aid["v"].(string);

		tm := aid["time"].(string);
		url := v;
		url += "&txtDqrq=";
		url += tm;

		url += "&random=";
		url += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url: url,
			Rule: k,
			Temp: aid,
		});
		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		v := ctx.GetTemp("v","").(string);
		loop := ctx.GetTemp("loop",false).(bool);

		tm := ctx.GetTemp("time","").(string);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + tm + ".html",ctx.GetText());

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
}

//历史行情排名
var Time3Rule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := rss_ShenZhengBond[k];

		tm := aid["time"].(string);

		url := v;
		url += "&txtBeginDate=";
		url += tm;
		url += "&txtEndDate=";
		url += tm;

		var stockCode string = "";
		if _, ok := aid["stockCode"]; ok {
			stockCode = aid["stockCode"].(string);
		}
		if(stockCode != ""){
			url += "&txtDMorJC=";
			url += stockCode;
		}

		pageno := 1;
		if _, ok := aid["pageno"]; ok {
			pageno = aid["pageno"].(int);
		}
		url += "&PAGENO=";
		url += strconv.Itoa(pageno);

		url += "&random=";
		url += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url: url,
			Rule: k,
			Temp: map[string]interface{}{"time": tm,"name": k,"v":v,"pageno":pageno,"stockCode":stockCode},
		});

		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		pageno := ctx.GetTemp("pageno",0).(int);
		stockCode := ctx.GetTemp("stockCode","");

		tm := ctx.GetTemp("time","").(string);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + tm + "_" + strconv.Itoa(pageno) + ".html",ctx.GetText());

		//拉取历史记录
		if(history_ShenZhengBond[k] == true){
			if(pageno == 1){
				loc, _ := time.LoadLocation("Local")
				t ,_ := time.ParseInLocation("2006-01-02",tm,loc);

				t_ := OpeningQuotationDay(LastDay(t));
				if(ValidDay(t_) == false){
					return;
				}
				tm_ := GetDayString(t_);

				ctx.Aid(map[string]interface{}{"time": tm_,"name": k,"stockCode":stockCode}, k)
			}
		}

		pagecount := gjson.Get(ctx.GetText(),"0.metadata.pagecount").Int();
		fmt.Println(pagecount);
		if(int64(pageno + 1) > pagecount){
			return;
		}

		ctx.Aid(map[string]interface{}{"time": tm,"name": k,"pageno":pageno + 1,"stockCode":stockCode}, k)
	},
}

//历史股票数据
var DataRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);
		// tm := aid["time"].(string);
		cycle := aid["cycle"].(string);
		code := aid["code"].(string);

		url := v;
		url += "&cycleType=";
		url += cycle;
		url += "&code=";
		url += code;

		url += "&random=";
		url += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url: url,
			Rule: k,
			Temp: aid,
		});

		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		// v := ctx.GetTemp("v","").(string);
		tm := ctx.GetTemp("time","").(string);
		code := ctx.GetTemp("code","").(string);
		cycle := ctx.GetTemp("cycle","").(string);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + code + "_" + cycle + "_" + tm + ".html",ctx.GetText());
	},
}

//分时
var TimeShareRule = &Rule{
	AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
		k := aid["name"].(string);
		v := aid["v"].(string);

		// tm := aid["time"].(string);
		code := aid["code"].(string);

		url := v;
		url += "&code=";
		url += code;

		url += "&random=";
		url += strconv.FormatFloat(rand.Float64(),'f',16,64);

		ctx.AddQueue(&request.Request{
			Url: url,
			Rule: k,
			Temp: aid,
		});

		return nil;
	},
	ParseFunc:func(ctx * Context){
		k := ctx.GetTemp("name","").(string);
		tm := ctx.GetTemp("time","").(string);
		code := ctx.GetTemp("code","").(string);

		dir := DATA + k;
		os.MkdirAll(dir,os.ModeDir);
		cache(dir + "/" + code + "_" + tm + ".html",ctx.GetText());
	},
};

var ShenZhengBond = &Spider{
    Name: "深圳证券交易所",
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

            for k,v:= range rss_ShenZhengBond {
				ctx.SetTimer(k, time.Minute*5, nil)
				ctx.Aid(map[string]interface{}{"time": tm,"name": k,"v":v}, k)
			}
        },
        Trunk: map[string]*Rule{
            "市场总貌-证券类别统计":TimeRule,
			"上市公司-上市公司列表":PageNORule,
			"上市公司-上市公司更名-全称":PageNORule,
			"上市公司-上市公司更名-简称":PageNORule,
			"上市公司-暂停上市公司":PageNORule,
			"上市公司-终止上市公司":PageNORule,
			"股票-A股列表":PageNORule,
			"股票-B股列表":PageNORule,
			"股票-A＋B股列表":PageNORule,
			"股票-基本指数-深圳市场":TimeRule,
			"股票-基本指数-深市主板":TimeRule,
			"股票-基本指数-中小企业板":TimeRule,
			"股票-基本指数-创业板":TimeRule,
			"股票-指标排名":{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					k := aid["name"].(string);
					v := rss_ShenZhengBond[k];
					tm := aid["time"].(string);
					loop := history_ShenZhengBond[k];

					for ks,vs := range select_index{
						for km,vm := range module_index {
							url := v + "&TABKEY=" + vm + "&selectModule=" + vs;
							name := k + "-" + ks + "-" + km;
							ctx.Aid(map[string]interface{}{"time": tm,"name": name,"v":url,"loop":loop}, name);
						}
					}

					return nil;
				},
			},
			"股票-指标排名-全部-总股本":Time2Rule,
			"股票-指标排名-全部-总市值":Time2Rule,
			"股票-指标排名-全部-流通股本":Time2Rule,
			"股票-指标排名-全部-流通市值":Time2Rule,
			"股票-指标排名-全部-市盈利":Time2Rule,

			"股票-指标排名-主板-总股本":Time2Rule,
			"股票-指标排名-主板-总市值":Time2Rule,
			"股票-指标排名-主板-流通股本":Time2Rule,
			"股票-指标排名-主板-流通市值":Time2Rule,
			"股票-指标排名-主板-市盈利":Time2Rule,

			"股票-指标排名-中小企业板-总股本":Time2Rule,
			"股票-指标排名-中小企业板-总市值":Time2Rule,
			"股票-指标排名-中小企业板-流通股本":Time2Rule,
			"股票-指标排名-中小企业板-流通市值":Time2Rule,
			"股票-指标排名-中小企业板-市盈利":Time2Rule,

			"股票-指标排名-创业板-总股本":Time2Rule,
			"股票-指标排名-创业板-总市值":Time2Rule,
			"股票-指标排名-创业板-流通股本":Time2Rule,
			"股票-指标排名-创业板-流通市值":Time2Rule,
			"股票-指标排名-创业板-市盈利":Time2Rule,

			"股票-行业统计-全部":TimeRule,
			"股票-行业统计-主板":TimeRule,
			"股票-行业统计-中小企业版":TimeRule,
			"股票-行业统计-创业板":TimeRule,

			"股票-行业市盈利":TimeRule,

			"行情-历史行情":Time3Rule,
			"行情-历史数据":&Rule{
				AidFunc:func(ctx * Context,aid map[string]interface{})interface{}{
					k := aid["name"].(string);
					tm := aid["time"].(string);

					code := "000001";
					v:="http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1";
					ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-分时","v":v,"code":code}, k + "-分时");

					v = rss_ShenZhengBond[k];
					ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-日","v":v,"code":code,"cycle":"32"}, k + "-日");
					ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-周","v":v,"code":code,"cycle":"33"}, k + "-周");
					ctx.Aid(map[string]interface{}{"time": tm,"name": k + "-月","v":v,"code":code,"cycle":"34"}, k + "-月");
					return nil;
				},
			},
			"行情-历史数据-分时":TimeShareRule,
			"行情-历史数据-日":DataRule,
			"行情-历史数据-周":DataRule,
			"行情-历史数据-月":DataRule,
        },
    },
}
