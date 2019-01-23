package kindleland

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// NewKeyboardListener takes the path to a /dev/input/* device and returns
// a chan on which events will be pushed when they occur.
func NewKeyboardListener(path string) (chan KeyboardEvent, error) {
	channel := make(chan KeyboardEvent)

	keyboard, err := os.Open(path)
	if err != nil {
		close(channel)
		return channel, err
	}
	buf := make([]byte, 16)

	go func() {
		for {
			if _, err := keyboard.Read(buf); err != nil {
				fmt.Println(err)
				break
			}
			var event Event
			if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &event); err != nil {
				break
			}

			kevent := KeyboardEvent{
				Time: time.Unix(int64(event.Time.Seconds), int64(event.Time.Microseconds)*1000),
				Type: KeyEventType(event.Value),
				Key:  KeyType(event.Code),
			}

			channel <- kevent
		}

		close(channel)
		keyboard.Close()
	}()

	return channel, nil
}
