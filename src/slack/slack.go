package slack

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ToshihitoKon/slack-ctf-score-bot/src/constants"
	mydb "github.com/ToshihitoKon/slack-ctf-score-bot/src/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	api = slack.New(
		constants.SlackBotToken,
		//slack.OptionDebug(true),
		//slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(constants.SlackAppToken),
	)
	client = socketmode.New(
		api,
		//socketmode.OptionDebug(true),
		//socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
)

func Run() {
	go runner(api, client)

	fmt.Println("[INFO] ctf-score-bot")
	fmt.Println("[INFO] run websocket")
	client.Run()
}

func runner(api *slack.Client, client *socketmode.Client) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			fmt.Println("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			fmt.Println("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			fmt.Println("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				fmt.Printf("Ignored %+v\n", evt)
				continue
			}
			client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEventAction(eventsAPIEvent.InnerEvent)
			}
		}
	}
}

func innerEventAction(innerEvent slackevents.EventsAPIInnerEvent) {
	switch ev := innerEvent.Data.(type) {
	case *slackevents.ReactionAddedEvent:
		reactionAddedEvent, ok := innerEvent.Data.(*slackevents.ReactionAddedEvent)
		if !ok {
			log.Println("err: slackevents.ReactionAddedEvent")
			return
		}
		log.Println("ReactionAddedEvent: ", reactionAddedEvent)

	case *slackevents.MessageEvent:
		messageEvent, ok := innerEvent.Data.(*slackevents.MessageEvent)
		if !ok {
			log.Println("err: slackevents.MessageEvent")
			return
		}
		if messageEvent.BotID != "" {
			return
		}

		messageEventAction(messageEvent)

	case *slackevents.AppMentionEvent:
		_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("hey", false))
		if err != nil {
			log.Println("slackevents.AppMentionEvent: ", err.Error())
		}

	default:
		log.Println("default: ", innerEvent.Type)
	}
}

func messageEventAction(messageEvent *slackevents.MessageEvent) {
	var err error
	if !strings.HasPrefix(messageEvent.Text, "散財") {
		return
	}

	// 分割
	separate := regexp.MustCompile(`[\n| ]+`)
	text := separate.Split(messageEvent.Text, -1)
	if len(text) < 3 {
		_, _, err = api.PostMessage(messageEvent.Channel, slack.MsgOptionText("[金額] [メモ]", false))
		if err != nil {
			log.Println("messageEventAction: ", err.Error())
		}
		return
	}

	// 金額
	price, err := strconv.Atoi(text[1])
	if err != nil {
		log.Println("err: strconv.Atoi: ", err.Error())
	}

	// メモ
	comment := text[2]

	slackResponse := fmt.Sprint(
		"メモった",
		"\n額: ", strconv.Itoa(price),
		"\nメモ: ", comment,
	)

	err = mydb.InsertTransaction(price, comment, messageEvent.User, messageEvent.Channel, messageEvent.TimeStamp)
	if err != nil {
		slackResponse = fmt.Sprint("失敗したわ ", err.Error())
	}

	_, _, err = api.PostMessage(messageEvent.Channel, slack.MsgOptionText(slackResponse, false))
	if err != nil {
		log.Println("messageEventAction: ", err.Error())
	}
}
