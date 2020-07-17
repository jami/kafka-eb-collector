# kafka-eb-collector

Experimental Kafka event bus component that is able to collect asynchronous group events

## Approach

There are several common operations in a decentralized event system. One of them is that one event triggers multiple async events followed by a state that is waiting for the result of all of the triggered events.

![example process](./doc/process-example.png "Example process")

In some programming languages there exist a similar technique to collect the result of parallel thread routines.

golang:

    func main() {
        var wg sync.WaitGroup

        for i := 1; i <= 5; i++ {
            wg.Add(1)
            go worker(i, &wg)
        }

        wg.Wait()
    }

c++

    boost::thread_group

node

    Promise.all([p1, p2, p3]).then((v) => {
        resolved(v);
    })

## Examples

Install tools

    brew install kafkacat

### Example 1 - simple group configuration

Scenario: A process_A have to dispatch multiple jobs to decentralized services. After all of the jobs finished, a process_B should be triggered. A unidirectional communication strategy

                  |---- job A (async) ---->|
    process_A --->|---- job B (async) ---->|---> process_B
                  |---- job C (async) ---->|

Start the environment

    make compose/up

Listen to the default kafka topic for the collector with a consumer

    kafkacat -C -b localhost:9092 -t eventbus-collector -o beginning

Start a producer

    kafkacat -P -b localhost:9092 -t eventbus-collector

Create a collector group. We take 'F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81' as unique group id in this case. The default TTL of a group is 30 seconds.

    {
        "id": "F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "type": "EventGroupCollectorCreate",
        "expected": [
            "A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
            "366B21B5-E499-40E1-882A-B3EC99B27E47",
            "11ACC3BA-F0E5-4962-BC00-0E54497AB296"
        ],
        "ttl": 120,
        "payload": {
            "group": "payload"
        }
    }

Minify the json

    {"id":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","type":"EventGroupCollectorCreate","expected":["A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","366B21B5-E499-40E1-882A-B3EC99B27E47","11ACC3BA-F0E5-4962-BC00-0E54497AB296"],"ttl":120,"payload":{"group":"payload"}}

And post it through the producer.
Next step is to create the 3 jobs that completes the group

    {
        "id": "A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "type":"EventGroupCollectorEntityDone",
        "group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "result": {
            "job": "A",
            "bar": true
        }
    }

    {
        "id": "366B21B5-E499-40E1-882A-B3EC99B27E47",
        "type":"EventGroupCollectorEntityDone",
        "group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "result": {
            "job": "B",
            "bar": true
        }
    }

    {
        "id": "11ACC3BA-F0E5-4962-BC00-0E54497AB296",
        "type":"EventGroupCollectorEntityDone",
        "group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "result": {
            "job": "C",
            "bar": {
                "buzz": 42
            }
        }
    }

Minify and send it with the producer within the ttl time duration

    {"id":"A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","type":"EventGroupCollectorEntityDone","group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","result":{"job":"A","bar":true}}
    {"id":"366B21B5-E499-40E1-882A-B3EC99B27E47","type":"EventGroupCollectorEntityDone","group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","result":{"job":"B","bar":true}}
    {"id":"11ACC3BA-F0E5-4962-BC00-0E54497AB296","type":"EventGroupCollectorEntityDone","group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","result":{"job":"C","bar":{"buzz":42}}}

You should hit the success message

    {
        "type":"EventGroupCollectorSuccess",
        "id":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "created":"2020-07-17 22:36:26.402756 +0000 UTC m=+261.390776401",
        "resolved":"2020-07-17 22:44:58.9145535 +0000 UTC m=+774.492444901",
        "results":{
            "11ACC3BA-F0E5-4962-BC00-0E54497AB296":{
                "bar":{
                    "buzz":42
                },
                "job":"C"
            },
            "366B21B5-E499-40E1-882A-B3EC99B27E47":{
                "bar":true,
                "job":"B"
            },
            "A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81":{
                "bar":true,
                "job":"A"
            }
        },
        "payload":{
            "group":"payload"
        }
    }

### Propagated and time dependent failures

As a process you can propagate a failure to the group.

    {
        "id": "11ACC3BA-F0E5-4962-BC00-0E54497AB296",
        "type": "EventGroupCollectorEntityDone",
        "group": "F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "error": "Error while processing operation. Insufficient space on bounded volume"
    }

    {"id":"11ACC3BA-F0E5-4962-BC00-0E54497AB296","type":"EventGroupCollectorEntityDone","group":"F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81","error":"Error while processing operation. Insufficient space on bounded volume"}

The result is a failure event in the kafka topic similar to this

    {
        "type": "EventGroupCollectorFailure",
        "id": "F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "error": "Error while processing operation. Insufficient space on bounded volume",
        "created": "2020-07-17 22:52:19.474696 +0000 UTC m=+1215.573108301",
        "results": {
            "11ACC3BA-F0E5-4962-BC00-0E54497AB296": null,
            "366B21B5-E499-40E1-882A-B3EC99B27E47": {
                "bar": true,
                "job": "B"
            },
            "A0CABCA1-64AE-40A3-ABEB-0741ACF1FF81": {
                "bar": true,
                "job": "A"
            }
        },
        "payload": {
            "group": "payload"
        }
    }

A failure because of a timeout event looks very similar

    {
        "type": "EventGroupCollectorFailure",
        "id": "F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "entity": "F0CABCA1-64AE-40A3-ABEB-0741ACF1FF81",
        "error": "Operation timed out after 10 seconds",
        "created": "2020-07-17T23:38:21Z",
        "results": {},
        "payload": {
            "group": "payload"
        }
    }