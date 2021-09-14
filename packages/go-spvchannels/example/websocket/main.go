package main

import (
	"context"
	"encoding/json"
	"fmt"

	spv "github.com/libsv/go-spvchannels"
)

var channelid = "b1j-Vd94XrU9NJlnrtQPfkzOgFLxun5oWZLUhXfvnZk2cekKe4QY7YKh_hbXivAroApDtVn3pmmOo848R6BhAw"
var tok = "hqaDcOY-3svYqZv5RXID7AphKp9Bm8obQ_74K7mFLjcjq_Bw-Vwng6Q0q7PvJqhKawikfmd0Kr2OpYFKpFrKcg"

// PullUnreadMessages pull the notified unread messages, mark them as read
func PullUnreadMessages(t int, msg []byte, err error) error {
	// If notification error, then return the error
	if err != nil {
		return err
	}

	// Pull unread messages
	restClient := spv.NewClient(
		spv.WithBaseURL("localhost:5010"),
		spv.WithVersion("v1"),
		spv.WithUser("dev"),
		spv.WithPassword("dev"),
		spv.WithToken("tok"),
		spv.WithInsecure(),
	)

	r := spv.MessagesRequest{
		ChannelID: channelid,
		UnRead:    true,
	}

	unreadMsg, err := restClient.Messages(context.Background(), r)
	if err != nil {
		return fmt.Errorf("unable to read new messages : %w", err)
	}

	for _, msg := range *unreadMsg {
		msgSeq := msg.Sequence
		r2 := spv.MessageMarkRequest{
			ChannelID: channelid,
			Sequence:  msgSeq,
			Older:     true,
			Read:      true,
		}

		err := restClient.MessageMark(context.Background(), r2)
		if err != nil {
			return fmt.Errorf("unable mark message as read : %w", err)
		}
	}

	bReply, _ := json.MarshalIndent(unreadMsg, "", "    ")
	fmt.Println("\nNew unread messages ===================")
	fmt.Println(string(bReply))

	return nil
}

// This program run a websocket notification listener
// Anytime a new (unread) message is notified, it pull the new messages, mark them as read
func main() {

	ws := spv.NewWSClient(
		spv.WithBaseURL("localhost:5010"),
		spv.WithVersion("v1"),
		spv.WithChannelID(channelid),
		spv.WithToken(tok),
		spv.WithInsecure(),
		spv.WithWebsocketCallBack(PullUnreadMessages),
		spv.WithMaxNotified(10),
	)

	err := ws.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("Exit Success : total processed ", ws.NbNotified(), " messages")
}
