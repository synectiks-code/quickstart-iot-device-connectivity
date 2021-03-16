package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	geohash "github.com/mmcloughlin/geohash"
)

var CLOUDRACK_DATEFORMAT string = "20060102"

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Chunk(array []interface{}, chunkSize int) [][]interface{} {
	var divided [][]interface{}
	for i := 0; i < len(array); i += chunkSize {
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
		divided = append(divided, array[i:end])
	}
	return divided
}

//Buld a double array of integer from a max value and a chunkSize
func ChunkArray(size int64, chunkSize int64) [][]interface{} {
	toChunk := make([]interface{}, 0, 0)
	for i := int64(0); i < size; i++ {
		toChunk = append(toChunk, i)
	}
	return Chunk(toChunk, int(chunkSize))
}

//Generats a usinque ID based timestamp
func GeneratUniqueId() string {
	return strconv.FormatInt(int64(Hash(time.Now().Format("20060102150405"))), 10)
}

/**********
*Geo Hasshing
**************/
func GeoHash(lat, lng float64) string {
	return geohash.Encode(lat, lng)
}
func GeoHashWithPrecision(lat, lng, radius float64) string {
	prec := builCaracterLength(radius)
	return geohash.EncodeWithPrecision(lat, lng, prec)
}

//geohasg precision: https://en.wikipedia.org/wiki/Geohash
//1   ±2500
//2   ±630
//3   ±78
//4   ±20
//5   ±2.4
//6   ±0.61
//7   ±0.076
//8   ±0.019
func builCaracterLength(radius float64) uint {
	precisions := []float64{2500, 630, 78, 20, 2.4, 0.61, 0.076, 0.019}
	ind := 0
	if radius >= precisions[ind] {
		return uint(ind + 1)
	}
	ind = ind + 1
	for ind < len(precisions)-2 {
		if radius >= precisions[ind+1] && radius < precisions[ind] {
			return uint(ind + 1)
		}
		ind = ind + 1
	}
	log.Printf("[CORE][UTILS] Geohash precision for %f KM is  %v \n", radius, ind)
	return uint(ind)
}

func GeterateDateRange(startDate string, endDate string) []string {
	start, _ := time.Parse(CLOUDRACK_DATEFORMAT, startDate)
	end, _ := time.Parse(CLOUDRACK_DATEFORMAT, endDate)
	diff := int(end.Sub(start).Hours() / 24.0)
	log.Printf("[CORE][UTILS] Diff between %v and %v => %+v \n", start, end, diff)
	dates := make([]string, 0, 0)
	for i := 0; i < diff; i++ {
		t := start.AddDate(0, 0, i)
		dates = append(dates, t.Format(CLOUDRACK_DATEFORMAT))
	}
	log.Printf("[CORE][UTILS] Generated date range from %s to %s => %+v \n", startDate, endDate, dates)
	return dates
}

/****************
* HTTP WRAPPERS
******************/

type HttpService struct {
	Endpoint string
}

type RestOptions struct {
	SubEndpoint string
	Headers     map[string]string
}

func HttpInit(endpoint string) HttpService {
	log.Printf("[CORE][HTTP] initializing HTTP client with endpoint %v\n", endpoint)
	return HttpService{Endpoint: endpoint}
}

func processResponse(resp *http.Response) (*http.Response, error) {
	if resp.StatusCode >= 400 {
		return resp, errors.New(resp.Status)
	}
	return resp, nil
}

func (s HttpService) HttpPut(id string, object interface{}) (map[string]interface{}, error) {
	endpoint := s.Endpoint
	if id != "" {
		endpoint = endpoint + "/" + id
	}
	log.Printf("[CORE][HTTP] PUT %v\n", endpoint)
	log.Printf("[CORE][HTTP] Body: %+v\n", object)

	var resp *http.Response
	bytesRepresentation, err := json.Marshal(object)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(bytesRepresentation))
	resp, err = client.Do(req)
	log.Printf("[CORE][HTTP] PUT RESPONSE %v\n", resp)
	if err == nil {
		resp, err = processResponse(resp)

	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Printf("[CORE][HTTP] PUT RESPONSE DECODED %v\n", result)
	return result, err
}

func (s HttpService) HttpPost(object interface{}, tpl CloudrackObject, options RestOptions) (interface{}, error) {
	endpoint := s.Endpoint
	if options.SubEndpoint != "" {
		endpoint = endpoint + "/" + options.SubEndpoint
	}
	log.Printf("[CORE][HTTP] POST %v\n", endpoint)
	log.Printf("[CORE][HTTP] Body: %+v\n", object)

	var resp *http.Response
	bytesRepresentation, err := json.Marshal(object)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Content-Type", "application/json")
	for headerName, headerValue := range options.Headers {
		req.Header.Add(headerName, headerValue)
	}
	resp, err = client.Do(req)
	log.Printf("[CORE][HTTP] POST RESPONSE %v\n", resp)
	if err == nil && resp != nil {
		resp, err = processResponse(resp)
		res, _ := ParseBody(resp.Body, tpl)
		log.Printf("[CORE][HTTP] POST RESPONSE DECODED %v\n", res)
		return res, err
	}
	return tpl, err
}

func (s HttpService) HttpGet(params map[string]string, tpl CloudrackObject, options RestOptions) (interface{}, error) {
	endpoint := s.Endpoint
	if options.SubEndpoint != "" {
		endpoint = endpoint + "/" + options.SubEndpoint
	}
	log.Printf("[CORE][HTTP] GET %v\n", s.Endpoint)
	var resp *http.Response
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//Adding params
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
	//Workaround: uncode comma to support comma separated params
	req.URL.RawQuery = strings.ReplaceAll(req.URL.RawQuery, "%2C", ",")
	log.Printf("[CORE][HTTP] GET REQUEST %v\n", req)
	resp, err = client.Do(req)
	log.Printf("[CORE][HTTP] GET RESPONSE %v\n", resp)
	if err == nil {
		resp, err = processResponse(resp)
		//json.NewDecoder(resp.Body).Decode(&res)
		res, _ := ParseBody(resp.Body, tpl)
		log.Printf("[CORE][HTTP] GET RESPONSE DECODED %v\n", res)
		return res, err
	}
	return tpl, err
}

func ParseBody(body io.ReadCloser, p CloudrackObject) (interface{}, error) {
	log.Printf("[CORE][HTTP] parsing json body into %v", reflect.TypeOf(p))
	dec := json.NewDecoder(body)
	var err error = nil
	for {
		//if err = dec.Decode(&p); err == io.EOF {
		if err, p = p.Decode(*dec); err == io.EOF {
			break
		} else if err != nil {
			break
		}
	}
	if err == io.EOF {
		err = nil
	}
	return p, err
}

//type of object that is manageable by a micro service
type CloudrackObject interface {
	//ParseBody(body io.ReadCloser) (interface{}, error)
	Decode(dec json.Decoder) (error, CloudrackObject)
}
