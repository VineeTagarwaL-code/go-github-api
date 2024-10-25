package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	Date  string `json:"date"`
	Count int    `json:"count"`
	Level int    `json:"level"`
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

func mapContributionLevel(level string) int {
	switch level {
	case "FIRST_QUARTILE":
		return 1
	case "SECOND_QUARTILE":
		return 2
	case "THIRD_QUARTILE":
		return 3
	case "FOURTH_QUARTILE":
		return 4
	default:
		return 0 // Use 0 for "NONE" or unknown levels
	}
}

func GithubHandler(r *gin.Context) {
	name := r.Param("name")

	if name == "" {
		r.JSON(http.StatusBadRequest, gin.H{
			"error": "Year parameter is required",
		})
		return
	}

	currentTime := time.Now()

	// Set the 'from' date to the same date last year
	from := time.Date(currentTime.Year()-1, currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	// Set the 'to' date to the present date and time
	to := currentTime

	// Format the dates in the desired format
	fromStr := from.Format("2006-01-02T15:04:05Z")
	toStr := to.Format("2006-01-02T15:04:05Z")

	fmt.Printf("From: %s\n", fromStr)
	fmt.Printf("To: %s\n", toStr)

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
				Date:  day.Date,
				Count: day.ContributionCount,
				Level: mapContributionLevel(day.ContributionLevel),
			})
		}
	}

	defer resp.Body.Close()
	r.JSON(http.StatusOK, gin.H{
		"data": apiResponse,
	})
}
