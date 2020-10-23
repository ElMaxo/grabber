package grabber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var lentaQuery = Query{
	ItemsSelector: Selector{
		XPath: "//div[@class='span4']/div[@class='item']/a",
	},
	TitleSelector: Selector{
		Value: Text,
	},
	LinkSelector: Selector{
		Value: Attribute,
		Attr:  "href",
	},
	DescriptionSelector: Selector{
		XPath: "//div[@itemprop='articleBody']",
		Value: InnerText,
	},
	FollowLinkForDescription: true,
}

var ramblerQuery = Query{
	ItemsSelector: Selector{
		XPath: "//div[@class='top-main__news-item']",
	},
	TitleSelector: Selector{
		XPath: "//div[@class='top-card__title']",
		Value: Text,
	},
	LinkSelector: Selector{
		XPath: "//a",
		Value: Attribute,
		Attr:  "href",
	},
	DescriptionSelector: Selector{
		XPath: "//meta[@itemprop='articleBody']",
		Value: Attribute,
		Attr:  "content",
	},
	FollowLinkForDescription: true,
}

func TestHtmlGrabber_GrabNews(t *testing.T) {
	grab := NewHtmlGrabber()
	lentaArticles, err := grab.GrabNews("https://lenta.ru", lentaQuery)
	require.NoError(t, err)
	require.NotEmpty(t, lentaArticles)
	require.NotEmpty(t, lentaArticles[0].Title)
	require.NotEmpty(t, lentaArticles[0].Link)

	ramblerArticles, err := grab.GrabNews("https://news.rambler.ru", ramblerQuery)
	require.NoError(t, err)
	require.NotEmpty(t, ramblerArticles)
	require.NotEmpty(t, ramblerArticles[0].Title)
	require.NotEmpty(t, ramblerArticles[0].Link)
}
