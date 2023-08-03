package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


func checkFile(filename string) error {
    _, err := os.Stat(filename)
        if os.IsNotExist(err) {
            _, err := os.Create(filename)
                if err != nil {
                    return err
                }
        }
        return nil
}

type NaverQADB struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Question string `json:"question"`
	Answers []string `json:"answer"`
	Topic string `json:"topic"`
}


func main() {
	filename := "webdata.json"
	err := checkFile(filename)
	if (err != nil){
		log.Fatal(err)
	}
	
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://kin.naver.com/qna/detail.naver?d1id=6&dirId=60205&docId=450471426&qb=7J246rO17KeA64ql&enc=utf8&section=kin&rank=1&search_sort=0&spq=0";

	resp, err := http.Get(url);
	if err != nil {
		log.Fatal(err);
	}
	defer resp.Body.Close();

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
	 log.Fatal(err)
	}

	data := []NaverQADB{}
	json.Unmarshal(file, &data)

	title := html.Find("div.c-heading__title-inner").Find("div.title").Text()
	title = strings.TrimSpace(title)

	topic := html.Find("a.tag-list__item.tag-list__item--category").Text()
	topic = strings.Join(strings.Split(topic, " ")[2:], " ")

	var temp_question []string
	question := ""
	
	_ = html.Find("div.c-heading__content").Find("div").Each(func(idx int, sel *goquery.Selection){
		t := sel.Text()
		t = strings.TrimSpace(t)

		temp_question = append(temp_question, t)
	})

	question = strings.Join(temp_question, "\n")
	fmt.Println(question)

	var answers []string
	
	_ = html.Find("div._answerListArea").Find("div.answer-content__list._answerList").Each(func(idx int, s *goquery.Selection){
		var temp_answer []string
		answer := ""

		s.Find(".se-module.se-module-text").Find("p").Each(func(idx int, sel *goquery.Selection) {
			t := sel.Text()
			t = strings.TrimSpace(t)

			temp_answer = append(temp_answer, t)
			fmt.Println(t)
		})

		answer = strings.Join(temp_answer, "\n")

		if (answer != ""){
			answers = append(answers, answer)
		}
	})
	
	newRow := &NaverQADB{
		ID: "1234",
		Title: title,
		Question: question,
		Answers: answers,
		Topic: topic,
	}

	data = append(data, *newRow)

	dataBytes, err := json.Marshal(data)
    if err != nil {
        log.Fatal(err)
    }

    err = ioutil.WriteFile(filename, dataBytes, 0644)
    if err != nil {
        log.Fatal(err)
    }
}