package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type NotionClient struct {
	Token string
}

type QueryFilter struct {
	Filter *Filter `json:"filter,omitempty"`
	Sorts  *[]Sort `json:"sorts,omitempty"`
}

type Filter struct {
	Property *string     `json:"property,omitempty"`
	Date     *DateFilter `json:"date,omitempty"`
}

type DateFilter struct {
	Before *string `json:"before,omitempty"`
}

type Sort struct {
	Property  *string `json:"property,omitempty"`
	Direction *string `json:"direction,omitempty"`
}

type Page struct {
	Object         *string `json:"object,omitempty"`
	Id             *string `json:"id,omitempty"`
	CreatedTime    *string `json:"createdTime,omitempty"`
	LastEditedTime *string `json:"lastEditedTime,omitempty"`
	Archived       *bool   `json:"archived,omitempty"`
	Icon           *FileObject
	Cover          *FileObject
	Properties     *map[string]Property `json:"properties,omitempty"`
	Url            *string              `json:"url,omitempty"`
}

type Property struct {
	Id      *string      `json:"id,omitempty"`
	Type    *string      `json:"type,omitempty"`
	Name    *string      `json:"name,omitempty"`
	Formula *FormulaProp `json:"formula,omitempty"`
	Title   *[]TitleProp `json:"title,omitempty"`
	Date    *DateProp    `json:"date,omitempty"`
}

type DateProp struct {
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
}

type TitleProp struct {
	Type      *string   `json:"type,omitempty"`
	Text      *TextProp `json:"text,omitempty"`
	PlainText *string   `json:"plain_text,omitempty"`
	Href      *string   `json:"href,omitempty"`
}

type TextProp struct {
	Content *string `json:"content,omitempty"`
	Link    *string `json:"link,omitempty"`
}

type FormulaProp struct {
	Type   *string `json:"type,omitempty`
	String *string `json:"string,omitempty`
	Number *int    `json:"number,omitempty`
}

type FileObject struct {
	Type *string         `json:"type,omitempty"`
	File *FileUrlWrapper `json:"file,omitempty"`
}

type FileUrlWrapper struct {
	Url        *string `json:"url,omitempty"`
	ExpiryTime *string `json:"expiry_time,omitempty"`
}

type QueryResult struct {
	Object  *string `json:"object,omitempty"`
	Results *[]Page `json:"results,omitempty"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
}

func (c *NotionClient) callApi(path string, method string, headers *map[string]string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("https://api.notion.com/v1%v", path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "(callApi) create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("Notion-Version", "2021-05-13")
	if headers != nil {
		for idx, el := range *headers {
			req.Header.Set(idx, el)
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(callApi) exec request")
	}
	defer res.Body.Close()

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(callApi) ioutil.readall")
	}
	return out, nil
}

func (c *NotionClient) QueryDatabase(databaseId string, queryFilter *QueryFilter) (*QueryResult, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	var qr QueryResult
	var body *bytes.Buffer

	if queryFilter != nil {
		jsonData, err := json.Marshal(queryFilter)
		if err != nil {
			return nil, errors.Wrap(err, "(QueryDatabase) marshal json")
		}
		body = bytes.NewBuffer(jsonData)
	}

	path := fmt.Sprintf("/databases/%v/query", databaseId)
	bytes, err := c.callApi(path, "POST", &headers, body)
	if err != nil {
		return nil, errors.Wrap(err, "(QueryDatabase) callApi")
	}

	err = json.Unmarshal(bytes, &qr)
	if err != nil {
		return nil, errors.Wrap(err, "(QueryDatabase) unmarshal response")
	}

	return &qr, nil
}

func (c *NotionClient) UpdatePage(pageId string, updates Page) error {
	fmt.Println("pageId", pageId)
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return errors.Wrap(err, "(FetchPage) marshal json")
	}

	client := &http.Client{}
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%v", pageId)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "(FetchPage) create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("Notion-Version", "2021-05-13")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "(FetchPage) exec request")
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		defer res.Body.Close()
		out, _ := ioutil.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Failed with status %v, body: %v", res.StatusCode, string(out)))
	}

	return nil
}

func (c *NotionClient) FetchPage(pageId string) (*Page, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%v", pageId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "(FetchPage) create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("Notion-Version", "2021-05-13")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(FetchPage) exec request")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(FetchPage) read body")
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, errors.Wrap(err, "(FetchPage) unmarshal json")
	}

	return &page, nil
}
