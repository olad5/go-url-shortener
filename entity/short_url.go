package entity

type ShortenUrl struct {
	ShortUrl    string `bson:"short_url"  json:"short_url"`
	OriginalUrl string `bson:"original_url" json:"original_url"`
	ClickCount  int    `bson:"click_count" json:"click_count"`
	UniqueId    int    `bson:"unique_id" json:"unique_id"`
}
