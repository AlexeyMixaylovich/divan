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

	return getDivanDataFromHtml(divansHtml), nil

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

func getDivanDataFromHtml(items []*html.Node) []divanData {
	var divans []divanData
	for _, divanHtml := range items {
		var data divanData

		for _, a := range divanHtml.FirstChild.Attr {
			if a.Key == "href" {
				data.url = a.Val
			}
		}
		data.name = divanHtml.FirstChild.FirstChild.Data
		secondChild := divanHtml.FirstChild.NextSibling
		data.price = getPrice(secondChild.FirstChild.FirstChild.Data)
		fmt.Println(secondChild.FirstChild.NextSibling.FirstChild.Data)
		data.oldPrice = getPrice(secondChild.FirstChild.NextSibling.FirstChild.Data)

		data.sale = secondChild.FirstChild.NextSibling.NextSibling.FirstChild.FirstChild.Data

		divans = append(divans, data)

	}
	return divans
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
