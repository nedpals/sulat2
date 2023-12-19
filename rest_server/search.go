package rest_server

import "net/http"

type SiteSearchResultDataType string

const (
	SiteSearchResultDataTypeCollection SiteSearchResultDataType = "collection"
	SiteSearchResultDataTypeRecord     SiteSearchResultDataType = "record"
)

type SiteSearchResultData struct {
	Id    string
	Label string
	Type  SiteSearchResultDataType
}

type SiteSearchResult struct {
	Total       int
	Collections []SiteSearchResultData
	Records     []SiteSearchResultData
}

func searchWithinSite(w http.ResponseWriter, r *http.Request) error {
	site := getCurrentSite(r)
	query := getDecodedQuery(r)
	results := SiteSearchResult{
		Collections: []SiteSearchResultData{},
		Records:     []SiteSearchResultData{},
	}

	if collections, err := site.Collections(); err == nil {
		for _, collection := range collections {
			results.Collections = append(results.Collections, SiteSearchResultData{
				Id:    collection.Id,
				Label: collection.Name,
				Type:  SiteSearchResultDataTypeCollection,
			})

			records, err := collection.Find(query, nil)
			if err != nil {
				continue
			}

			for _, record := range records {
				results.Records = append(results.Records, SiteSearchResultData{
					Id:    record.Id,
					Label: record.Title(),
					Type:  SiteSearchResultDataTypeRecord,
				})
			}
		}
	}

	results.Total = len(results.Collections) + len(results.Records)
	return returnJson(w, results)
}
