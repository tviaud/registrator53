package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/registrator53/ec2"
	"github.com/registrator53/route53"
)

func main() {
	log.Println("Starting Registrator53")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("v1.37"))
	if err != nil {
		panic(err)
	}

	//Creating AWS Session
	sess := session.Must(session.NewSession())

	//Getting Private IP of the instance
	privateIP := ec2.GetPrivateIP(sess)
	events, _ := cli.Events(ctx, types.EventsOptions{})
	// if errors != nil {
	// 	panic(errors)
	// }
	for msg := range events {
		dnsName := msg.Actor.Attributes["name"] + ".els.vpc.local"
		switch msg.Status {
		case "start":
			log.Println("Container Starting:", msg.Actor.Attributes["name"])
			route53.CreateRoute53Record(sess, dnsName, privateIP, "ContainerTest", "DNSZONEID", "UPSERT")
		case "die":
			log.Println("Container Killing:", msg.Actor.Attributes["name"])
			route53.CreateRoute53Record(sess, dnsName, privateIP, "ContainerTest", "DNSZONEID", "DELETE")
		}
	}
}
