package structs

type TrendingNowResponse struct {
	Default Default `json:"default"`
}
type Title struct {
	Query       string `json:"query"`
	ExploreLink string `json:"exploreLink"`
}
type Image struct {
	NewsURL  string `json:"newsUrl"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
}
type Articles struct {
	Title   string `json:"title"`
	TimeAgo string `json:"timeAgo"`
	Source  string `json:"source"`
	Image   Image  `json:"image"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}
type TrendingSearches struct {
	Title            Title      `json:"title"`
	FormattedTraffic string     `json:"formattedTraffic"`
	RelatedQueries   []any      `json:"relatedQueries"`
	Image            Image      `json:"image"`
	Articles         []Articles `json:"articles"`
	ShareURL         string     `json:"shareUrl"`
}
type TrendingSearchesDays struct {
	Date             string             `json:"date"`
	FormattedDate    string             `json:"formattedDate"`
	TrendingSearches []TrendingSearches `json:"trendingSearches"`
}
type Default struct {
	TrendingSearchesDays  []TrendingSearchesDays `json:"trendingSearchesDays"`
	EndDateForNextRequest string                 `json:"endDateForNextRequest"`
	RssFeedPageURL        string                 `json:"rssFeedPageUrl"`
}




