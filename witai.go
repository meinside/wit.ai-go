// https://wit.ai/docs/http/20160330

package witai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultVersion = "20141022"
)

type Client struct {
	Token   *string
	Version *string

	headerAuth   *string
	headerAccept *string
}

// https://wit.ai/docs/http/20160330#response-format-link
type ResponseError struct {
	Error *string `json:"error,omitempty"`
	Code  *string `json:"code,omitempty"`
}

// https://wit.ai/docs/http/20160330#converse-link
type ResponseConverse struct {
	ResponseError

	Type       *string     `json:"type"`
	Message    *string     `json:"msg,omitempty"`
	Action     *string     `json:"action,omitempty"`
	Entities   interface{} `json:"entities,omitempty"`
	Confidence float32     `json:"confidence"`
}

func (r ResponseConverse) String() string {
	return fmt.Sprintf("{Type: %s, Message: %s, Action: %s, Entities: %v, Confidence: %.6f}", *r.Type, *r.Message, *r.Action, r.Entities, r.Confidence)
}

// https://wit.ai/docs/http/20160330#context-link
type Context struct {
	State         interface{} `json:"state,omitempty"`
	ReferenceTime *string     `json:"reference_time,omitempty"`
	TimeZone      *string     `json:"timezone,omitempty"`
	Entities      *Entities   `json:"entities,omitempty"`
	Location      *Location   `json:"location,omitempty"`
}

func (c Context) String() string {
	return fmt.Sprintf("{State: %v, ReferenceTime: %s, TimeZone: %s, Entities: %v, Location: %v}", c.State, *c.ReferenceTime, *c.TimeZone, c.Entities, c.Location)
}

type Entities struct {
	Id     *string       `json:"id"`
	Doc    *string       `json:"doc,omitempty"`
	Values []interface{} `json:"values,omitempty"`
}

func (e Entities) String() string {
	return fmt.Sprintf("{Id: %s, Doc: %s, Values: %v}", *e.Id, *e.Doc, e.Values)
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (l Location) String() string {
	return fmt.Sprintf("{Latitude: %.6f, Longitude: %.6f}", l.Latitude, l.Longitude)
}

// https://wit.ai/docs/http/20160330#get-intent-via-text-link
type ResponseMessage struct {
	ResponseError

	MessageId *string   `json:"msg_id"`
	Text      *string   `json:"_text"`
	Outcomes  []Outcome `json:"outcomes"`
}

func (r ResponseMessage) String() string {
	return fmt.Sprintf("{MessageId: %s, Text: %s, Outcomes: %v}", *r.MessageId, *r.Text, r.Outcomes)
}

type Outcome struct {
	Text       *string     `json:"_text"`
	Intent     *string     `json:"intent"`
	Entities   interface{} `json:"entities"`
	Confidence int         `json:"confidence"`
}

func (o Outcome) String() string {
	return fmt.Sprintf("{Text: %s, Intent: %s, Entities: %v, Confidence: %d}", *o.Text, *o.Intent, o.Entities, o.Confidence)
}

// new client with default version
func NewClient(token *string) *Client {
	version := DefaultVersion
	return NewClientWithVersion(token, &version)
}

// new client with other version
func NewClientWithVersion(token, version *string) *Client {
	headerAuth := fmt.Sprintf("Bearer %s", *token)
	headerAccept := fmt.Sprintf("application/vnd.wit.%s+json", *version)

	return &Client{
		Token:        token,
		Version:      version,
		headerAuth:   &headerAuth,
		headerAccept: &headerAccept,
	}
}

// send http request with given method, url, and body data
func (c *Client) request(method, url string, body interface{}) (res []byte, err error) {
	var data []byte
	if data, err = json.Marshal(body); err == nil {
		var req *http.Request
		if req, err = http.NewRequest(method, url, bytes.NewBuffer(data)); err == nil {
			// headers
			req.Header.Set("Authorization", *c.headerAuth)
			req.Header.Set("Accept", *c.headerAccept)
			req.Header.Set("Content-Type", "application/json")

			var resp *http.Response
			client := &http.Client{}
			if resp, err = client.Do(req); err == nil {
				defer resp.Body.Close()

				res, _ = ioutil.ReadAll(resp.Body)
			} else {
				log.Printf("Error while sending request: %s\n", err.Error())
			}
		} else {
			log.Printf("Error while building request: %s\n", err.Error())
		}
	} else {
		log.Printf("Error while building request body: %s\n", err.Error())
	}

	return res, err
}

// upload voice file
func (c *Client) upload(method, url, filepath, contentType string) (res []byte, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filepath); err == nil {
		var req *http.Request
		if req, err = http.NewRequest(method, url, bytes.NewBuffer(data)); err == nil {
			// headers
			req.Header.Set("Authorization", *c.headerAuth)
			req.Header.Set("Accept", *c.headerAccept)
			req.Header.Set("Content-Type", contentType)

			var resp *http.Response
			client := &http.Client{}
			if resp, err = client.Do(req); err == nil {
				defer resp.Body.Close()

				res, _ = ioutil.ReadAll(resp.Body)
			} else {
				log.Printf("Error while sending request: %s\n", err.Error())
			}
		} else {
			log.Printf("Error while building request: %s\n", err.Error())
		}
	}

	return res, err
}

