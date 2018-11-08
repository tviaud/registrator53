package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/registrator53/ec2"
	"github.com/registrator53/route53"
)

func main() {
	log.Println("Starting Registrator53")

	//Getting Environment variables
	domainName, ok := os.LookupEnv("DOMAIN_NAME")
	if !ok {
		log.Println("Domain Name not set")
		panic("Domain Name not set")
	} else {
		//Prepending . for DNS convention and Route53
		domainName = "." + domainName
	}
	hostedZoneID, ok := os.LookupEnv("HOSTED_ZONE_ID")
	if !ok {
		panic("Zone_ID not set")
	}

	ctx := context.Background()
	//Creating docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("v1.37"))
	if err != nil {
		panic(err)
	}

	//Creating AWS Session
	sess := session.Must(session.NewSession())

	//Getting Private IP of the instance
	privateIP := ec2.GetPrivateIP(sess)
	events, _ := cli.Events(ctx, types.EventsOptions{})

	for msg := range events {
		containerName := msg.Actor.Attributes["name"]
		dnsName := containerName + domainName
		switch msg.Status {
		case "start":
			log.Println("Container Starting:", containerName)
			result := route53.CreateRoute53Record(sess, dnsName, privateIP, containerName, hostedZoneID, "UPSERT")
			if result == false {
				log.Println("Record not Created")
			}
		case "die":
			log.Println("Container Killing:", containerName)
			result := route53.CreateRoute53Record(sess, dnsName, privateIP, containerName, hostedZoneID, "DELETE")
			if result == false {
				log.Println("Record not Removed")
			}
		}
	}
}
