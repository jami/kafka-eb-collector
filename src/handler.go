package src

import (
	"encoding/json"
	"fmt"
	"time"
)

// GroupCreateHandler creates a new group
func GroupCreateHandler(store CollectorStore, b []byte, eventbusProducer *Producer) error {
	header := EventEnvelopGroup{}
	if err := header.FromBytes(b); err != nil {
		return err
	}

	header.CreatedTime = time.Now().String()
	header.Results = map[string]interface{}{}

	data, _ := json.Marshal(header)
	store.Put(header.ID, data)
	return nil
}

// EntityDoneHandler closes a entity on the group
func EntityDoneHandler(store CollectorStore, b []byte, eventbusProducer *Producer) error {
	header := EventEnvelopEntityDone{}
	if err := header.FromBytes(b); err != nil {
		return err
	}

	if data, err := store.Get(header.GroupID); err == nil {
		var group EventEnvelopGroup
		if err := json.Unmarshal(data, &group); err != nil {
			return fmt.Errorf("Error while unmarshaling group '%s'. %s", header.GroupID, err.Error())
		}

		group.Resolve(&header)

		if header.Error != "" {
			failure := EventEnvelopFailure{
				EventEnvelop: EventEnvelop{
					Type: "EventGroupCollectorFailure",
					ID:   header.GroupID,
				},
				Error:       header.Error,
				CreatedTime: group.CreatedTime,
				Results:     group.Results,
				Payload:     group.Payload,
			}
			eventbusProducer.SendJSON(failure)
			store.Delete(header.GroupID)
			return nil
		}

		if group.IsComplete() {
			success := EventEnvelopSuccess{
				EventEnvelop: EventEnvelop{
					Type: "EventGroupCollectorSuccess",
					ID:   header.GroupID,
				},
				CreatedTime:  group.CreatedTime,
				ResolvedTime: time.Now().String(),
				Results:      group.Results,
				Payload:      group.Payload,
			}
			eventbusProducer.SendJSON(success)
			store.Delete(header.GroupID)
			return nil
		}

		data, _ := json.Marshal(group)
		store.Put(header.GroupID, data)
	} else {
		return fmt.Errorf("Entity '%s' got unknown group id '%s'", header.ID, header.GroupID)
	}

	return nil
}
