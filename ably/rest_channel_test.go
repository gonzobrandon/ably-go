package ably_test

import (
	"fmt"
	"testing"

	"github.com/ably/ably-go/ably"
	"github.com/ably/ably-go/ably/ablytest"
	"github.com/ably/ably-go/ably/proto"

	. "github.com/ably/ably-go/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/ably/ably-go/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("RestChannel", func() {
	const (
		event   = "sendMessage"
		message = "A message in a bottle"
	)

	Describe("publishing a message", func() {
		It("does not raise an error", func() {
			err := channel.Publish(event, message)
			Expect(err).NotTo(HaveOccurred())
		})

		It("is available in the history", func() {
			page, err := channel.History(nil)
			Expect(err).NotTo(HaveOccurred())

			messages := page.Messages()
			Expect(len(messages)).NotTo(Equal(0))
			Expect(messages[0].Name).To(Equal(event))
			Expect(messages[0].Data).To(Equal(message))
			Expect(messages[0].Encoding).To(Equal(proto.UTF8))
		})
	})

	Describe("History", func() {
		var historyRestChannel *ably.RestChannel

		BeforeEach(func() {
			historyRestChannel = client.Channel("channelhistory")

			for i := 0; i < 2; i++ {
				historyRestChannel.Publish("breakingnews", "Another Shark attack!!")
			}
		})

		It("returns a paginated result", func() {
			page1, err := historyRestChannel.History(&ably.PaginateParams{Limit: 1})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(page1.Messages())).To(Equal(1))
			Expect(len(page1.Items())).To(Equal(1))

			page2, err := page1.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(page2.Messages())).To(Equal(1))
			Expect(len(page2.Items())).To(Equal(1))
		})
	})

	Describe("PublishAll", func() {
		var encodingRestChannel *ably.RestChannel

		BeforeEach(func() {
			encodingRestChannel = client.Channel("this?is#an?encoding#channel")
		})

		It("allows to send multiple messages at once", func() {
			messages := []*proto.Message{
				{Name: "send", Data: "test data 1"},
				{Name: "send", Data: "test data 2"},
			}
			err := encodingRestChannel.PublishAll(messages)
			Expect(err).NotTo(HaveOccurred())

			page, err := encodingRestChannel.History(&ably.PaginateParams{Limit: 2})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(page.Messages())).To(Equal(2))
			Expect(len(page.Items())).To(Equal(2))
		})
	})
})

func TestRSL1f1(t *testing.T) {
	app, err := ablytest.NewSandbox(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	opts := app.Options()
	// RSL1f
	opts.UseTokenAuth = false
	client, err := ably.NewRestClient(opts)
	if err != nil {
		t.Fatal(err)
	}
	channel := client.Channels.Get("RSL1f", nil)
	id := "any_client_id"
	var msgs []*proto.Message
	size := 10
	for i := 0; i < size; i++ {
		msgs = append(msgs, &proto.Message{
			ClientID: id,
			Data:     fmt.Sprint(i),
		})
	}
	err = channel.PublishAll(msgs)
	if err != nil {
		t.Fatal(err)
	}
	res, err := channel.History(nil)
	if err != nil {
		t.Fatal(err)
	}
	m := res.Messages()
	n := len(m)
	if n != size {
		t.Errorf("expected %d messages got %d", size, n)
	}
	for _, v := range m {
		if v.ClientID != id {
			t.Errorf("expected clientId %s got %s data:%v", id, v.ClientID, v.Data)
		}
	}
}
