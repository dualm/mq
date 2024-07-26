package main

import (
	"context"
	"gitee.com/ziegate/mq/rabbitmq"
	"log"
	"sync"
)

func main() {
	producer, err := rabbitmq.NewTopicProducer(
		&rabbitmq.DialOption{
			Username: "manager",
			Password: "admin123",
			Host:     "10.1.36.190",
			Port:     "5671",
			VHost:    "/mes",
		},
		&rabbitmq.ExchangeOption{
			Name:       "MES.PRD.EAP.SVR",
			Kind:       "topic",
			AutoDelete: false,
			Durable:    true,
			NoWait:     false,
			Args:       nil,
		},
		&rabbitmq.ExchangeOption{
			Name:       "MES.PRD.EAP.SVR",
			Kind:       "topic",
			AutoDelete: false,
			Durable:    true,
			NoWait:     false,
			Args:       nil,
		},
		&rabbitmq.QueueOption{
			Name:       "MES.DEV.EDGE.G2",
			AutoDelete: false,
			Durable:    false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		"KEY.MES.REPLYTO.EGATE",
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	count := 1024 * 16
	wg.Add(count)
	for range count {
		go func(ctx context.Context) {
			//ctx1, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			//defer cancel()

			defer wg.Done()
			_, err := producer.SendRequest(
				ctx,
				"MES.PRD.EAP.PEXsvr",
				false, false,
				[]byte(
					`			<Message>
			<Header>
			<MESSAGENAME>AreYouThereRequest</MESSAGENAME>
			<SHOPNAME>Z10000</SHOPNAME>
			<MACHINENAME>Z1COCK01</MACHINENAME>
			<TRANSACTIONID>20240726151621152951</TRANSACTIONID>
			<ORIGINALSOURCESUBJECTNAME>PKAOWRQPJH-Z1TYYT01-AreYouThereRequest</ORIGINALSOURCESUBJECTNAME>
			<SOURCESUBJECTNAME>ZL.FMM.MES.TEST.8MMPI91</SOURCESUBJECTNAME>
			<TARGETSUBJECTNAME>ZL.FMM.MES.TEST.PEXsvr</TARGETSUBJECTNAME>
			<EVENTUSER>Z1TYYT01</EVENTUSER>
			<EVENTCOMMENT>AreYouThereRequest</EVENTCOMMENT>
			</Header>
			<Body>
			<MACHINENAME>Z1TYYT01</MACHINENAME>
			</Body>
			<Return>
			<RETURNCODE></RETURNCODE>
			<RETURNMESSAGE></RETURNMESSAGE>
			</Return>
			</Message>
`),
			)
			if err != nil {
				log.Println(err)

				return
			}
		}(ctx)
	}

	wg.Wait()
}
