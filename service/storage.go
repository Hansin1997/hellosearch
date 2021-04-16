package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Storage interface {
	Search(query string, indexes []string, page int) (*SearchResults, error)
	List() ([]string, error)
}

type ElasticStorage struct {
	client  *elasticsearch.Client
	context context.Context
}

func NewElasticStorage(address []string,
	username string,
	password string,
	context context.Context) (*ElasticStorage, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: address,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return &ElasticStorage{client: es, context: context}, nil
}

func (e ElasticStorage) Search(query string, indexes []string, page int) (*SearchResults, error) {
	if page <= 0 {
		page = 1
	}
	res, err := e.client.Search(
		e.client.Search.WithIndex(indexes...),
		e.client.Search.WithBody(buildQuery(query, (page-1)*10)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	var r envelopeResponse

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	var results SearchResults
	results.Total = r.Hits.Total.Value

	if len(r.Hits.Hits) < 1 {
		results.Hits = []*Hit{}
		return &results, nil
	}
	for _, hit := range r.Hits.Hits {
		var h Hit
		h.Id = hit.ID
		h.Sort = hit.Sort

		if err := json.Unmarshal(hit.Source, &h); err != nil {
			return &results, err
		}

		if len(hit.Highlights) > 0 {
			if err := json.Unmarshal(hit.Highlights, &h.Highlights); err != nil {
				return &results, err
			}
		}

		results.Hits = append(results.Hits, &h)
	}

	return &results, nil
}

func (e ElasticStorage) List() ([]string, error) {
	resp, err := e.client.Cat.Indices(e.client.Cat.Indices.WithIndex("spider-*"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(bytes), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		items := strings.Split(line, " ")
		if len(items) >= 2 {
			result = append(result, strings.TrimSpace(items[2]))
		}
	}
	return result, nil
}

func buildQuery(query string, from int) io.Reader {
	var b strings.Builder

	b.WriteString("{\n")

	b.WriteString(fmt.Sprintf(searchMatch, query, from))

	b.WriteString("\n}")

	//fmt.Printf("%s\n", b.String())
	return strings.NewReader(b.String())
}

type envelopeResponse struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []struct {
			ID         string          `json:"_id"`
			Source     json.RawMessage `json:"_source"`
			Highlights json.RawMessage `json:"highlight"`
			Sort       []interface{}   `json:"sort"`
		}
	}
} // SearchResults wraps the Elasticsearch search response.
//
type SearchResults struct {
	Total int    `json:"total"`
	Hits  []*Hit `json:"hits"`
}

// Hit wraps the document returned in search response.
//
type Hit struct {
	Document
	Sort       []interface{} `json:"sort"`
	Highlights *struct {
		Title []string `json:"Title"`
		Alt   []string `json:"Content"`
	} `json:"highlights,omitempty"`
}

type Document struct {
	Id             string
	Url            string
	Title          string
	ResponseHeader http.Header
	FetchAt        time.Time
}

const searchMatch = `	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["Title", "Content^10","Keywords","Description"],
			"operator" : "and"
		}
	},
	"_source": ["Url","FetchAt","Title","ResponseHeader"]
	,
	"highlight" : {
		"fields" : {
			"Title" : { "number_of_fragments" : 0 },
			"Content" : { }
		}
	},
	"size" : 10,
	"from": %d,
	"sort" : [ { "_score" : "desc" }, { "_doc" : "asc" } ]`
