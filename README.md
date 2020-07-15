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
