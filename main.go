package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// input: 和訳したい記事のURL, output: 和訳したい記事のhtmlファイル
func httpRequest(url *string) []byte {
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	//fmt.Printf("%s", b) // type of b is byte
	return b

}

// input: 和訳したい記事のhtml, output:<title>, <body>を含んだStrings
//func parse()

// input: Strings, output: 和訳が含まれるJSON
//func deeplRequest()

// input: 和訳が含まれるJSON、output: htmlファイル
//func produceHtml()

func main() {
	//var url string = "https://news.crunchbase.com/news/what-this-years-seed-funding-tells-us-about-the-startup-future/"
	var url string = "https://yorkn.info"
	var txByte []byte = httpRequest(&url)

	fmt.Printf("%s", txByte)

}
