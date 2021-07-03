package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

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

	r = regexp.MustCompile(`– Crunchbase News`)
	od.title = r.ReplaceAllString(od.title, "")

	return od

}

// input: Strings, output: 和訳が含まれるJSON
// 以下は適当に貼り付けしただけだから、Deeplのhttp reqの形に合わせて変更する
func HttpPost(url, token, device string) error {
	url =
	values := url.Values{}
	values.Set("token", token)
	values.Add("device", device)

	req, err := http.NewRequest(
			"POST",
			url,
			strings.NewReader(values.Encode()),
	)
	if err != nil {
			return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
			return err
	}
	defer resp.Body.Close()

	return err
}

// input: 和訳が含まれるJSON、output: htmlファイル
//func produceHtml()

func main() {
	var url string = "https://news.crunchbase.com/news/what-this-years-seed-funding-tells-us-about-the-startup-future/"
	//var url string = "https://yorkn.info"
	var od DocTitleBody = getDocument(&url)
	od = parse(od)

	//var body string = parse(txByte)
	fmt.Printf("title is %s\n", od.title)
	fmt.Printf("body is %s\n", od.body)

}
