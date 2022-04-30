package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Session struct {
	token string
}

var session Session

type Project struct {
	FullName string `json:"full_name"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("input token")
		return
	}
	unstar(os.Args[1])
}

func unstar(token string) {
	session = Session{token: token}
	repos, err := session.getAllRepos()
	if err != nil {
		fmt.Println(err)
	}
	session.unStar(repos)
}

func (s Session) unStar(projects []Project) {
	for index, p := range projects {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://api.github.com/user/starred/%s", p.FullName), nil)
		req.Header.Add("Authorization", fmt.Sprintf("token %s", s.token))
		if err != nil {
			log.Printf("[%d]unstar %s %v", index, p.FullName, err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[%d]unstar %s %v", index, p.FullName, err)
		}
		if res.StatusCode < 300 {
			log.Printf("[%d]unstar %-30s success", index, p.FullName)
		}
	}
}

func (s Session) getAllRepos() ([]Project, error) {
	projects := []Project{}
	page := 1
	pageSize := 100
	for {
		t, err := s.getRepos(page, pageSize)
		if err != nil {
			log.Println(err)
			return projects, nil
		}
		if len(t) == 0 {
			break
		}
		page += 1
		projects = append(projects, t...)
	}
	return projects, nil
}

func (s Session) getRepos(page int, pageSize int) ([]Project, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/user/starred?per_page=%d&page=%d", pageSize, page), nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", s.token))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		log.Printf("%s", bs)
		return nil, err
	}
	var projects = []Project{}
	err = json.Unmarshal(bs, &projects)
	return projects, err
}
