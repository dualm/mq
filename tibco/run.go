package tibco

import (
	"fmt"

	"github.com/dualm/mq"
)

func Send(targetSubjectName, fieldName string, c chan<- []byte, msg []mq.MqMessage) {
	for _, m := range msg {
		if m.IsEvent {
			err := tibSend(targetSubjectName, fieldName, string(m.Msg))

			if err != nil {
				fmt.Println(err)

				return
			}
		} else {
			re, err := tibSendRequest(targetSubjectName, fieldName, string(m.Msg))
			if err != nil {
				fmt.Printf("%s, %s\n", err, m.Msg)

				c <- nil

				continue
			}

			c <- []byte(re)
		}
	}
}
