package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

var sessionTest = Session{token: token()}

func token() string {
	token, err := ioutil.ReadFile("token")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", token)
}
func Test_getRepos(t *testing.T) {
	got, err := sessionTest.getRepos(1, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	sessionTest.unStar(got)
}

func Test_getAllRepos(t *testing.T) {
	got, err := sessionTest.getAllRepos()
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range got {
		fmt.Printf("%d %s\n", i, v.FullName)
	}
	sessionTest.unStar(got)
}
