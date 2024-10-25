package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const githubEndpoint = "https://api.github.com/graphql"

type GraphqlrRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type ContributionResponse struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
							ContributionLevel string `json:"contributionLevel"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

type ApiResponse struct {
	Day     string `json:"day"`
	Commits int    `json:"commits"`
	Level   string `json:"level"`
}

type ErrorResponse struct {
	Errors []struct {
		Extensions struct {
			Value    interface{} `json:"value"`
			Problems []struct {
				Path        []interface{} `json:"path"`
				Explanation string        `json:"explanation"`
			} `json:"problems"`
		} `json:"extensions"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Message string `json:"message"`
	} `json:"errors"`
}

func GithubHandler(r *gin.Context) {
	name := r.Param("name")
	year := r.Query("year")
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env")
	}
	if name == "" || year == "" {
		r.JSON(http.StatusBadRequest, gin.H{
			"error": "Year parameter is required",
		})
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		r.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid year parameter",
		})
		return
	}

	from := fmt.Sprintf("%d-01-01T00:00:00Z", yearInt)
	to := fmt.Sprintf("%d-12-31T23:59:59Z", yearInt)

	query := `query ($username: String!, $from: DateTime!, $to: DateTime!) {
		user(login: $username) {
			contributionsCollection(from: $from, to: $to) {
				contributionCalendar {
					weeks {
						contributionDays {
							date
							contributionCount
							contributionLevel
						}
					}
				}
			}
		}
	}`
	variables := map[string]interface{}{
		"username": name,
		"from":     from,
		"to":       to,
	}
	RequestBody := GraphqlrRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(RequestBody)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
	}
	req, err := http.NewRequest("POST", githubEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from github api"})
	}

	var contriResponse ContributionResponse
	err = json.Unmarshal([]byte(body), &contriResponse)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal the response from github api"})
	}
	var apiResponse []ApiResponse
	for _, week := range contriResponse.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			apiResponse = append(apiResponse, ApiResponse{
				Day:     day.Date,
				Commits: day.ContributionCount,
				Level:   day.ContributionLevel,
			})
		}
	}

	defer resp.Body.Close()
	r.JSON(http.StatusOK, gin.H{
		"data": apiResponse,
	})
}
