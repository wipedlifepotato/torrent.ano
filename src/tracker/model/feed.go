package model

import (
	"encoding/xml"
	"net/url"
	"time"
)

type AtomFeed struct {
	Domain   string
	BaseURL  *url.URL
	Torrents []Torrent
	Title    string
	ID       string
}

type Link struct {
	URL string `xml:"href,attr"`
}

func NewLink(domain, path string) Link {
	u := &url.URL{
		Scheme: "http",
		Host:   domain,
		Path:   path,
	}
	return Link{
		URL: u.String(),
	}
}

type AtomFeedImpl struct {
	Title    string    `xml:"title"`
	SubTitle string    `xml:"subtitle"`
	ID       string    `xml:"id"`
	Link     Link      `xml:"link"`
	Updated  time.Time `xml:"updated"`
	Torrents []Torrent `xml:"entry"`
}

func (feed *AtomFeed) toFeed() *AtomFeedImpl {

	latest := time.Unix(0, 0)
	for _, t := range feed.Torrents {
		u := t.UploadedAt()
		if u.After(latest) {
			latest = u
		}
	}

	u := feed.BaseURL

	return &AtomFeedImpl{
		Title:    feed.Title,
		SubTitle: feed.Title,
		Link:     NewLink(feed.Domain, u.RequestURI()),
		ID:       feed.ID,
		Torrents: feed.Torrents,
		Updated:  latest,
	}
}

func (feed *AtomFeed) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = "feed"
	start.Name.Space = "http://www.w3.org/2005/Atom"
	err = e.EncodeElement(feed.toFeed(), start)
	return
}