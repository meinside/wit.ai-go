// types_helper.go: helper functions for types

package witai

import (
	"fmt"
	"strings"
)

// https://wit.ai/docs/http/20160330#converse-link
func (c Converse) String() string {
	attrs := []string{}
	if c.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *c.Error))
	}
	if c.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *c.Code))
	}
	if c.Type != nil {
		attrs = append(attrs, fmt.Sprintf("Type: %s", *c.Type))
	}
	if c.Message != nil {
		attrs = append(attrs, fmt.Sprintf("Message: %s", *c.Message))
	}
	if c.Entities != nil {
		attrs = append(attrs, fmt.Sprintf("Entities: %v", c.Entities))
	}
	attrs = append(attrs, fmt.Sprintf("Confidence: %.6f", c.Confidence))

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

// https://wit.ai/docs/http/20160330#context-link
func (c Context) String() string {
	attrs := []string{}
	attrs = append(attrs, fmt.Sprintf("State: %v", c.State))
	if c.ReferenceTime != nil {
		attrs = append(attrs, fmt.Sprintf("ReferenceTime: %s", *c.ReferenceTime))
	}
	if c.TimeZone != nil {
		attrs = append(attrs, fmt.Sprintf("TimeZone: %s", *c.TimeZone))
	}
	if c.Entities != nil {
		attrs = append(attrs, fmt.Sprintf("Entities: %v", c.Entities))
	}
	if c.Location != nil {
		attrs = append(attrs, fmt.Sprintf("Location: %v", c.Location))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (e Entities) String() string {
	attrs := []string{}
	if e.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *e.Id))
	}
	if e.Doc != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *e.Doc))
	}
	if len(e.Values) > 0 {
		attrs = append(attrs, fmt.Sprintf("Values: %v", e.Values))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (l Location) String() string {
	return fmt.Sprintf("{Latitude: %.6f, Longitude: %.6f}", l.Latitude, l.Longitude)
}

// https://wit.ai/docs/http/20160330#get-intent-via-text-link
func (m Message) String() string {
	attrs := []string{}
	if m.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *m.Error))
	}
	if m.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *m.Code))
	}
	if m.MessageId != nil {
		attrs = append(attrs, fmt.Sprintf("MessageId: %s", *m.MessageId))
	}
	if m.Text != nil {
		attrs = append(attrs, fmt.Sprintf("Text: %s", *m.Text))
	}
	if len(m.Outcomes) > 0 {
		attrs = append(attrs, fmt.Sprintf("Outcomes: %v", m.Outcomes))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (o Outcome) String() string {
	attrs := []string{}
	if o.Text != nil {
		attrs = append(attrs, fmt.Sprintf("Text: %s", *o.Text))
	}
	if o.Intent != nil {
		attrs = append(attrs, fmt.Sprintf("Intent: %s", *o.Intent))
	}
	if o.Entities != nil {
		attrs = append(attrs, fmt.Sprintf("Entities: %v", o.Entities))
	}
	attrs = append(attrs, fmt.Sprintf("Confidence: %.6f", o.Confidence))

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i IntentExpression) String() string {
	attrs := []string{}
	if i.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *i.Id))
	}
	if i.Body != nil {
		attrs = append(attrs, fmt.Sprintf("Body: %s", *i.Body))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i Intent) String() string {
	attrs := []string{}
	if i.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *i.Error))
	}
	if i.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *i.Code))
	}
	if i.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *i.Id))
	}
	if i.Name != nil {
		attrs = append(attrs, fmt.Sprintf("Name: %s", *i.Name))
	}
	if i.Doc != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *i.Doc))
	}
	if i.Metadata != nil {
		attrs = append(attrs, fmt.Sprintf("Metadata: %s", *i.Metadata))
	}
	if len(i.Expressions) > 0 {
		attrs = append(attrs, fmt.Sprintf("Expressions: %v", i.Expressions))
	}
	if i.Meta != nil {
		attrs = append(attrs, fmt.Sprintf("Meta: %v", i.Meta))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i Intents) String() string {
	attrs := []string{}
	if i.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *i.Error))
	}
	if i.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *i.Code))
	}
	if len(i.Intents) > 0 {
		attrs = append(attrs, fmt.Sprintf("Intents: %v", i.Intents))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i IntentDetail) String() string {
	attrs := []string{}
	if i.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *i.Error))
	}
	if i.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *i.Code))
	}
	if i.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *i.Id))
	}
	if i.Name != nil {
		attrs = append(attrs, fmt.Sprintf("Name: %s", *i.Name))
	}
	if i.Doc != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *i.Doc))
	}
	if len(i.Expressions) > 0 {
		attrs = append(attrs, fmt.Sprintf("Expressions: %s", i.Expressions))
	}
	if i.Entities != nil {
		attrs = append(attrs, fmt.Sprintf("Entities: %s", i.Entities))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i IntentDetailExpression) String() string {
	attrs := []string{}
	if i.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *i.Id))
	}
	if i.Body != nil {
		attrs = append(attrs, fmt.Sprintf("Body: %s", *i.Body))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i IntentExpressionCreated) String() string {
	attrs := []string{}
	if i.IntentId != nil {
		attrs = append(attrs, fmt.Sprintf("IntentId: %s", *i.IntentId))
	}
	if i.Body != nil {
		attrs = append(attrs, fmt.Sprintf("Body: %s", *i.Body))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (i IntentAttributes) String() string {
	attrs := []string{}
	if i.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *i.Error))
	}
	if i.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *i.Code))
	}
	if i.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *i.Id))
	}
	if i.Name != nil {
		attrs = append(attrs, fmt.Sprintf("Name: %s", *i.Name))
	}
	if i.Metadata != nil {
		attrs = append(attrs, fmt.Sprintf("Metadata: %s", *i.Metadata))
	}
	if i.Doc != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *i.Doc))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (e Entity) String() string {
	attrs := []string{}
	if e.Error != nil {
		attrs = append(attrs, fmt.Sprintf("Error: %s", *e.Error))
	}
	if e.Code != nil {
		attrs = append(attrs, fmt.Sprintf("Code: %s", *e.Code))
	}
	if e.Id != nil {
		attrs = append(attrs, fmt.Sprintf("Id: %s", *e.Id))
	}
	if e.Name != nil {
		attrs = append(attrs, fmt.Sprintf("Name: %s", *e.Name))
	}
	if e.Doc != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *e.Doc))
	}
	if e.Lang != nil {
		attrs = append(attrs, fmt.Sprintf("Lang: %s", *e.Lang))
	}
	attrs = append(attrs, fmt.Sprintf("Closed: %t", e.Closed))
	attrs = append(attrs, fmt.Sprintf("Exotic: %t", e.Exotic))
	attrs = append(attrs, fmt.Sprintf("Builtin: %t", e.Builtin))
	if len(e.Values) > 0 {
		attrs = append(attrs, fmt.Sprintf("Values: %v", e.Values))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

func (e EntityValue) String() string {
	attrs := []string{}
	if len(e.Expressions) > 0 {
		attrs = append(attrs, fmt.Sprintf("Expressions: %v", e.Expressions))
	}
	if e.Value != nil {
		attrs = append(attrs, fmt.Sprintf("Doc: %s", *e.Value))
	}

	return fmt.Sprintf("{%s}", strings.Join(attrs, ", "))
}

// helper functions

func (r ResponseError) HasError() bool { // XXX - due to inconsistency in response formats
	if r.Error != nil || len(r.Errors) > 0 || r.Body != nil {
		return true
	}
	return false
}

func (r ResponseError) ErrorMessage() string { // XXX - due to inconsistency in response formats
	errors := append([]string{}, r.Errors...)
	if r.Error != nil {
		errors = append(errors, *r.Error)
	}
	if r.Body != nil {
		errors = append(errors, *r.Body)
	}
	return strings.Join(errors, ",")
}

func NewIntentExpression(body string) IntentExpression {
	return IntentExpression{
		Body: &body,
	}
}

func NewIntentExpressions(body ...string) []IntentExpression {
	expressions := []IntentExpression{}
	for _, b := range body {
		expressions = append(expressions, NewIntentExpression(b))
	}
	return expressions
}
