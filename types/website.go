package types

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
