package services

import (
	"blog/models"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fpay/gopress"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	// ElasticServiceName is the identity of elastic service
	ElasticServiceName = "elastic"
	// Index elastic Index
	Index = "blog"
)

// ElasticService type
type ElasticService struct {
	EsClient *elastic.Client
}

// ElasticOption elastic options
type ElasticOption struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// NewElasticService returns instance of elastic service
func NewElasticService() *ElasticService {
	es := new(ElasticService)
	var err error
	es.EsClient, err = elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}

	return es
}

// ServiceName is used to implements gopress.Service
func (s *ElasticService) ServiceName() string {
	return ElasticServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *ElasticService) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

// SearchPosts Search
func (s *ElasticService) SearchPosts(keyword string, limit, page int) ([]*models.Post, error) {
	//s.EsClient.Index(Index).Type("posts").
	// Search with a term query
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	matchPhraseQuery := elastic.NewMatchPhraseQuery("Title", keyword)
	searchResult, err := s.EsClient.Search().
		Index(Index).
		Type("posts").                        // search in index
		Query(matchPhraseQuery).              // specify the query
		From((page - 1) * limit).Size(limit). // take documents 0-9
		Do(ctx)

	if err != nil {
		return nil, err
	}

	posts := []*models.Post{}

	if searchResult.Hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			p := &models.Post{}
			err := json.Unmarshal(*hit.Source, p)
			if err != nil {
				return nil, err
			}
			posts = append(posts, p)
		}
	} else {
		// No hits
		return nil, errors.New(" not found")
	}

	return posts, nil
}
