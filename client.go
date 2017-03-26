package fatsecret

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	fatSecretAPIURL = "http://platform.fatsecret.com/rest/server.api"
)

// Client is the top-level FatSecret client which is used to
// fetch data from the FatSecret API using Oauth1 authentication
type Client struct {
	consumerKey    string
	consumerSecret string
	apiURL         string
	escapedAPIURL  string
	randSrc        rand.Source
	signer         Signer
}

// NewClient creates and returns a new FatSecret client instance
func NewClient(consumerKey string, consumerSecret string) (*Client, error) {
	// validate the given key and secret
	if consumerKey == "" {
		return nil, errors.New("Invalid consumer key given")
	}
	if consumerSecret == "" {
		return nil, errors.New("Invalid consumer key given")
	}

	// return the new client
	return &Client{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		apiURL:         fatSecretAPIURL,
		escapedAPIURL:  url.QueryEscape(fatSecretAPIURL),
		randSrc:        rand.NewSource(time.Now().UnixNano()),
		signer:         NewHMACSigner(consumerSecret),
	}, nil
}

// InvokeAPI calls the FatSecret API and returns the response body.
// This lower-level function is used by all higher-level API functions (ie: FoodSearch)
func (c *Client) InvokeAPI(apiMethod string, params map[string]string) ([]byte, error) {
	// build the oauth api url
	apiURL, err := c.buildURL(apiMethod, params)
	if err != nil {
		return nil, err
	}

	// invoke the http api call
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	// read the response message body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// return the response message body
	return body, nil
}

// buildURL builds and returns the oauth API URL based on the given parameters
func (c *Client) buildURL(apiMethod string, params map[string]string) (string, error) {
	// get the oauth time parameters
	ts := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("%d", rand.New(c.randSrc).Int63())

	// build the base message
	m := map[string]string{
		"method":                 apiMethod,
		"oauth_consumer_key":     c.consumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        ts,
		"oauth_version":          "1.0",
		"format":                 "json",
	}

	// add the given parameters to the message
	for k, v := range params {
		m[k] = v
	}

	// create a sorted array of oauth name keys
	oauthNames := make([]string, len(m))
	i := 0
	for k := range m {
		oauthNames[i] = k
		i++
	}
	sort.Strings(oauthNames)

	// build the oauth base signature string
	sigQuery := ""
	for _, k := range oauthNames {
		sigQuery += fmt.Sprintf("&%s=%s", k, m[k])
	}
	sigQuery = sigEscape(sigQuery[1:])
	sigBase := fmt.Sprintf("GET&%s&%s", c.escapedAPIURL, sigQuery)

	// generate the oauth sha1 base64 signature (no token key)
	oauthSig := c.signer.Sign("", sigBase)
	/*
		// generate the oauth sha1 base64 signature
		mac := hmac.New(sha1.New, []byte(c.consumerSecret+"&"))
		mac.Write([]byte(sigBase))
		oauthSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	*/

	// add the oauth signature to the map
	m["oauth_signature"] = oauthSig
	oauthNames = append(oauthNames, "oauth_signature")

	// re-sort the keys after adding the signature
	sort.Strings(oauthNames)

	// build the api request url
	apiURL := fmt.Sprintf("%s?", c.apiURL)
	apiQuery := ""
	for _, k := range oauthNames {
		apiQuery += fmt.Sprintf("&%s=%s", k, sigEscape(m[k]))
	}
	apiURL += apiQuery[1:]

	// return the api url
	return apiURL, nil
}

// escape the given string using url-escape plus some extras
func sigEscape(s string) string {
	s = url.QueryEscape(s)
	s = strings.Replace(s, "%7E", "~", -1)
	s = strings.Replace(s, "+", "%20", -1)
	s = strings.Replace(s, "!", "%21", -1)
	s = strings.Replace(s, "\\", "%27", -1)
	s = strings.Replace(s, "(", "%28", -1)
	s = strings.Replace(s, ")", "%29", -1)
	s = strings.Replace(s, "*", "%2A", -1)
	return s
}
