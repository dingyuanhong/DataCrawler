package main

import (
    "fmt"
    // "io"
    // "strings"
    "github.com/anaskhan96/soup"
    "github.com/PuerkitoBio/goquery"
    "github.com/henrylee2cn/pholcus/exec"
    _ "./self_pholcus_lib"
)

func soupTest(){
    fmt.Println("soup:");
    response,err := soup.Get("http://www.baidu.com");
    // fmt.Println(response);
    fmt.Println(err);
    doc := soup.HTMLParse(response);
    div := doc.Find("meta");
    fmt.Println(div);
}

func goqueryTest(){
    fmt.Println("goquery:");
    doc ,error := goquery.NewDocument("http://www.baidu.com");
    fmt.Println(error);
    doc.Find("div").EachWithBreak(func(i int,s * goquery.Selection) bool{
        if(s.Length() > 1){
            s.EachWithBreak(func(j int,sc *goquery.Selection)bool{
                fmt.Println("b:",sc.Size());
                str,_ := sc.Html();
                fmt.Println("a:",str);
                return true;
            })
        }else{
            str,_ := s.Html();
            fmt.Println("c:",str);
        }

        return true;
    })
}

func main(){
    // soupTest();
    // goqueryTest();

    exec.DefaultRun("web")
}
