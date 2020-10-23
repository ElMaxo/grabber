package main

import (
	"flag"
	"fmt"
	"grabber/internal/grabber"
	"log"
	"os"
)

var lentaQuery = grabber.Query{
	ItemsSelector: grabber.Selector{
		XPath: "//div[@class='span4']/div[@class='item']/a",
	},
	TitleSelector: grabber.Selector{
		Value: grabber.Text,
	},
	LinkSelector: grabber.Selector{
		Value: grabber.Attribute,
		Attr:  "href",
	},
	DescriptionSelector: grabber.Selector{
		XPath: "//div[@itemprop='articleBody']",
		Value: grabber.InnerText,
	},
	FollowLinkForDescription: true,
}

var ramblerQuery = grabber.Query{
	ItemsSelector: grabber.Selector{
		XPath: "//div[@class='top-main__news-item']",
	},
	TitleSelector: grabber.Selector{
		XPath: "//div[@class='top-card__title']",
		Value: grabber.Text,
	},
	LinkSelector: grabber.Selector{
		XPath: "//a",
		Value: grabber.Attribute,
		Attr:  "href",
	},
	DescriptionSelector: grabber.Selector{
		XPath: "//meta[@itemprop='articleBody']",
		Value: grabber.Attribute,
		Attr:  "content",
	},
	FollowLinkForDescription: true,
}

func main() {
	var url string
	//flag.StringVar(&url, "url", "https://lenta.ru", "News page URL")
	flag.StringVar(&url, "url", "https://news.rambler.ru/", "News page URL")
	flag.Parse()
	if url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	grab := grabber.NewHtmlGrabber()
	articles, err := grab.GrabNews(url, ramblerQuery)
	if err != nil {
		log.Fatal("failed to grab articles: ", err)
	}
	for _, article := range articles {
		fmt.Println("TITLE: ", article.Title)
		fmt.Println("LINK: ", article.Link)
		fmt.Println("DESC: ", article.Description)
	}
}
