package grabber

import (
	"bytes"
	"runtime"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type ValueType string

const (
	Attribute = ValueType("attr")
	Text      = ValueType("text")
	InnerText = ValueType("inner_text")
)

type Article struct {
	Link        string
	Title       string
	Description string
}
type Selector struct {
	XPath string    `json:"xpath"`
	Value ValueType `json:"value"`
	Attr  string    `json:"attr"`
}

type Query struct {
	ItemsSelector            Selector `json:"itemsSelector"`
	TitleSelector            Selector `json:"titleSelector"`
	LinkSelector             Selector `json:"linkSelector"`
	DescriptionSelector      Selector `json:"descriptionSelector"`
	FollowLinkForDescription bool     `json:"followLinkForDescription"`
}

type asyncDescJob struct {
	id       int
	url      string
	selector Selector
}

type asyncDescResult struct {
	id   int
	desc string
	err  error
}

type Grabber interface {
	GrabNews(url string, query Query) ([]Article, error)
}

type htmlGrabber struct {
}

func NewHtmlGrabber() Grabber {
	return &htmlGrabber{}
}

func (hg *htmlGrabber) GrabNews(url string, query Query) ([]Article, error) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}
	items, err := htmlquery.QueryAll(doc, query.ItemsSelector.XPath)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []Article{}, nil
	}
	articles := make(map[int]Article)
	for i, item := range items {
		article := Article{}
		title, err := nodeValue(item, query.TitleSelector)
		if err != nil {
			return nil, err
		}
		article.Title = strings.TrimSpace(title)
		link, err := nodeValue(item, query.LinkSelector)
		if err != nil {
			return nil, err
		}
		article.Link = articleURL(url, link)
		if !query.FollowLinkForDescription {
			article.Description, err = nodeValue(item, query.DescriptionSelector)
			if err != nil {
				return nil, err
			}
		}
		articles[i] = article
	}

	articlesWithDesc := make([]Article, 0, len(items))
	if query.FollowLinkForDescription {
		jobsChan := make(chan asyncDescJob, len(items))
		resChan := make(chan asyncDescResult, len(items))
		for i := 0; i < runtime.NumCPU(); i++ {
			go asyncDescWorker(jobsChan, resChan)
		}

		for i := 0; i < len(items); i++ {
			jobsChan <- asyncDescJob{
				id:       i,
				url:      articles[i].Link,
				selector: query.DescriptionSelector,
			}
		}
		close(jobsChan)

		for i := 0; i < len(items); i++ {
			res := <-resChan
			if res.err != nil {
				continue
			}
			article := articles[res.id]
			article.Description = res.desc
			articlesWithDesc = append(articlesWithDesc, article)
		}
		close(resChan)
	}
	return articlesWithDesc, nil
}

func asyncDescWorker(jobsChan <-chan asyncDescJob, resChan chan<- asyncDescResult) {
	for job := range jobsChan {
		description, err := getDescriptionByLink(job.url, job.selector)
		resChan <- asyncDescResult{
			id:   job.id,
			desc: description,
			err:  err,
		}
	}
}

func getDescriptionByLink(url string, selector Selector) (string, error) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return "", err
	}
	return nodeValue(doc, selector)
}

func articleURL(baseURL string, url string) string {
	if !strings.HasPrefix(url, "http") {
		return strings.TrimSuffix(baseURL, "/") + url
	}
	return url
}

func nodeValue(node *html.Node, selector Selector) (string, error) {
	valNode := node
	var err error
	if selector.XPath != "" {
		valNode, err = htmlquery.Query(node, selector.XPath)
		if err != nil {
			return "", err
		}
	}
	if valNode == nil {
		return "", nil
	}
	switch selector.Value {
	case Attribute:
		return htmlquery.SelectAttr(valNode, selector.Attr), nil
	case Text:
		for child := valNode.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data, nil
			}
		}
	case InnerText:
		return innerText(valNode), nil
	}
	return "", nil
}

func innerText(n *html.Node) string {
	var output func(*bytes.Buffer, *html.Node)
	output = func(buf *bytes.Buffer, n *html.Node) {
		if n.Data == "style" || n.Data == "script" || n.Data == "div" {
			return
		}
		switch n.Type {
		case html.TextNode:
			buf.WriteString(n.Data)
			return
		case html.CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}

	var buf bytes.Buffer
	output(&buf, n)
	return buf.String()
}
