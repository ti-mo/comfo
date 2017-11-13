package comfoclient

import (
	"encoding/json"
	"fmt"
	"github.com/ti-mo/comfo/libcomfo"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	userAgent = "comfoclient"
)

type Client struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func New(baseUrlStr string) Client {

	baseUrl, err := url.Parse(baseUrlStr)
	if err != nil {
		log.Fatalln("Invalid base URL for comfoclient", baseUrlStr)
	}

	return Client{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		baseUrl: baseUrl,
	}
}

// req prepares and executes an HTTP request against
// the endpoint configured in cl.baseUrl.
func (cl *Client) req(path string, method string, reqBody io.Reader) (respBody []byte, err error) {

	// Safely build URL
	url := cl.baseUrl.ResolveReference(&url.URL{Path: path})

	// Prepare the request
	hr, err := http.NewRequest(method, url.String(), reqBody)
	if err != nil {
		return
	}

	// Set headers
	hr.Header.Set("User-Agent", userAgent)

	// Execute the HTTP request
	resp, err := cl.httpClient.Do(hr)
	if err != nil {
		return
	}

	// Read the body
	respBody, err = ioutil.ReadAll(resp.Body)

	return
}

func (cl *Client) get(path string) (body []byte, err error) {

	body, err = cl.req(path, http.MethodGet, nil)
	return
}

func (cl *Client) set(path string) (body []byte, err error) {

	body, err = cl.req(path, http.MethodPut, nil)
	return
}

func (cl *Client) SetFans(speed string) (fans libcomfo.Fans, fp libcomfo.FanProfiles, err error) {

	body, err := cl.set(fmt.Sprintf("/fans/%v", speed))
	if err != nil {
		return
	}

	temp := struct {
		libcomfo.Fans        `json:"fans"`
		libcomfo.FanProfiles `json:"profiles"`
	}{}

	err = json.Unmarshal(body, &temp)

	return temp.Fans, temp.FanProfiles, err
}

func (cl *Client) GetTemps() (temps libcomfo.Temps, err error) {

	body, err := cl.get("/temps")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &temps)

	return
}

func (cl *Client) GetFans() (fans libcomfo.Fans, err error) {

	body, err := cl.get("/fans")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &fans)

	return
}

func (cl *Client) GetFanProfiles() (fp libcomfo.FanProfiles, err error) {

	body, err := cl.get("/profiles")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &fp)

	return
}
