package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayush6624/go-chatgpt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/Sonmezturk/telegram-bot-formatter/structs"
)


var OPENAI_KEY string

type ArticleInfo struct {
	Title string
	URL   string
}

type ConsolidatedTrend struct {
	TrendingSearchTitle string
	Article             ArticleInfo
}


func trend(bot *tgbotapi.BotAPI) {
	err := godotenv.Load()
	OPENAI_KEY = os.Getenv("OPENAI_KEY")

	if OPENAI_KEY == "" {
		log.Fatalln("OPENAI KEY NOT FOUND")
	}

	url := "https://trends.google.com/trends/api/dailytrends?hl=en-US&tz=-180&geo=US&hl=en-US&ns=15"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making the request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading the response body:", err)
	}

	bodyStr := string(body)
	bodyStr = strings.TrimPrefix(bodyStr, ")]}',")
	var response structs.TrendingNowResponse
	err = json.Unmarshal([]byte(bodyStr), &response)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	trends := consolidateTrends(response)

	for _, trend := range trends[0:6] {
		fmt.Printf("Trending Search: %s\t Article Title: %s\t Article Url: %s\n", trend.TrendingSearchTitle, trend.Article.Title, trend.Article.URL)
		//articleContent := fetchArticleContent(trend.Article.URL, "test")
		//summary := getArticleSummary(trend.TrendingSearchTitle, articleContent)
		//fmt.Println(summary)
		//sendMessageToChannel(bot, "channelID", fmt.Sprintf(`Trend keyword %s,   Idea: %s`, trend.TrendingSearchTitle, summary))
		//image, err := dallE.GenerateImageGenerateImage(summary, 1, "1024x1024")
		if err != nil {
			log.Panic(err)
		}
		//fmt.Println(image)

		//channelID, err := strconv.ParseInt("-1002125172322", 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		//sendImageToChannel(bot, channelID, image.Data[0].URL)
	}
}

func getArticleSummary(articleContent string, keyword string) string {

	template := `Analyze the article and extract its main theme, message, or any unique elements it presents. Based on these findings, generate a creative idea for a t-shirt design that represents the essence or most striking aspect of the article but stick to the keyword: %s. The design idea should be engaging and capture the article's spirit in a visually appealing way.
	The idea needs to be less than 1000 character and there should be only T-shirt design in your output.
	Example Output:
	
	Article on environmental conservation: 'T-shirt design: An artistic depiction of the Earth with a heart in the center, surrounded by diverse flora and fauna. The caption reads, "Love Your Mother Earth – Conserve and Preserve."'
	
	Article on technological advancements: 'T-shirt design: A futuristic robot hand reaching out to a human hand, symbolizing unity between humans and technology. The slogan, "Building a Better Tomorrow, Together," is featured below the image
	
	%s`
	result := fmt.Sprintf(template, keyword, articleContent)
	return gpt(result)
}

func gpt(message string) string {
	client, err := chatgpt.NewClient(OPENAI_KEY)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	res, err := client.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: message,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return res.Choices[0].Message.Content

}

func consolidateTrends(response structs.TrendingNowResponse) []ConsolidatedTrend {
	var consolidatedTrends []ConsolidatedTrend

	for _, day := range response.Default.TrendingSearchesDays {
		for _, trendingSearch := range day.TrendingSearches {
			article := ArticleInfo{
				Title: trendingSearch.Articles[0].Title,
				URL:   trendingSearch.Articles[0].URL,
			}

			consolidatedTrend := ConsolidatedTrend{
				TrendingSearchTitle: trendingSearch.Title.Query,
				Article:             article,
			}
			consolidatedTrends = append(consolidatedTrends, consolidatedTrend)
		}
	}

	return consolidatedTrends
}

func fetchArticleContent(url string, headline string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return headline
	}

	req.Header.Set("User-Agent", "")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return headline
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return headline
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
		return headline
	}

	articleText := ""
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		articleText += s.Text() + " "
	})

	return articleText
}