package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/olivere/elastic/v7"
	"log"
	"my-gin/libraries/config"
	zap "my-gin/libraries/log"
	"strings"
)

/**
 * 官方推荐
 */
func InitEs() (esClient *elasticsearch7.Client) {
	var (
		host = config.UnmarshalConfig.Elastic.Host
		err  error
		r    map[string]interface{}
	)

	cfg := elasticsearch7.Config{
		Addresses: []string{
			host,
		},
	}
	if esClient, err = elasticsearch7.NewClient(cfg); err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", err.Error())
		panic(err)
	}

	res, err := esClient.Info()
	if err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", fmt.Sprintf("Error getting response: %s", err))
	}
	// Check response status
	if res.IsError() {
		zap.InitLog("elastic").Errorf("init", "msg", fmt.Sprintf("Error: %s", res.String()))
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", fmt.Sprintf("Error parsing the response body: %s", err))
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch7.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
	return
}

func Init() *elastic.Client {
	var (
		host     = config.UnmarshalConfig.Elastic.Host
		err      error
		esClient *elastic.Client
	)
	esClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetHealthcheck(false), elastic.SetSniff(false))
	if err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", err.Error())
		panic(err)
	}
	info, code, err := esClient.Ping(host).Do(context.Background())
	if err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", err.Error())
		panic(err)
	}
	zap.InitLog("elastic").Infof("init", "msg", fmt.Sprintf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number))

	esversion, err := esClient.ElasticsearchVersion(host)
	if err != nil {
		zap.InitLog("elastic").Errorf("init", "msg", err.Error())
		panic(err)
	}
	zap.InitLog("elastic").Infof("init", "msg", fmt.Sprintf("Elasticsearch version %s\n", esversion))
	return esClient
}
