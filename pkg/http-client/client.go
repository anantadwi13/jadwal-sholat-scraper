package http_client

import (
	"io"
	"net/http"
	_url "net/url"
	"strings"
	"sync"
)

type MapCookies map[string]*http.Cookie

type FormData map[string]string

type HttpClient struct {
	client  http.Client
	cookies MapCookies
	m       sync.Mutex
}

var instance *HttpClient
var once sync.Once

func HttpClientInstance() *HttpClient {
	once.Do(func() {
		instance = &HttpClient{}
		instance.client = http.Client{}
		instance.cookies = make(map[string]*http.Cookie)
	})
	return instance
}

func NewHttpClientInstance() *HttpClient {
	client := &HttpClient{}
	client.client = http.Client{}

	globalClient := HttpClientInstance()

	if globalClient != nil {
		client.cookies.Append(&globalClient.cookies)
	} else {
		client.cookies = make(MapCookies)
	}

	return client
}

func (m *MapCookies) Append(cookies *MapCookies) {
	if *m == nil {
		*m = make(MapCookies)
	}

	for key, value := range *cookies {
		(*m)[key] = value
	}
}

func (h *HttpClient) Client() *http.Client {
	return &h.client
}

func (h *HttpClient) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)

	res, err := h.sendRequest(req)

	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (h *HttpClient) Post(url string, formData FormData) (string, error) {
	_formData := _url.Values{}

	for key, value := range formData {
		_formData.Add(key, value)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(_formData.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := h.sendRequest(req)

	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (h *HttpClient) sendRequest(req *http.Request) (*http.Response, error) {
	for _, cookie := range h.cookies {
		req.AddCookie(cookie)
	}

	res, err := h.Client().Do(req)

	if err != nil {
		return nil, err
	}

	h.m.Lock()
	defer h.m.Unlock()

	for _, cookie := range res.Cookies() {
		h.cookies[cookie.Name] = cookie
	}

	return res, err
}
