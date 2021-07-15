package main

import (
	"testing"
)

func TestClientCreation(t *testing.T) {
	ctr, stargazers_url := clientCreation("Personal-Projects")

	if ctr != 29 {
		t.Error("Wrong {} Number of Repositories should be 29", ctr)
	}
	checkUrl := "https://api.github.com/repos/rohitchauraisa1997/Personal-Projects/stargazers"
	if stargazers_url != checkUrl {
		t.Error("Wrong {} stargazers url retrieved for rohitchauraisa1997 should be {}", stargazers_url, checkUrl)
	}
}

func TestStargazersInfo(t *testing.T) {
	url := "https://api.github.com/repos/rohitchauraisa1997/Personal-Projects/stargazers"
	stargazers, numberOfStargazers := getStargazersInfo(url)
	if stargazers[0] != "https://api.github.com/users/shagunchaurasia/followers" {
		t.Error("retrieved wrong stargazers follower url")
	}
	if numberOfStargazers != 3 {
		t.Error("retrieved wrong number of stargazers")
	}

}

func TestExtractUsers(t *testing.T) {
	stargazers := []string{"https://api.github.com/users/shagunchaurasia/followers"}
	extractedusers := extractUsers(stargazers)
	if extractedusers[0] != "shagunchaurasia" {
		t.Error("retrieved wrong user from url")
	}
}