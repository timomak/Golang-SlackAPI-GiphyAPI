package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "type in '@BOT_NAME <command_arg_1> <command_arg_2>'"

type ImageData struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
	Frames string `json:"frames"`
}
type Gif struct {
	Type               string `json:"type"`
	Id                 string `json:"id"`
	URL                string `json:"url"`
	Tags               string `json:"tags"`
	BitlyGifURL        string `json:"bitly_gif_url"`
	BitlyFullscreenURL string `json:"bitly_fullscreen_url"`
	BitlyTiledURL      string `json:"bitly_tiled_url"`
	EmbedURL           string `json:"embed_url`
	Images             struct {
		Original               ImageData `json:"original"`
		FixedHeight            ImageData `json:"fixed_height"`
		FixedHeightStill       ImageData `json:"fixed_height_still"`
		FixedHeightDownsampled ImageData `json:"fixed_height_downsampled"`
		FixedWidth             ImageData `json:"fixed_width"`
		FixedwidthStill        ImageData `json:"fixed_width_still"`
		FixedwidthDownsampled  ImageData `json:"fixed_width_downsampled"`
	} `json:"images"`
}

type paginatedResults struct {
	Data       []*Gif `json:"data"`
	Pagination struct {
		TotalCount int `json:"total_count"`
	} `json:"pagination"`
}

type singleResult struct {
	Data *Gif `json:"data"`
}

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
   DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.

   EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			sendResponse(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!

func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)
	println("[RECEIVED] sendResponse:", command)

	// START SLACKBOT CUSTOM CODE
	// ===============================================================
	// TODO:
	//      1. Implement sendResponse for one or more of your custom Slackbot commands.
	//         You could call an external API here, or create your own string response. Anything goes!
	//      2. STRETCH: Write a goroutine that calls an external API based on the data received in this function.
	// ===============================================================
	// END SLACKBOT CUSTOM CODE\
	url := fmt.Sprintf("http://api.giphy.com/v1/gifs/search?api_key=4AZiEXqeJDmw6I1tzPAWobx790tH98f4&q=%s&limit=10", message)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// q := req.URL.Query()
	// q.Add("api_key", "4AZiEXqeJDmw6I1tzPAWobx790tH98f4")
	// q.Add("q", message)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	var data paginatedResults

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(data.Data[rand.Intn(9)].Images.FixedHeightDownsampled.URL, slackChannel))
}
