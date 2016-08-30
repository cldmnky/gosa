package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SaltClient struct {
	Hostname string
	Port     string
	Client   http.Client
	Token    string
}

type LoginResp struct {
	Return []struct {
		Perms  []interface{} `json:"perms"`
		Start  float64       `json:"start"`
		Token  string        `json:"token"`
		Expire float64       `json:"expire"`
		User   string        `json:"user"`
		Eauth  string        `json:"eauth"`
	} `json:"return"`
}

// Request represents a single request made to the Salt API
type Request struct {
	Client         string `url:"client"`
	Target         string `url:"tgt"`
	Function       string `url:"fun"`
	Arguments      string `url:"arg,omitempty"`
	ExpressionForm string `url:"expr_form,omitempty"`
}

func NewSaltClient(hostname string, port string) *SaltClient {
	s := new(SaltClient)
	s.Hostname = hostname
	s.Port = port
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s.Client = http.Client{
		Transport: tr,
	}
	return s
}

func (s *SaltClient) Login(username string, password string, eauth string) (string, error) {
	// logon and return token
	creds := url.Values{
		"username": {username},
		"password": {password},
		"eauth":    {eauth},
	}
	resp, err := s.Client.PostForm(s.Hostname+":"+s.Port+"/login", creds)
	if err != nil {
		return "", errors.New(fmt.Sprintf("gosa: PostForm error %v", err))
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("gosa: HTTP error %v", resp.StatusCode))
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	token, err := getTokenFromLogin([]byte(body))
	if err != nil {
		return "", errors.New(fmt.Sprintf("gosa: Token error %v", err))
	}
	return token, nil
}

func getTokenFromLogin(body []byte) (string, error) {
	var r = new(LoginResp)
	err := json.Unmarshal(body, &r)
	if err != nil {
		return "", errors.New(fmt.Sprintf("gosa: Unmarshal error %v", err))
	}
	t := r.Return[0].Token
	return t, nil
}

func (s *SaltClient) SetToken(token string) *SaltClient {
	// use a token already aquired and return a client
	s.Token = token
	return s
}

func (s *SaltClient) Run(target string, function string, arguments string) (string, error) {
	fmt.Println("Run")
	payload := Request{
		Client:    "local",
		Target:    target,
		Function:  function,
		Arguments: arguments,
	}
	values, _ := query.Values(payload)
	req, err := http.NewRequest("POST", s.Hostname+":"+s.Port, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return "", errors.New(fmt.Sprintf("gosa: Run error %v", err))
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("gosa: Run error %v", err))
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("gosa: HTTP error %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
