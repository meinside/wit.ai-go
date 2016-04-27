// types.go: responses and data types

package witai

type Client struct {
	Token   *string
	Version *string

	headerAuth   *string
	headerAccept *string

	Verbose bool
}

// https://wit.ai/docs/http/20160330#response-format-link
type ResponseError struct {
	Error *string `json:"error,omitempty"`
	Code  *string `json:"code,omitempty"`
}

// https://wit.ai/docs/http/20160330#converse-link
type Converse struct {
	ResponseError

	Type       *string                `json:"type,omitempty"`
	Message    *string                `json:"msg,omitempty"`
	Action     *string                `json:"action,omitempty"`
	Entities   map[string]interface{} `json:"entities,omitempty"`
	Confidence float32                `json:"confidence"`
}

// https://wit.ai/docs/http/20160330#context-link
type Context struct {
	State         interface{} `json:"state,omitempty"`
	ReferenceTime *string     `json:"reference_time,omitempty"`
	TimeZone      *string     `json:"timezone,omitempty"`
	Entities      *Entities   `json:"entities,omitempty"`
	Location      *Location   `json:"location,omitempty"`
}

type Entities struct {
	Id     *string       `json:"id,omitempty"`
	Doc    *string       `json:"doc,omitempty"`
	Values []EntityValue `json:"values,omitempty"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// https://wit.ai/docs/http/20160330#get-intent-via-text-link
type Message struct {
	ResponseError

	MessageId *string   `json:"msg_id"`
	Text      *string   `json:"_text"`
	Outcomes  []Outcome `json:"outcomes"`
}

type Outcome struct {
	Text       *string                `json:"_text"`
	Intent     *string                `json:"intent"`
	Entities   map[string]interface{} `json:"entities"`
	Confidence float32                `json:"confidence"`
}

type Intent struct {
	ResponseError

	Id          *string            `json:"id,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Doc         *string            `json:"doc,omitempty"`
	Metadata    *string            `json:"metadata,omitempty"`
	Expressions []IntentExpression `json:"expressions,omitempty"`
	Meta        interface{}        `json:"meta,omitempty"`
}

type IntentExpression struct {
	Id   *string `json:"id,omitempty"`
	Body *string `json:"body,omitempty"`
}

type Intents struct {
	ResponseError

	Intents []Intent `json:"intents"`
}

type IntentDetail struct {
	ResponseError

	Id          *string                  `json:"id"`
	Name        *string                  `json:"name"`
	Doc         *string                  `json:"doc"`
	Expressions []IntentDetailExpression `json:"expressions"`
	Entities    []interface{}            `json:"entities"`
}

type IntentDetailExpression struct {
	Id   *string `json:"id,omitempty"`
	Body *string `json:"body,omitempty"`
}

type IntentExpressionCreated struct {
	IntentId *string `json:"intent_id,omitempty"`
	Body     *string `json:"body,omitempty"`
}

type IntentAttributes struct {
	ResponseError

	Id       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Metadata *string `json:"metadata,omitempty"`
	Doc      *string `json:"doc,omitempty"`
}

type Entity struct {
	ResponseError

	Id      *string       `json:"id"`
	Name    *string       `json:"name"`
	Doc     *string       `json:"doc"`
	Lang    *string       `json:"lang"`
	Closed  bool          `json:"closed"`
	Exotic  bool          `json:"exotic"`
	Builtin bool          `json:"builtin"`
	Values  []EntityValue `json:"values"`
}

type EntityValue struct {
	Expressions []string `json:"expressions"`
	Value       *string  `json:"value"`
}
