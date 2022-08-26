//+build !test

package nats

import (
	"log"
	"os"
	"sync"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
	"github.com/nats-io/nats.go"
)

//printMsg : To print when a msg is recieved
func printMsg(m *nats.Msg, i int) {
	defer util.Panic()
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

//Subscribe : This function is used to initiate subscriber
func Subscribe() {
	defer util.Panic()
	var wg sync.WaitGroup
	nc := config.NC
	i := 0
	wg.Add(1)
	subject := config.Conf.NatsServer.Subject
	queue := config.Conf.NatsServer.Queue
	var payload []byte
	nc.QueueSubscribe(subject+".documents.*", queue, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		go func(count int) {
			switch msg.Subject {
			case subject + ".documents.checklink":
				payload = Getlink(msg.Data)
			case subject + ".documents.getlink":
				payload = GetPresignedUrl(msg.Data)
			case subject + ".documents.delete":
				payload = DeleteFile(msg.Data)
			}
			err := msg.Respond(payload)
			if err != nil {
				log.Println(err)
			}
		}(i)
	})
	log.Println("listening on XSOnBoarding document manage service")
	wg.Wait()
}
