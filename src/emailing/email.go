package emailing

import (
	"fmt"
	"log"
	"os"

	//   "os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

var (
	Emails      emails
	apikey      string = "e3666c5ac12e90ce7826cd2382f767f2"
	secretkey   string = "479a0020bef2201427f63e146b3e33e8"
	Customid    string = "1234567"
	Owner       string = "MY Synca"
	OwnerEmail  string = "myrachanto@gmal.com"
	WebsiteLink string = ""
	Phone       string = "123456777"
)

type emails struct {
}

func (e *emails) Emailing(status bool) {
	if status {
		Customid := os.Getenv("Customid")
		Owner := os.Getenv("Owner")
		OwnerEmail := os.Getenv("OwnerEmail")
		mailjetClient := mailjet.NewMailjetClient(apikey, secretkey)
		messagesInfo := []mailjet.InfoMessagesV31{
			mailjet.InfoMessagesV31{
				From: &mailjet.RecipientV31{
					Email: OwnerEmail,
					Name:  Owner,
				},
				To: &mailjet.RecipientsV31{
					mailjet.RecipientV31{
						Email: OwnerEmail,
						Name:  Owner,
					},
				},
				Subject:  "Hellos",
				TextPart: "thank you forshoppign with us",
				HTMLPart: "<h3>Database Synchronizing done successifully!</h3><br />Have a lovevely day!",
				CustomID: Customid,
			},
		}
		messages := mailjet.MessagesV31{Info: messagesInfo}
		res, err := mailjetClient.SendMailV31(&messages)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Data: %+v\n", res)

	} else {
		Customid := os.Getenv("Customid")
		Owner := os.Getenv("Owner")
		OwnerEmail := os.Getenv("OwnerEmail")
		mailjetClient := mailjet.NewMailjetClient(apikey, secretkey)
		messagesInfo := []mailjet.InfoMessagesV31{
			mailjet.InfoMessagesV31{
				From: &mailjet.RecipientV31{
					Email: OwnerEmail,
					Name:  Owner,
				},
				To: &mailjet.RecipientsV31{
					mailjet.RecipientV31{
						Email: OwnerEmail,
						Name:  Owner,
					},
				},
				Subject:  "Hellos",
				TextPart: "thank you forshoppign with us",
				HTMLPart: "<h3>Database Synchronizing Failed!</h3><br />Have a lovevely day!",
				CustomID: Customid,
			},
		}
		messages := mailjet.MessagesV31{Info: messagesInfo}
		res, err := mailjetClient.SendMailV31(&messages)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Data: %+v\n", res)
	}
}
