package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var APIKey string

const (
	GooglePlaceAPIURL = "https://maps.googleapis.com/maps/api/place"
)

type Place struct {
	*googleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}
type googleResponse struct {
	Results []*Place `json:"results"`
	Status  string   `json:"status"`
}
type googleGeometry struct {
	*googleLocation `json:"location"`
}
type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type googlePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}
type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":   p.Name,
		"icon":   p.Icon,
		"photos": p.Vicinity,
		"lat":    p.Lat,
		"lng":    p.Lng,
	}
}

func (q *Query) find(types string) (*googleResponse, error) {
	u := GooglePlaceAPIURL + "/nearbysearch/json"
	// to construct url get params
	vals := make(url.Values)
	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("type", types)
	vals.Set("key", APIKey)
	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		vals.Set("minprice", fmt.Sprintf("%d", int(r.From-1)))
		vals.Set("minprice", fmt.Sprintf("%d", int(r.To-1)))
	}
	res, err := http.Get(u + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var response googleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (q *Query) Run() []interface{} {
	rand.Seed(time.Now().UnixNano())
	var w sync.WaitGroup
	// var mux sync.Mutex
	places := make([]interface{}, len(q.Journey))
	for i, r := range q.Journey {
		w.Add(1)
		// to get all result quickly, use goroutine to parallelize
		go func(types string, i int) {
			// close this routine and notify to waitgroup
			defer w.Done()
			response, err := q.find(types)
			if err != nil {
				log.Println("failed to search facilities", err)
				return
			}
			if len(response.Results) == 0 {
				log.Println("cannot find any facility", err)
				log.Println("status: ", response.Status)
				return
			}
			for _, result := range response.Results {
				for _, photo := range result.Photos {
					photo.URL = GooglePlaceAPIURL + "/photo?" +
						"maxWidth=1000&photoreference=" + photo.PhotoRef + "&key=" + APIKey
				}
			}
			randT := rand.Intn(len(response.Results))
			// no need to lock because each goroutine write into another item of array
			// mux.Lock()
			places[i] = response.Results[randT]
			//mux.Unlock()
		}(r, i)
	}
	// wait for all request
	w.Wait()
	return places
}
