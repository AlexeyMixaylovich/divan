package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type divanData struct {
	url, name, sale string
	price, oldPrice int
}

func main() {
	res, err := getHtml()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
	}
	fmt.Println(res)

}

func getHtml() ([]divanData, error) {
	url := "https://www.divan.ru/category/stok-mebeli?categories[]=2"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("пирсинг %s: as HTML: %v", url, err)
	}
	divansHtml := getDivansHtml(nil, doc)

	return getDivansDataFromHtml(divansHtml), nil

}

func getDivansHtml(items []*html.Node, n *html.Node) []*html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "lsooF") {
				items = append(items, n)
			}
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		items = getDivansHtml(items, c)
	}

	return items
}

func getDivansDataFromHtml(items []*html.Node) []divanData {
	var divans []divanData
	for _, divanHtml := range items[0:1] {
		data := getDivanDataFromHtml(divanHtml)

		divans = append(divans, data)

	}
	return divans
}

func getDivanDataFromHtml(divanHtml *html.Node) divanData {
	var data divanData

	//  const name = _.get(element, 'childNodes[0].childNodes[0].text')?.trim();
	// const href = _.get(element, 'childNodes[0].attrs.href');
	// const price = _.get(element, 'childNodes[1].childNodes[0].childNodes[0].text')?.trim();
	// const oldPrice = _.get(element, 'childNodes[1].childNodes[1].childNodes[0].text')?.trim();
	// const sale = _.get(element, 'childNodes[1].childNodes[2].childNodes[0].childNodes[0].text')?.trim();
	for _, a := range divanHtml.FirstChild.Attr {
		if a.Key == "href" {
			data.url = a.Val
		}
	}
	data.name = divanHtml.FirstChild.FirstChild.Data
	// secondChild := divanHtml.FirstChild.NextSibling

	fmt.Println("name")
	fmt.Println(divanHtml.FirstChild.FirstChild.Data)
	fCh := divanHtml.FirstChild
	fmt.Println(fCh.NextSibling)

	fmt.Println("name2")
	// if secondFirst == nil {
	data.price = 0
	data.oldPrice = 0
	data.sale = "0%"
	// } else {
	// 	data.price = getPrice(secondChild.FirstChild.FirstChild.Data)
	// 	data.oldPrice = getPrice(secondChild.FirstChild.NextSibling.FirstChild.Data)
	// 	data.sale = secondChild.FirstChild.NextSibling.NextSibling.FirstChild.FirstChild.Data
	// }

	return data
}

func getPrice(n string) int {

	nTrim := strings.ReplaceAll(n, " ", "")
	res, err := strconv.Atoi(nTrim)
	if err == nil {
		return res
	} else {
		fmt.Println(err)
		return 0
	}

}
