package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Session struct {
	token string
}

var session = Session{}

type Project struct {
	FullName string `json:"full_name"`
}

func main() {
	fmt.Println("hello,world!")
}

func (s Session) unStar(projects []Project) {
	for _, p := range projects {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://api.github.com/user/starred/%s", p.FullName), nil)
		req.Header.Add("Authorization", fmt.Sprintf("token %s", s.token))
		if err != nil {
			log.Println(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		}
		if res.StatusCode < 300 {
			log.Printf("unstar %s success", p.FullName)
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

/*

curl \
  -H "Authorization: token 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"\
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/user/starred


curl \
  -X DELETE \
  -H "Authorization: token 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"\
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/user/starred/airbnb/HorizonCalendar

*/
