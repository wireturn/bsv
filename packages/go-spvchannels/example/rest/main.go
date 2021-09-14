package main

import (
	"context"
	"encoding/json"
	"fmt"

	spv "github.com/libsv/go-spvchannels"
)

func main() {

	client := spv.NewClient(
		spv.WithBaseURL("localhost:5010"),
		spv.WithVersion("v1"),
		spv.WithUser("dev"),
		spv.WithPassword("dev"),
		spv.WithInsecure(),
	)

	r := spv.ChannelCreateRequest{
		AccountID:   "1",
		PublicRead:  true,
		PublicWrite: true,
		Sequenced:   true,
		Retention: struct {
			MinAgeDays int  "json:\"min_age_days\""
			MaxAgeDays int  "json:\"max_age_days\""
			AutoPrune  bool "json:\"auto_prune\""
		}{
			MinAgeDays: 0,
			MaxAgeDays: 99999,
			AutoPrune:  true,
		},
	}

	reply, err := client.ChannelCreate(context.Background(), r)
	if err != nil {
		panic("Problem with the request")
	}

	bReply, _ := json.MarshalIndent(reply, "", "    ")
	fmt.Println(string(bReply))
}
