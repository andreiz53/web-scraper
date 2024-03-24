package types

import (
	"gorm.io/gorm"
)

type Website struct {
	gorm.Model
	URL            string
	UserID         uint
	User           User
	Robots         bool
	Keywords       string
	Description    string
	HasSitemap     bool
	HasOpenGraph   bool
	HasTwitterCard bool
	HasDublinCore  bool
	HasJsonLd      bool
	HasMetaPixel   bool
}

func NewWebsite() *Website {
	return &Website{
		Robots:         false,
		Keywords:       "",
		Description:    "",
		HasSitemap:     false,
		HasOpenGraph:   false,
		HasTwitterCard: false,
		HasDublinCore:  false,
		HasJsonLd:      false,
		HasMetaPixel:   false,
	}
}
