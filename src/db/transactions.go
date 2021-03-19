package db

type Transaction struct {
	Price          int
	Comment        string
	SlackUserId    string
	SlackChannelId string
	SlackTimestamp string
}

func InsertTransaction(price int, comment, slackUserId, slackChannelId, slackTimestamp string) error {
	return DB().Insert(&Transaction{
		Price:          price,
		Comment:        comment,
		SlackUserId:    slackUserId,
		SlackChannelId: slackChannelId,
		SlackTimestamp: slackTimestamp,
	})
}
