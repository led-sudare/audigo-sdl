package net

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	HTTPQUEUE_SIZE = 50
	ENDPOINT       = "audio/v1"
)

type Action struct {
	Act  string
	Args interface{}
}

type Requester struct {
	domain string
	id     string

	httpQueue  chan interface{}
	httpClient *http.Client
	// requestPool *audigoPool
}

// type audigoPool struct {
// 	playPool   *sync.Pool
// 	stopPool   *sync.Pool
// 	pausePool  *sync.Pool
// 	resumePool *sync.Pool
// 	volumePool *sync.Pool
// }

func NewRequester(domain string, id string) *Requester {
	r := &Requester{
		domain:     domain,
		id:         id,
		httpClient: &http.Client{},
		httpQueue:  make(chan interface{}, HTTPQUEUE_SIZE),
	}
	// r.initPool()

	go r.work()
	return r
}

func newRequest(url string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Errorf("error getting request: %s", err.Error())
	}
	return req
}

func (r *Requester) GetQueue() chan<- interface{} {
	return r.httpQueue
}

func (r *Requester) work() {
	for {
		select {
		case v := <-r.httpQueue:
			go r.request(v.(*Action))
		}
	}
}

func (r *Requester) request(v *Action) {
	// req := r.getRequest(v.Act)

	json, err := json.Marshal(v.Args)
	if err != nil {
		log.Errorf("dont marshal args: %s\n%s", v.Args, err.Error())
	}
	body := bytes.NewBuffer(json)
	uri := r.createUrl(v.Act)
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		log.Errorf("dont create new request: %s\n%s", v.Act, err.Error())
	}
	printRequest(req)
	res, _ := r.httpClient.Do(req)
	printResponse(res)
}

// func (r *Requester) getRequest(act string) *http.Request {
// 	var f interface{}
// 	switch strings.ToLower(act) {
// 	case "play":
// 		f = r.requestPool.playPool.Get()
// 	case "stop":
// 		f = r.requestPool.stopPool.Get()
// 	case "pause":
// 		f = r.requestPool.pausePool.Get()
// 	case "resume":
// 		f = r.requestPool.resumePool.Get()
// 	case "volume":
// 		f = r.requestPool.volumePool.Get()
// 	default:
// 		return nil
// 	}
// 	return f.(*http.Request)
// }

// func (r *Requester) initPool() {
// 	r.requestPool = &audigoPool{}
// 	r.requestPool.playPool = &sync.Pool{
// 		New: r.getCreationPoolFunc("play"),
// 	}
// 	r.requestPool.stopPool = &sync.Pool{
// 		New: r.getCreationPoolFunc("stop"),
// 	}
// 	r.requestPool.pausePool = &sync.Pool{
// 		New: r.getCreationPoolFunc("pause"),
// 	}
// 	r.requestPool.resumePool = &sync.Pool{
// 		New: r.getCreationPoolFunc("resume"),
// 	}
// 	r.requestPool.volumePool = &sync.Pool{
// 		New: r.getCreationPoolFunc("volume"),
// 	}
// }

// func (r *Requester) getCreationPoolFunc(act string) func() interface{} {
// 	return func() interface{} {
// 		url := createUrl(act)
// 		log.Debugf("*** %s", url)
// 		r := newRequest(url, nil)
// 		return r
// 	}
// }

func (r *Requester) createUrl(act string) string {
	uri, _ := url.Parse(r.domain)
	path := path.Join(ENDPOINT, act, r.id)
	uri, _ = uri.Parse(path)
	return uri.String()
}

func printRequest(r *http.Request) {
	if r == nil {
		log.Warn("not found request")
		return
	}

	log.Debug("show request.")
	log.Debugf("[method] %s", r.Method)
	log.Debugf("[uri] %s", r.URL.String())
	for k, v := range r.Header {
		log.Debugf("[header] %s: %s", k, strings.Join(v, ","))
	}
}

func printResponse(r *http.Response) {
	if r == nil {
		log.Warn("not found response")
		return
	}

	log.Debug("show response.")
	log.Debugf("[status] %d", r.StatusCode)
	for k, v := range r.Header {
		log.Debugf("[header] %s: %s", k, strings.Join(v, ","))
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Warn("dont open response body: %s", err.Error())
	}
	log.Debugf("[body] %s", string(body))
}
