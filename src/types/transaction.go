package types

type Transaction struct {
	Price          int64
	Comment        string
	SlackUserId    string
	SlackChannelId string
	SlackTimestamp string
}
