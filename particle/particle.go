package particle

import (
	"encoding/json"
	"log"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/simpleiot/simpleiot/data"
)

// Event from particle
type Event struct {
	Data      string    `json:"data"`
	TTL       uint32    `json:"ttl"`
	Timestamp time.Time `json:"published_at"`
	CoreID    string    `json:"coreid"`
}

const url string = "https://api.particle.io/v1/devices/events/"

// SampleReader does a streaming http read and returns when the connection closes
func SampleReader(eventPrefix, token string, callback func(string, []data.Sample)) error {
	urlAuth := url + eventPrefix + "?access_token=" + token

	stream, err := eventsource.Subscribe(urlAuth, "")

	if err != nil {
		return err
	}

	for {
		select {
		case event := <-stream.Events:
			var pEvent Event
			err := json.Unmarshal([]byte(event.Data()), &pEvent)
			if err != nil {
				log.Println("Got error decoding particle event: ", err)
				continue
			}

			var samples []data.Sample
			err = json.Unmarshal([]byte(pEvent.Data), &samples)
			if err != nil {
				log.Println("Got error decoding samples: ", err)
				continue
			}

			callback(pEvent.CoreID, samples)

		case err := <-stream.Errors:
			log.Println("Got error: ", err)
		}
	}
}
