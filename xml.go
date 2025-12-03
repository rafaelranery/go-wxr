package wxr

import "encoding/xml"

// wxr represents the WordPress eXtended RSS structure (internal).
type wxr struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	BaseSiteURL string     `xml:"http://wordpress.org/export/1.2/ base_site_url"`
	BaseBlogURL string     `xml:"http://wordpress.org/export/1.2/ base_blog_url"`
	Items       []item     `xml:"item"`
	Authors     []wpAuthor `xml:"http://wordpress.org/export/1.2/ author"`
}

type item struct {
	Title           string       `xml:"title"`
	Link            string       `xml:"link"`
	GUID            string       `xml:"guid"`
	PubDate         string       `xml:"pubDate"`
	DCCreator       string       `xml:"http://purl.org/dc/elements/1.1/ creator"`
	ContentEncoded  string       `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	ExcerptEncoded  string       `xml:"http://wordpress.org/export/1.2/excerpt/ encoded"`
	PostID          int          `xml:"http://wordpress.org/export/1.2/ post_id"`
	PostDate        string       `xml:"http://wordpress.org/export/1.2/ post_date"`
	PostDateGMT     string       `xml:"http://wordpress.org/export/1.2/ post_date_gmt"`
	PostModified    string       `xml:"http://wordpress.org/export/1.2/ post_modified"`
	PostModifiedGMT string       `xml:"http://wordpress.org/export/1.2/ post_modified_gmt"`
	PostParent      int          `xml:"http://wordpress.org/export/1.2/ post_parent"`
	PostName        string       `xml:"http://wordpress.org/export/1.2/ post_name"`
	PostType        string       `xml:"http://wordpress.org/export/1.2/ post_type"`
	Status          string       `xml:"http://wordpress.org/export/1.2/ status"`
	AttachmentURL   string       `xml:"http://wordpress.org/export/1.2/ attachment_url"`
	PostMeta        []postMeta   `xml:"http://wordpress.org/export/1.2/ postmeta"`
	Categories      []wpCategory `xml:"category"`
}

type wpCategory struct {
	Domain   string `xml:"domain,attr"`
	Value    string `xml:",chardata"`
	NiceName string `xml:"nicename,attr"`
}

type wpAuthor struct {
	ID          int    `xml:"http://wordpress.org/export/1.2/ author_id"`
	Login       string `xml:"http://wordpress.org/export/1.2/ author_login"`
	DisplayName string `xml:"http://wordpress.org/export/1.2/ author_display_name"`
}

type postMeta struct {
	Key   string `xml:"http://wordpress.org/export/1.2/ meta_key"`
	Value string `xml:"http://wordpress.org/export/1.2/ meta_value"`
}
