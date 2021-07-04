package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DocTitleBody struct {
	title string
	body  string
}

// input: 和訳したい記事のURL, output: 和訳したい記事のhtmlファイル
func getDocument(url *string) DocTitleBody {
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var od DocTitleBody
	od.title = doc.Find("title").Text()
	od.body = doc.Find("div > div.col-lg-9.col-md-9.col-mod-single.col-mod-main > div.row > div.col-lg-10.col-md-10.col-sm-10").Text()

	resp.Body.Close()
	return od

}

//input: 和訳したい記事のtitle, body, output:<title>の中身, <body>（メイン文書の中身）をextractしたStrings
func parse(od DocTitleBody) DocTitleBody {

	r := regexp.MustCompile(`Subscribe to the Crunchbase Daily`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`Shares`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`Email`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`Facebook`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`Twitter`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`LinkedIn`)
	od.body = r.ReplaceAllString(od.body, " ")
	r = regexp.MustCompile(`\s{2}`)
	od.body = r.ReplaceAllString(od.body, " ")

	fmt.Println(od.body)

	r = regexp.MustCompile(`– Crunchbase News`)
	od.title = r.ReplaceAllString(od.title, "")

	od.title = strings.Replace(od.title, " ", "%20", -1)
	od.body = strings.Replace(od.title, " ", "%20", -1)
	return od

}

// input: Strings, output: 和訳が含まれるJSON
func deeplPost(sentences string) string {

	url := "https://api-free.deepl.com/v2/translate?auth_key=a0474d62-bbcf-5b7d-c298-da9bbac4ecdd:fx&text=" + sentences + "&target_lang=JA"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ここよ" + string(body))
	return string(body)
}

func spritBody(bc string) []string {
	var sslice []string = strings.Split(bc, ".")
	return sslice
}

// input: 和訳が含まれるJSON、output: htmlファイル
//func produceHtml()

func main() {
	var url string = "https://news.crunchbase.com/news/when-90-is-young-what-a-moonshot-vc-thinks-about-radical-longevity/"
	//var url string = "https://yorkn.info"
	var od DocTitleBody = getDocument(&url)
	od = parse(od)
	var td DocTitleBody // 翻訳後のtitle, body

	fmt.Printf("title is %s\n", od.title)
	fmt.Printf("body is %s\n", td.title)

	td.title = deeplPost(od.title)

	fmt.Println(td.title)
	//td.body = deeplPost(od.body)

}
