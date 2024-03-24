package scraper

import (
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	URL       string
	Collector colly.Collector
	Data      WebsiteData
}

type ExternalLink struct {
	Text         string
	Link         string
	OpenInNewTab bool
	Relationship string
}

type WebsiteData struct {
	Robots         bool
	Keywords       []string
	Description    string
	HasSitemap     bool
	HasOpenGraph   bool
	HasTwitterCard bool
	HasDublinCore  bool
	HasJsonLd      bool
	HasMetaPixel   bool
	ExternalLinks  []ExternalLink
}

func NewScraper(URL string) (*Scraper, error) {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return nil, err
	}
	return &Scraper{
		URL:       URL,
		Collector: *colly.NewCollector(),
		Data:      WebsiteData{},
	}, nil
}

func (s *Scraper) SetSEO() {
	s.GetRobots()
	s.GetKeywords()
	s.GetDescription()
	s.HasSitemapLinked()
	s.HasOpenGraph()
	s.HasTwitterCard()
	s.HasDublinCore()
	s.hasJsonLd()
	s.GetExternalLinks()
	s.HasMetaPixel()
}

func (s *Scraper) GetRobots() {
	s.Collector.OnHTML("head meta[name='robots']", func(h *colly.HTMLElement) {
		if strings.Contains(h.Attr("content"), "index") && !strings.Contains(h.Attr("content"), "noindex") {
			s.Data.Robots = true
		}
	})
}

func (s *Scraper) GetKeywords() {
	s.Collector.OnHTML("head meta[name='keywords']", func(h *colly.HTMLElement) {
		if h.Attr("content") != "" {
			s.Data.Keywords = strings.Split(h.Attr("content"), ",")
		}
	})
}

func (s *Scraper) GetDescription() {
	s.Collector.OnHTML("head meta[name='description']", func(h *colly.HTMLElement) {
		s.Data.Description = h.Attr("content")
	})
}

func (s *Scraper) HasSitemapLinked() {
	s.Collector.OnHTML("head link[rel='alternate']", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		if strings.Contains(href, "sitemap.xml") {
			s.Data.HasSitemap = true
		}
	})
}

func (s *Scraper) HasOpenGraph() {
	s.Collector.OnHTML("head meta[property]", func(h *colly.HTMLElement) {
		property := h.Attr("property")
		if strings.Contains(property, "og:") {
			s.Data.HasOpenGraph = true
			return
		}
	})
}

func (s *Scraper) HasTwitterCard() {
	s.Collector.OnHTML("head meta[property]", func(h *colly.HTMLElement) {
		property := h.Attr("property")
		if strings.Contains(property, "twitter:") {
			s.Data.HasTwitterCard = true
			return
		}
	})
}

func (s *Scraper) HasDublinCore() {
	s.Collector.OnHTML("head meta[name]", func(h *colly.HTMLElement) {
		name := h.Attr("name")
		if strings.Contains(name, "dc.") {
			s.Data.HasDublinCore = true
			return
		}
	})
}

func (s *Scraper) hasJsonLd() {
	s.Collector.OnHTML("script[type='application/ld+json']", func(h *colly.HTMLElement) {
		if h.Text != "" {
			s.Data.HasJsonLd = true
		}

	})
}

func (s *Scraper) GetExternalLinks() {
	s.Collector.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link, err := url.Parse(h.Attr("href"))
		if err != nil {
			return
		}
		local_link, err := url.Parse(s.URL)
		if err != nil {
			return
		}

		if !strings.Contains(link.Hostname(), local_link.Hostname()) && link.Host != "" {
			new_link := ExternalLink{
				Link:         h.Attr("href"),
				OpenInNewTab: h.Attr("target") == "_blank",
				Text:         strings.TrimSpace(h.Text),
				Relationship: h.Attr("rel"),
			}

			s.Data.ExternalLinks = append(s.Data.ExternalLinks, new_link)
		}

	})
}

func (s *Scraper) HasMetaPixel() {
	s.Collector.OnHTML("meta[property='fb:app_id']", func(h *colly.HTMLElement) {
		if h.Attr("content") != "" {
			s.Data.HasMetaPixel = true
		}
	})
}
