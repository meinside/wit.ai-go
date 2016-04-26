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
	DefaultVersion = "20160330" // last update: 2016.04.26.
)

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

	url := baseUrl
	if len(params) > 0 {
		url = url + "?" + strings.Join(queries, "&")
	}

	return &url
}

// get next steps
//
// https://wit.ai/docs/http/20160330#converse-link
func (c *Client) ConverseFirst(sessionId, query string, context interface{}) (response Converse, err error) {
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
		var converseRes Converse
		if err = json.Unmarshal(bytes, &converseRes); err == nil {
			if converseRes.Error == nil {
				response = converseRes
			} else {
				err = fmt.Errorf("converse response error: %s", *converseRes.Error)
			}
		} else {
			err = fmt.Errorf("converse parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("converse request error: %s", err)
	}

	return response, err
}

func (c *Client) ConverseNext(sessionId string, context interface{}) (response Converse, err error) {
	return c.ConverseFirst(sessionId, "", context)
}

// get meaning of a sentence
//
// https://wit.ai/docs/http/20160330#get-intent-via-text-link
func (c *Client) Message(query string, context interface{}, messageId, threadId string, n int) (response Message, err error) {
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
		var msgRes Message
		if err = json.Unmarshal(bytes, &msgRes); err == nil {
			if msgRes.Error == nil {
				response = msgRes
			} else {
				err = fmt.Errorf("message response error: %s", *msgRes.Error)
			}
		} else {
			err = fmt.Errorf("message parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("message request error: %s", err)
	}

	return response, err
}

// get meaning of audio (mp3 format)
//
// https://wit.ai/docs/http/20160330#get-intent-via-speech-link
func (c *Client) SpeechMp3(filepath string, context interface{}, messageId, threadId string, n int) (response Message, err error) {
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
		var speechRes Message
		if err = json.Unmarshal(bytes, &speechRes); err == nil {
			if speechRes.Error == nil {
				response = speechRes
			} else {
				err = fmt.Errorf("speech response error: %s", *speechRes.Error)
			}
		} else {
			err = fmt.Errorf("speech parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("speech request error: %s", err)
	}

	return response, err
}

// create new intents
//
// https://wit.ai/docs/http/20160330#intents-post-link
func (c *Client) CreateNewIntent(intents ...Intent) (response Intents, err error) {
	var data interface{}

	if len(intents) > 1 {
		data = intents
	} else {
		data = intents[0]
	}

	url := c.makeUrl("https://api.wit.ai/intents", nil)

	var bytes []byte
	if bytes, err = c.request("POST", *url, data); err == nil {
		var intentsRes Intents
		if err = json.Unmarshal(bytes, &intentsRes); err == nil {
			if intentsRes.Error == nil {
				response = intentsRes
			} else {
				err = fmt.Errorf("new intents response error: %s", *intentsRes.Error)
			}
		} else {
			err = fmt.Errorf("new intents parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("new intents request error: %s", err)
	}

	return response, err
}

// retrieve the list of all intents
//
// https://wit.ai/docs/http/20160330#intents-index-link
func (c *Client) GetAllIntents() (response []Intent, err error) {
	url := c.makeUrl("https://api.wit.ai/intents", nil)

	var bytes []byte
	if bytes, err = c.request("GET", *url, nil); err == nil {
		var intentsRes []Intent
		if err = json.Unmarshal(bytes, &intentsRes); err == nil {
			response = intentsRes
		} else {
			err = fmt.Errorf("intent list parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("intent list request error: %s", err)
	}

	return response, err

}

// retrieve all entities and expressions of an intent
//
// https://wit.ai/docs/http/20160330#intent-show-link
func (c *Client) ShowIntent(intentId *string) (response IntentDetail, err error) {
	// TODO
	return IntentDetail{}, nil
}

// update intent attributes
//
// https://wit.ai/docs/http/20160330#intent-put-link
func (c *Client) UpdateIntent(intentId *string) (response IntentAttributes, err error) {
	// TODO
	return IntentAttributes{}, nil
}

// add new expressions to an intent
//
// https://wit.ai/docs/http/20160330#create-intent-expressions-link
func (c *Client) CreateIntentExpressions(intentId *string, expressions []interface{}) (response []IntentExpression, err error) {
	// TODO
	return []IntentExpression{}, nil
}

// TODO

// remove an expression from an intent
//
// https://wit.ai/docs/http/20160330#destroy-intent-expression-link
func (c *Client) DeleteExpression(intentId *string, expressionId *string) (response interface{}, err error) {
	// TODO
	return nil, nil
}

// retrieve the list of all available entities
//
// https://wit.ai/docs/http/20160330#entities-index-link
func (c *Client) GetAllEntities() (response []string, err error) {
	// TODO
	return nil, nil
}

// create a new entity
//
// https://wit.ai/docs/http/20160330#entities-post-link
func (c *Client) CreateNewEntity(id *string, doc *string, values []interface{}) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// retrieve all values of an entity
//
// https://wit.ai/docs/http/20160330#entities-show-link
func (c *Client) ShowEntity(entityId *string) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// update the values of an entity
//
// https://wit.ai/docs/http/20160330#entities-put-link
func (c *Client) UpdateEntity(entityId *string, doc *string, values []interface{}) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// delete an entity
//
// https://wit.ai/docs/http/20160330#entities-destroy-link
func (c *Client) DeleteEntity(entityId *string) (response interface{}, err error) {
	// TODO
	return nil, nil
}

// add a new value to an entity
//
// https://wit.ai/docs/http/20160330#create-entity-value-link
func (c *Client) CreateEntityValue(entityId *string, value *string, expressions []string, metadata *string) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// remove a given value from an entity
//
// https://wit.ai/docs/http/20160330#delete-entity-value-link
func (c *Client) DeleteEntityValue(entityId *string, entityValue *string) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// create a new expression for an entity
//
// https://wit.ai/docs/http/20160330#create-entity-expression-link
func (c *Client) CreateEntityExpression(entityId *string, entityValue *string, expression *string) (response Entity, err error) {
	// TODO
	return Entity{}, nil
}

// remove an expression from an entity
//
// https://wit.ai/docs/http/20160330#destroy-entity-expression-link
func (c *Client) DeleteEntityExpression(entityId *string, entityValue *string, expressionValue *string) (response interface{}, err error) {
	// TODO
	return nil, nil
}

// retrieve an existing message
//
// https://wit.ai/docs/http/20160330#get-message-link
func (c *Client) GetMessage(messageId *string) (response Message, err error) {
	// TODO
	return Message{}, nil
}
