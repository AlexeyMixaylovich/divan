package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var domain = "https://www.divan.ru"

type divanData struct {
	url, name, sale, art string
	price, oldPrice      int
}

func getDivans() []divanData {
	var divans []divanData

	paths := []string{
		"/category/stok-mebeli?categories[]=2",
		"/category/stok-mebeli/page-2?categories[]=2",
		"/category/stok-mebeli/page-3?categories[]=2",
	}
	for _, path := range paths {
		res, err := getHtml(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}
		divans = append(divans, res...)
	}
	return divans
}

func getHtml(path string) ([]divanData, error) {
	url := domain + path
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
	for _, divanHtml := range items {
		data := getDivanDataFromHtml(divanHtml)

		divans = append(divans, data)

	}
	return divans
}

func getDivanDataFromHtml(divanHtml *html.Node) divanData {
	var data divanData

	for _, a := range divanHtml.FirstChild.Attr {
		if a.Key == "href" {
			data.url = domain + a.Val
			_, art, found := strings.Cut(a.Val, "art--")
			if found {
				data.art = art
			}
		}
	}

	data.name = divanHtml.FirstChild.FirstChild.Data
	thirdChild := divanHtml.FirstChild.NextSibling.NextSibling

	if thirdChild == nil {
		data.price = 0
		data.oldPrice = 0
		data.sale = "0%"
	} else {
		data.price = getPrice(thirdChild.FirstChild.FirstChild.Data)
		saleNode := thirdChild.FirstChild.NextSibling
		if saleNode == nil {
			data.oldPrice = 0
			data.sale = "0%"
		} else {
			data.oldPrice = getPrice(saleNode.FirstChild.Data)
			data.sale = saleNode.NextSibling.FirstChild.FirstChild.Data
		}
	}

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
