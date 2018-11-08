package route53

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

//CreateRoute53Record Create AWS session
func CreateRoute53Record(sess *session.Session, dnsName string, ip string, containerName string, hostedZoneID string, actionRoute53 string) bool {
	svc := route53.New(sess)
	recordSetInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(actionRoute53),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(dnsName),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String("1 10 5269 " + dnsName),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("SRV"),
					},
				},
			},
			Comment: aws.String("Insert DNS record for " + containerName),
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	log.Println(actionRoute53 + " " + dnsName + " Record in Route53 with IP:" + ip)

	_, err := svc.ChangeResourceRecordSets(recordSetInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				log.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
				return false
			case route53.ErrCodeInvalidChangeBatch:
				log.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
				return false
			case route53.ErrCodeInvalidInput:
				log.Println(route53.ErrCodeInvalidInput, aerr.Error())
				return false
			case route53.ErrCodePriorRequestNotComplete:
				log.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
				return false
			default:
				log.Println(aerr.Error())
				return false
			}
		}
	}
	return true
}
