package src

import (
	"encoding/json"
	"fmt"
)

// GroupCollectorEnumeration type
type GroupCollectorEnumeration int

const (
	// EventGroupCollectorCreate event
	EventGroupCollectorCreate GroupCollectorEnumeration = iota
	// EventGroupCollectorEntityDone event
	EventGroupCollectorEntityDone
	// EventGroupCollectorFailure event
	EventGroupCollectorFailure
	// EventGroupCollectorSuccess event
	EventGroupCollectorSuccess
)

// FromString reads enum from string
func (e *GroupCollectorEnumeration) FromString(s string) error {
	switch s {
	case "EventGroupCollectorCreate":
		*e = EventGroupCollectorCreate
	case "EventGroupCollectorEntityDone":
		*e = EventGroupCollectorEntityDone
	case "EventGroupCollectorFailure":
		*e = EventGroupCollectorFailure
	case "EventGroupCollectorSuccess":
		*e = EventGroupCollectorSuccess
	default:
		return fmt.Errorf("Unable to parse '%s' into GroupCollectorEnumeration", s)
	}

	return nil
}

func (e GroupCollectorEnumeration) String() string {
	ret := "undefined"

	switch e {
	case EventGroupCollectorCreate:
		ret = "EventGroupCollectorCreate"
	case EventGroupCollectorEntityDone:
		ret = "EventGroupCollectorEntityDone"
	case EventGroupCollectorFailure:
		ret = "EventGroupCollectorFailure"
	case EventGroupCollectorSuccess:
		ret = "EventGroupCollectorSuccess"
	}

	return ret
}

// EventEnvelop base of all event data
type EventEnvelop struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Sender string `json:"sender,omitempty"`
}

// FromBytes unmarshal byte stream to envelop
func (ee *EventEnvelop) FromBytes(b []byte) error {
	if err := json.Unmarshal(b, ee); err != nil {
		return err
	}

	return nil
}

// EnvelopGroupHandler model
type EnvelopGroupHandler struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

// EventEnvelopGroup model
type EventEnvelopGroup struct {
	EventEnvelop
	ExpectedIDs []string                       `json:"expected"`
	TTL         int                            `json:"ttl"`
	CreatedTime string                         `json:"created"`
	Handler     map[string]EnvelopGroupHandler `json:"handler"`
	Payload     interface{}                    `json:"payload"`
	Results     map[string]interface{}         `json:"results"`
}

// FromBytes creates
func (eec *EventEnvelopGroup) FromBytes(b []byte) error {
	return json.Unmarshal(b, eec)
}

// Resolve appends entity result data to group
func (eec *EventEnvelopGroup) Resolve(e *EventEnvelopEntityDone) {
	containsID := false
	for _, v := range eec.ExpectedIDs {
		if v == e.ID {
			containsID = true
			break
		}
	}

	if !containsID {
		return
	}

	eec.Results[e.ID] = e.Result
}

// IsComplete checks for completness
func (eec *EventEnvelopGroup) IsComplete() bool {
	isComplete := true
	for _, v := range eec.ExpectedIDs {
		if _, ok := eec.Results[v]; !ok {
			isComplete = false
			break
		}
	}

	return isComplete
}

// EventEnvelopEntityDone model
type EventEnvelopEntityDone struct {
	EventEnvelop
	Error   string      `json:"error,omitempty"`
	GroupID string      `json:"group"`
	Result  interface{} `json:"result"`
}

// FromBytes creates
func (eeed *EventEnvelopEntityDone) FromBytes(b []byte) error {
	return json.Unmarshal(b, eeed)
}

// EventEnvelopFailure model
type EventEnvelopFailure struct {
	EventEnvelop
	Error        string                 `json:"error,omitempty"`
	CreatedTime  string                 `json:"created"`
	ResolvedTime string                 `json:"resolved"`
	Results      map[string]interface{} `json:"results"`
	Payload      interface{}            `json:"payload"`
}

// EventEnvelopSuccess model
type EventEnvelopSuccess struct {
	EventEnvelop
	CreatedTime  string                 `json:"created"`
	ResolvedTime string                 `json:"resolved"`
	Results      map[string]interface{} `json:"results"`
	Payload      interface{}            `json:"payload"`
}
