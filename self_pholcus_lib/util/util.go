package self_pholcus_lib

import (
    // "math"
    "fmt"
	"time"
    "io/ioutil"
    "strings"
)

//缓存数据
func Cache(path string,data string){
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

//创建交易所日期
func GreaterByDay(t time.Time,limit map[string]int)bool{
	if(t.Year() < limit["year"]){
		return false;
	}else if(t.Year() == limit["year"]){
		if(int(t.Month()) < limit["month"]){
			return false;
		}
	}
	return true;
}

func TrimLeft(s string,count int)string{
    index := 0;
    s = strings.TrimLeftFunc(s,func(_ rune) bool{
        index++;
        if(index > count){
            return false;
        }
        return true;
    });
    return s;
}

func TrimRight(s string,count int)string{
    index := 0;
    s = strings.TrimRightFunc(s,func(_ rune) bool{
        index++;
        if(index > count){
            return false;
        }
        return true;
    });
    return s;
}