// make request url with given base url and GET parameters
func (c *Client) makeUrl(baseUrl string, params map[string]interface{}) *string {
	index := 0
	queries := make([]string, len(params))
	for k, v := range params {
		queries[index] = fmt.Sprintf("%s=%s", k, url.QueryEscape(fmt.Sprintf("%v", v)))
		index++
	}

	result := baseUrl + "?" + strings.Join(queries, "&")
	return &result
}

// get next steps
//
// https://wit.ai/docs/http/20160330#converse-link
func (c *Client) Converse(sessionId, query string, context interface{}) (response ResponseConverse, err error) {
	params := map[string]interface{}{
		"session_id": sessionId,
	}
	if context != nil {
		params["context"] = context
	}
	if len(query) > 0 {
		params["q"] = query
	}

	url := c.makeUrl("https://api.wit.ai/converse", params)

	var bytes []byte
	if bytes, err = c.request("POST", *url, context); err == nil {
		var converseRes ResponseConverse
		if err = json.Unmarshal(bytes, &converseRes); err == nil {
			if converseRes.Error == nil {
				response = converseRes
			} else {
				err = fmt.Errorf("converse request error: %s", *converseRes.Error)
			}
		} else {
			err = fmt.Errorf("converse parse error: %s", err)
		}
	}

	return response, err
}

func (c *Client) ConverseNext(sessionId string, context interface{}) (response ResponseConverse, err error) {
	return c.Converse(sessionId, "", context)
}

// get meaning of a sentence
//
// https://wit.ai/docs/http/20160330#get-intent-via-text-link
func (c *Client) Message(query string, context interface{}, messageId, threadId string, n int) (response ResponseMessage, err error) {
	params := map[string]interface{}{
		"q": query,
	}
	if context != nil {
		params["context"] = context
	}
	if len(messageId) > 0 {
		params["msg_id"] = messageId
	}
	if len(threadId) > 0 {
		params["thread_id"] = threadId
	}
	if n <= 0 {
		n = 1
	}
	params["n"] = n

	url := c.makeUrl("https://api.wit.ai/message", params)

	var bytes []byte
	if bytes, err = c.request("GET", *url, context); err == nil {
		var msgRes ResponseMessage
		if err = json.Unmarshal(bytes, &msgRes); err == nil {
			if msgRes.Error == nil {
				response = msgRes
			} else {
				err = fmt.Errorf("message request error: %s", *msgRes.Error)
			}
		} else {
			err = fmt.Errorf("message parse error: %s", err)
		}
	}

	return response, err
}

// get meaning of audio (mp3 format)
//
// https://wit.ai/docs/http/20160330#get-intent-via-speech-link
func (c *Client) SpeechMp3(filepath string, context interface{}, messageId, threadId string, n int) (response ResponseMessage, err error) {
	params := map[string]interface{}{}
	if context != nil {
		params["context"] = context
	}
	if len(messageId) > 0 {
		params["msg_id"] = messageId
	}
	if len(threadId) > 0 {
		params["thread_id"] = threadId
	}
	if n <= 0 {
		n = 1
	}
	params["n"] = n

	url := c.makeUrl("https://api.wit.ai/speech", params)

	var bytes []byte
	if bytes, err = c.upload("POST", *url, filepath, "audio/mpeg3"); err == nil {
		var speechRes ResponseMessage
		if err = json.Unmarshal(bytes, &speechRes); err == nil {
			if speechRes.Error == nil {
				response = speechRes
			} else {
				err = fmt.Errorf("speech request error: %s", *speechRes.Error)
			}
		} else {
			err = fmt.Errorf("speech parse error: %s", err)
		}
	}

	return response, err
}
