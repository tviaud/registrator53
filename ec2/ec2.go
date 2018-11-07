package ec2

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

//isMetadataAvailable Return if Metadata Service is available
func isMetadataAvailable(sess *session.Session) bool {
	client := ec2metadata.New(sess)
	isAvailable := client.Available()
	if isAvailable == false {
		log.Println("EC2 Metadata not available")
		return false
	}
	log.Println("EC2 Metadata available")
	return true
}

//GetPrivateIP Return private IP of the EC2 instance
func GetPrivateIP(sess *session.Session) string {
	client := ec2metadata.New(sess)
	if isMetadataAvailable(sess) {
		privateIP, err := client.GetMetadata("local-ipv4")
		if err != nil {
			log.Panicln("Error retrieving Private IP")
		}
		log.Println("Private IP:", privateIP)
		return privateIP
	}
	log.Println("EC2 Metadata service not available")
	return ""
}
