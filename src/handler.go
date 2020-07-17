package src

import (
	"encoding/json"
	"fmt"
	"time"
)

// GroupTimeoutHandler handles the overdue groups
func GroupTimeoutHandler(store CollectorStore, eventbusProducer *Producer) error {
	store.IterateAll(func(key string, data []byte) {
		fmt.Printf("Iterate key %s", key)
		header := EventEnvelopGroup{}
		if err := header.FromBytes(data); err != nil {
			return
		}

		if created, err := time.Parse(time.RFC3339, header.CreatedTime); err == nil {
			till := created.Add(time.Second * time.Duration(header.TTL))
			if till.Before(time.Now()) {
				fmt.Printf("Timeout key %s", key)

				failure := EventEnvelopFailure{
					EventEnvelop: EventEnvelop{
						Type: "EventGroupCollectorFailure",
						ID:   header.ID,
					},
					FailureEntityID: header.ID,
					Error:           fmt.Sprintf("Operation timed out after %d seconds", header.TTL),
					CreatedTime:     header.CreatedTime,
					ResolvedTime:    time.Now().Format(time.RFC3339),
					Results:         header.Results,
					Payload:         header.Payload,
				}
				eventbusProducer.SendJSON(failure)
				store.Delete(header.ID)
			}
		}
	})

	return nil
}

// GroupCreateHandler creates a new group
func GroupCreateHandler(store CollectorStore, b []byte, eventbusProducer *Producer) error {
	header := EventEnvelopGroup{}
	if err := header.FromBytes(b); err != nil {
		return err
	}

	if header.TTL <= 0 {
		header.TTL = DefaultGroupProcessingTTL
	}

	header.CreatedTime = time.Now().Format(time.RFC3339)
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
				FailureEntityID: header.ID,
				Error:           header.Error,
				CreatedTime:     group.CreatedTime,
				ResolvedTime:    time.Now().Format(time.RFC3339),
				Results:         group.Results,
				Payload:         group.Payload,
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
				ResolvedTime: time.Now().Format(time.RFC3339),
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
