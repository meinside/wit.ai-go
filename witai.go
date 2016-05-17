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
	DefaultVersion = "20160516" // last update: 2016.05.17.
)

// new client with default version
func NewClient(token string) *Client {
	version := DefaultVersion
	return NewClientWithVersion(token, version)
}

// new client with other version
func NewClientWithVersion(token, version string) *Client {
	headerAuth := fmt.Sprintf("Bearer %s", token)
	headerAccept := fmt.Sprintf("application/vnd.wit.%s+json", version)

	return &Client{
		Token:        &token,
		Version:      &version,
		headerAuth:   &headerAuth,
		headerAccept: &headerAccept,
	}
}

// send http request with given method, url, and body data
func (c *Client) request(method, url string, body interface{}) (res []byte, err error) {
	var data []byte
	if data, err = json.Marshal(body); err == nil {
		if c.Verbose {
			log.Printf("< HTTP request: %s %s, %s\n", method, url, string(data))
		}

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

				if c.Verbose {
					log.Printf("> HTTP response: %s\n", string(res))
				}
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
		if c.Verbose {
			log.Printf("< HTTP request: %s %s, %s (%s)\n", method, url, filepath, contentType)
		}

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

				if c.Verbose {
					log.Printf("> HTTP response: %s\n", string(res))
				}
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

// get meaning of a sentence
//
// https://wit.ai/docs/http/20160516#get--message-link
func (c *Client) QueryMessage(query string, context interface{}, messageId, threadId string) (response Message, err error) {
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

	url := c.makeUrl("https://api.wit.ai/message", params)

	var bytes []byte
	if bytes, err = c.request("GET", *url, context); err == nil {
		var msgRes Message
		if err = json.Unmarshal(bytes, &msgRes); err == nil {
			if !msgRes.HasError() {
				response = msgRes
			} else {
				err = fmt.Errorf("message response error: %s", msgRes.ErrorMessage())
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
// https://wit.ai/docs/http/20160516#post--speech-link
func (c *Client) QuerySpeechMp3(filepath string, context interface{}, messageId, threadId string, n int) (response Message, err error) {
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
			if !speechRes.HasError() {
				response = speechRes
			} else {
				err = fmt.Errorf("speech response error: %s", speechRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("speech parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("speech request error: %s", err)
	}

	return response, err
}

// get next steps
//
// https://wit.ai/docs/http/20160516#post--converse-link
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
			if !converseRes.HasError() {
				response = converseRes
			} else {
				err = fmt.Errorf("converse response error: %s", converseRes.ErrorMessage())
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

func (c *Client) ConverseAll(sessionId, query string, context interface{}) (responses []Converse, err error) {
	if result, err := c.ConverseFirst(sessionId, query, context); err == nil {
		responses = append(responses, result)

		for {
			if result, err := c.ConverseNext(sessionId, context); err == nil {
				responses = append(responses, result)

				if *result.Type != "stop" {
					continue
				}
			} else {
				return nil, err
			}
			break
		}

		return responses, nil
	} else {
		return nil, err
	}
}

// retrieve the list of all available entities
//
// https://wit.ai/docs/http/20160516#get--entities-link
func (c *Client) GetAllEntities() (response []string, err error) {
	url := c.makeUrl("https://api.wit.ai/entities", nil)

	var bytes []byte
	if bytes, err = c.request("GET", *url, nil); err == nil {
		var entitiesRes []string
		if err = json.Unmarshal(bytes, &entitiesRes); err == nil {
			response = entitiesRes
		} else {
			err = fmt.Errorf("get all entities parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("get all entities request error: %s", err)
	}

	return response, err
}

// create a new entity
//
// https://wit.ai/docs/http/20160516#post--entities-link
func (c *Client) CreateEntity(idOrName, doc *string, values ...EntityValue) (response Entity, err error) {
	url := c.makeUrl("https://api.wit.ai/entities", nil)

	data := map[string]interface{}{
		"id": *idOrName,
	}
	if doc != nil {
		data["doc"] = *doc
	}
	if len(values) > 0 {
		data["values"] = append([]EntityValue{}, values...)
	}

	var bytes []byte
	if bytes, err = c.request("POST", *url, data); err == nil {
		var entityRes Entity
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			if !entityRes.HasError() {
				response = entityRes
			} else {
				err = fmt.Errorf("new entity response error: %s", entityRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("new entity parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("new entity request error: %s", err)
	}

	return response, err
}

// retrieve all values of an entity
//
// https://wit.ai/docs/http/20160516#get--entities-:entity-id-link
func (c *Client) ShowEntity(entityId *string) (response Entity, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s", *entityId), nil)

	var bytes []byte
	if bytes, err = c.request("GET", *url, nil); err == nil {
		var entityRes Entity
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			if !entityRes.HasError() {
				response = entityRes
			} else {
				err = fmt.Errorf("show entity response error: %s", entityRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("show entity parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("show entity request error: %s", err)
	}

	return response, err
}

// update the values of an entity
//
// https://wit.ai/docs/http/20160516#put--entities-:entity-id-link
func (c *Client) UpdateEntity(entityId, doc *string, values ...EntityValue) (response Entity, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s", *entityId), nil)

	body := map[string]interface{}{}
	if doc != nil {
		body["doc"] = *doc
	}
	if len(values) > 0 {
		body["values"] = append([]EntityValue{}, values...)
	}

	var bytes []byte
	if bytes, err = c.request("PUT", *url, body); err == nil {
		var entityRes Entity
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			if !entityRes.HasError() {
				response = entityRes
			} else {
				err = fmt.Errorf("update entity response error: %s", entityRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("update entity parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("update entity request error: %s", err)
	}

	return response, err
}

// delete an entity
//
// https://wit.ai/docs/http/20160516#delete--entities-:entity-id-link
func (c *Client) DeleteEntity(entityId *string) (response map[string]string, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s", *entityId), nil)

	var bytes []byte
	if bytes, err = c.request("DELETE", *url, nil); err == nil {
		var entityRes map[string]string
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			response = entityRes
		} else {
			err = fmt.Errorf("delete entity parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("delete entity request error: %s", err)
	}

	return response, err
}

// add new values to an entity
//
// https://wit.ai/docs/http/20160516#post--entities-:entity-id-values-link
func (c *Client) CreateEntityValue(entityId, value *string, expressions []string, metadata *string) (response Entity, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s/values", *entityId), nil)

	body := map[string]interface{}{
		"value": *value,
	}
	if len(expressions) > 0 {
		body["expressions"] = expressions
	}
	if metadata != nil {
		body["metadata"] = *metadata
	}

	var bytes []byte
	if bytes, err = c.request("POST", *url, body); err == nil {
		var entityRes Entity
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			if !entityRes.HasError() {
				response = entityRes
			} else {
				err = fmt.Errorf("create entity value response error: %s", entityRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("create entity value parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("create entity value request error: %s", err)
	}

	return response, err
}

// remove a given value from an entity
//
// https://wit.ai/docs/http/20160516#delete--entities-:entity-id-values-link
func (c *Client) DeleteEntityValue(entityId, entityValue *string) (response map[string]string, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s/values/%s", *entityId, *entityValue), nil)

	var bytes []byte
	if bytes, err = c.request("DELETE", *url, nil); err == nil {
		var entityRes map[string]string
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			response = entityRes
		} else {
			err = fmt.Errorf("delete entity value parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("delete entity value request error: %s", err)
	}

	return response, err
}

// create a new expression for an entity
//
// https://wit.ai/docs/http/20160516#post--entities-:entity-id-values-:value-id-expressions-link
func (c *Client) CreateEntityExpression(entityId, entityValue, expression *string) (response Entity, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s/values/%s/expressions", *entityId, *entityValue), nil)

	body := map[string]interface{}{
		"expression": *expression,
	}

	var bytes []byte
	if bytes, err = c.request("POST", *url, body); err == nil {
		var entityRes Entity
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			if !entityRes.HasError() {
				response = entityRes
			} else {
				err = fmt.Errorf("create entity expression response error: %s", entityRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("create entity expression parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("create entity expression request error: %s", err)
	}

	return response, err
}

// remove an expression from an entity
//
// https://wit.ai/docs/http/20160516#delete--entities-:entity-id-values-:value-id-expressions-link
func (c *Client) DeleteEntityExpression(entityId, entityValue, expression *string) (response map[string]string, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/entities/%s/values/%s/expressions/%s", *entityId, *entityValue, *expression), nil)

	var bytes []byte
	if bytes, err = c.request("DELETE", *url, nil); err == nil {
		var entityRes map[string]string
		if err = json.Unmarshal(bytes, &entityRes); err == nil {
			response = entityRes
		} else {
			err = fmt.Errorf("delete entity expression parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("delete entity expression request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) create new intents
//
// https://wit.ai/docs/http/20160330#intents-post-link
// => https://wit.ai/docs/http/20160516#post--intents-(deprecated)-link
func (c *Client) CreateIntent_deprecated(intents ...Intent) (response Intents, err error) {
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
			if !intentsRes.HasError() {
				response = intentsRes
			} else {
				err = fmt.Errorf("new intents response error: %s", intentsRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("new intents parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("new intents request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) retrieve the list of all intents
//
// https://wit.ai/docs/http/20160330#intents-index-link
// => https://wit.ai/docs/http/20160516#get--intents-(deprecated)-link
func (c *Client) GetAllIntents_deprecated() (response []Intent, err error) {
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

// (DEPRECATED) retrieve all entities and expressions of an intent
//
// https://wit.ai/docs/http/20160330#intent-show-link
// => https://wit.ai/docs/http/20160516#get--intents-:intent-id-(deprecated)-link
func (c *Client) ShowIntent_deprecated(intentIdOrName *string) (response IntentDetail, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/intents/%s", *intentIdOrName), nil)

	var bytes []byte
	if bytes, err = c.request("GET", *url, nil); err == nil {
		var intentRes IntentDetail
		if err = json.Unmarshal(bytes, &intentRes); err == nil {
			if !intentRes.HasError() {
				response = intentRes
			} else {
				err = fmt.Errorf("show intent response error: %s", intentRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("show intent parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("show intent request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) update intent attributes
//
// https://wit.ai/docs/http/20160330#intent-put-link
// => https://wit.ai/docs/http/20160516#put--intents-:intent-id-(deprecated)-link
func (c *Client) UpdateIntentAttrs_deprecated(intentIdOrName, name, doc, metadata *string) (response IntentAttributes, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/intents/%s", *intentIdOrName), nil)

	body := map[string]interface{}{}
	if name != nil {
		body["name"] = *name
	}
	if doc != nil {
		body["doc"] = *doc
	}
	if metadata != nil {
		body["metadata"] = *metadata
	}

	var bytes []byte
	if bytes, err = c.request("PUT", *url, body); err == nil {
		var intentRes IntentAttributes
		if err = json.Unmarshal(bytes, &intentRes); err == nil {
			if !intentRes.HasError() {
				response = intentRes
			} else {
				err = fmt.Errorf("update intent attrs response error: %s", intentRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("update intent attrs parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("update intent attrs request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) add new expressions to an intent
//
// https://wit.ai/docs/http/20160330#create-intent-expressions-link
// => https://wit.ai/docs/http/20160516#post--intents-:intent-id-expressions-(deprecated)-link
func (c *Client) CreateIntentExpressions_deprecated(intentIdOrName *string, expressions ...string) (response []IntentExpressionCreated, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/intents/%s/expressions", *intentIdOrName), nil)

	body := []interface{}{}
	for _, expression := range expressions {
		body = append(body, map[string]string{"body": expression})
	}

	var bytes []byte
	if bytes, err = c.request("POST", *url, body); err == nil {
		var intentRes []IntentExpressionCreated
		if err = json.Unmarshal(bytes, &intentRes); err == nil {
			response = intentRes
		} else {
			err = fmt.Errorf("create intent expressions parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("create intent expressions request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) remove an expression from an intent
//
// https://wit.ai/docs/http/20160330#destroy-intent-expression-link
// => https://wit.ai/docs/http/20160516#delete--intents-:intent-id-expressions-:expression-id-(deprecated)-link
func (c *Client) DeleteIntentExpression_deprecated(intentIdOrName, expressionId *string) (response map[string]string, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/intents/%s/expressions/%s", *intentIdOrName, *expressionId), nil)

	var bytes []byte
	if bytes, err = c.request("DELETE", *url, nil); err == nil {
		var exprRes map[string]string
		if err = json.Unmarshal(bytes, &exprRes); err == nil {
			response = exprRes
		} else {
			err = fmt.Errorf("delete expression parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("delete expression request error: %s", err)
	}

	return response, err
}

// (DEPRECATED) retrieve an existing message
//
// https://wit.ai/docs/http/20160330#get-message-link
// => https://wit.ai/docs/http/20160516#get--messages-:msg-id-(deprecated)-link
func (c *Client) GetMessage_deprecated(messageId *string) (response Message, err error) {
	url := c.makeUrl(fmt.Sprintf("https://api.wit.ai/messages/%s", *messageId), nil)

	var bytes []byte
	if bytes, err = c.request("GET", *url, nil); err == nil {
		var msgRes Message
		if err = json.Unmarshal(bytes, &msgRes); err == nil {
			if !msgRes.HasError() {
				response = msgRes
			} else {
				err = fmt.Errorf("get message response error: %s", msgRes.ErrorMessage())
			}
		} else {
			err = fmt.Errorf("get message parse error: %s", err)
		}
	} else {
		err = fmt.Errorf("get message request error: %s", err)
	}

	return response, err
}
