package src

// GroupCollectorEnumeration type
type GroupCollectorEnumeration int

const (
	// EventGroupCollectorCreate event
	EventGroupCollectorCreate GroupCollectorEnumeration = iota
	// EventGroupCollectorAppend event
	EventGroupCollectorAppend
	// EventGroupCollectorFailure event
	EventGroupCollectorFailure
	// EventGroupCollectorSuccess event
	EventGroupCollectorSuccess
)

func (e GroupCollectorEnumeration) String() string {
	ret := "undefined"

	switch e {
	case EventGroupCollectorCreate:
		ret = "EventGroupCollectorCreate"
	case EventGroupCollectorAppend:
		ret = "EventGroupCollectorAppend"
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
	Sender string `json:"sender"`
}

// EventEnvelopCreate model
type EventEnvelopCreate struct {
	EventEnvelop
}

// EventEnvelopAppend model
type EventEnvelopAppend struct {
	EventEnvelop
}

// EventEnvelopFailure model
type EventEnvelopFailure struct {
	EventEnvelop
}

// EventEnvelopSuccess model
type EventEnvelopSuccess struct {
	EventEnvelop
}
