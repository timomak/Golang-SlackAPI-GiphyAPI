// DO NOT EDIT THIS FILE. This is a fully complete implementation.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/nlopes/slack"
	"gopkg.in/go-playground/webhooks.v5/github"

	// "github.com/droxey/goslackit/slack"
	_ "github.com/joho/godotenv/autoload"
)

func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

// sendHelp is a working help message, for reference.
func NotifySlackChannel(slackClient *slack.RTM, message string) {
	// if strings.ToLower(message) != "help" {
	// 	return
	// }
	slackClient.SendMessage(slackClient.NewOutgoingMessage(message, "portal-devs"))
}

// // GetSlackClient will return an authenticated slack
// // api object that you can then run slack methods on
// func GetSlackClient() *slack.Client {
//     return slack.New(os.Getenv("SLACK_TOKEN"))
// }
// // SiteDownMessage sends a message in slack channel when
// // a site returns an error code
// func SiteDownMessage(statusCode int, url, loadbalancerIP string) error {
//     client := GetSlackClient()
//     // get the error message from the status code
//     statusMsg := http.StatusText(statusCode)
//     // message to send
//     notifyMsg :=
//     slackMsg := slack.MsgOptionText(notifyMsg, false)
//     // send the message in the group
//     _, timestamp, err := client.PostMessage("portal-devs", slackMsg)
//     if err != nil {
//         return err
//     }

// Trying again
// main is our entrypoint, where the application initializes the Slackbot.
func main() {
	// port := ":" + os.Getenv("PORT")
	// go http.ListenAndServe(port, nil)
	hook, _ := github.New(github.Options.Secret("Wassup"))

	e := echo.New()
	e.POST("/test", func(c echo.Context) error {
		payload, err := hook.Parse(c.Request(), github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			release := payload.(github.PushPayload)
			// Do whatever you want from here...
			// fmt.Printf("%+v", release)
			// fmt.Printf("EMAIL: %+v", release.Pusher.Email)
			newMessage := string(release.Pusher.Name) + " just made a commit to the " + string(release.Repository.FullName) + "repo.\nLook at the changes: " + string(release.Repository.HTMLURL) + "\n"

			fmt.Printf("%+v just made a commit to the %+v repo.\nLook at the changes: %+v\n", release.Pusher.Name, release.Repository.FullName, release.Repository.HTMLURL)
			slackIt(newMessage)
			// case github.PullRequestPayload: .
			// 	pullRequest := payload.(github.PullRequestPayload)
			// 	// Do whatever you want from here...
			// 	fmt.Printf("%+v", pullRequest)

		}

		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}

// slackIt is a function that initializes the Slackbot.
func slackIt(message string) {
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := CreateSlackClient(botToken)
	fmt.Println("SENDING MESSASSAGE TO SLACK CHANNEL:", message)
	NotifySlackChannel(slackClient, message)
}
