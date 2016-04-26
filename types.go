// types

package witai

import (
	"fmt"
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
	Confidence float32     `json:"confidence"`
}

func (o Outcome) String() string {
	return fmt.Sprintf("{Text: %s, Intent: %s, Entities: %v, Confidence: %.6f}", *o.Text, *o.Intent, o.Entities, o.Confidence)
}

type Intent struct {
	Name        *string       `json:"name"`
	Doc         *string       `json:"doc,omitempty"`
	Metadata    *string       `json:"metadata,omitempty"`
	Expressions []interface{} `json:"expressions,omitempty"`
	Meta        interface{}   `json:"meta,omitempty"`
}

func (i Intent) String() string {
	return fmt.Sprintf("{Name: %s, Doc: %s, Metadata: %s, Expressions: %v, Meta: %v}", *i.Name, *i.Doc, *i.Metadata, i.Expressions, i.Meta)
}

type ResponseIntents struct {
	ResponseError

	Intents []interface{} `json:"intents"`
}
