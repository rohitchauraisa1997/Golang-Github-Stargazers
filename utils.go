package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type Response struct {
	StarCount int      `json:"star_count"`
	Followers []string `json:"followers"`
}

type StarGazerInfo struct {
	AvatarUrl           string `json:"avatar_url"`
	EventtsUrl          string `json:"events_url"`
	FollowersUrl        string `json:"followers_url"`
	FollowingsUrl       string `json:"following_url"`
	GistsUrl            string `json:"gists_url"`
	GravatarId          string `json:"gravatar_id"`
	HtmlUrl             string `json:"html_url"`
	Id                  int    `json:"id"`
	Login               string `json:"login"`
	NodeID              string `json:"node_id"`
	OrganizationUrl     string `json:"organizations_url"`
	Received_Events_Url string `json:"received_events_url"`
	ReposUrl            string `json:"repos_url"`
	SiteAdmin           bool   `json:"site_admin"`
	StarredUrl          string `json:"starred_url"`
	Subscriptions_Url   string `json:"subscriptions_url"`
	Type                string `json:"type"`
	Url                 string `json:"url"`
}

func clientCreation(reponame string) (int, string) {
	/*
		Creates client and returns the number of total repos alongwith
		the stargazers[people who have starred the r	if repo_count, gazers_url == 0, ""{

		}equired repo] url.
	*/
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	// log.Println(repos)
	if err != nil {
		log.Print("ERROR ", err)
		return 0, ""
	}
	stargazers_url := ""
	ctr := 0
	for _, repo := range repos {
		// parses json out for correct repo
		if *repo.Name == reponame {
			// log.Println(repo)
			stargazers_url = *repo.StargazersURL
			log.Println("stargazers_url", stargazers_url)
		}
		ctr += 1
	}
	log.Println(ctr)
	return ctr, stargazers_url
}

func getStargazersInfo(url string) ([]string, int) {
	/*
		sends a GET request to stargazers_url
		returns list of followers and number of followers.
	*/
	if url == "" {
		log.Print("No Url Resolved please crosscheck the Repository name")
		return []string{}, 0
	}
	var stargazers []string

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	q := req.URL.Query()
	req.Header.Set("Content-Type", "application/json")

	req.URL.RawQuery = q.Encode()
	log.Println("url to hit.", req.URL.String())

	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var star_ifno []StarGazerInfo
	jsonErr := json.Unmarshal(bodyText, &star_ifno)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, value := range star_ifno {
		// log.Println(_, value)
		stargazers = append(stargazers, value.FollowersUrl)
	}
	log.Println("list of stargazers with url", stargazers)
	return stargazers, len(stargazers)
}

func extractUsers(stargazers []string) []string {
	/*
		extracts users from the complete url
		""
	*/
	var extractedusers []string
	for index, user := range stargazers {
		log.Println("extracting users", index, user)
		splitUsername := strings.Split(user, "/")
		log.Println(splitUsername[4])
		// log.Printf("\n%T\n", user)
		extractedusers = append(extractedusers, splitUsername[4])
	}
	log.Println("Extracted Users Who Starred the repo", extractedusers)
	return extractedusers
}

func callFlow(checkRepoName string) map[string]interface{} {
	/*
		callFlow function.
	*/
	repo_count, gazers_url := clientCreation(checkRepoName)
	log.Println("repo_count", repo_count)
	log.Println("gazers_url", gazers_url)
	if repo_count == 0 && gazers_url == "" {
		answerDict := map[string]interface{}{"star_count": nil, "followers": []string{"Wrong Github Token"}}
		return answerDict
	} else {
		followersList, star_count := getStargazersInfo(gazers_url)
		if reflect.DeepEqual(followersList, []string{}) && star_count == 0 {
			answerDict := map[string]interface{}{"star_count": nil, "followers": []string{"Wrong Repo Name"}}
			return answerDict
		} else {
			log.Println(followersList)
			log.Println(star_count)

			log.Println("Exttracting Users")
			followersList = extractUsers(followersList)

			answer := Response{StarCount: star_count, Followers: followersList}
			log.Println(answer)

			answerDict := map[string]interface{}{"star_count": star_count, "followers": followersList}
			log.Println(answerDict)
			return answerDict
		}
	}
}
