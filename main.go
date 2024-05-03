package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func fetchGitHubUserEvents(username, token string) string {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public?per_page=2", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	return string(body)
}

func formatDate(dateStr string) string {
	layout := time.RFC3339
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}
	return t.Format("2006-01-02") // YYYY-MM-DD形式で出力
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	username := os.Getenv("USERNAME")
	token := os.Getenv("TOKEN")

	res := fetchGitHubUserEvents(username, token)
	events, err := ParseGitHubEvents(res)
	if err != nil {
		log.Fatal("Error parsing GitHub events:", err)
	}

	for _, event := range events {
		fmt.Printf("Event ID: %s, Type: %s, Created At: %s\n", event.ID, event.Type, formatDate(event.CreatedAt))

	}
}
