package rpc_test

import (
	"context"
	"fmt"
	"frost/internal/sigag"
	client "frost/internal/sigag/sigagclient"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var SigAgClient client.SigAgClient

var _ = Describe("Rpc", Ordered, func() {
	BeforeAll(func() {
		go func() {
			sigAg := sigag.New(sigag.Options{
				Logger: logrus.New(),
				Port:   "8080",
			})
			sigAg.StartSignatureAggregator(context.Background(), 10*time.Second)
		}()
		time.Sleep(5 * time.Second) // await until server is up

		SigAgClient = client.New(fmt.Sprintf("http://%s:%s%s", "localhost", "8080", "/"))
	})

	Context("While Party Client Interaction with SigAg Rpc", func() {
		It("should be able to check uptime", func() {
			isAlive, err := SigAgClient.CheckUptime()
			Expect(err).To(BeNil())
			Expect(isAlive).To(BeTrue())
		})

		It("should be able to register", func() {
			err := SigAgClient.Register("1", "127.0.0.1", "3", "4")
			Expect(err).To(BeNil())

			err = SigAgClient.Register("2", "127.0.0.1", "5", "6")
			Expect(err).To(BeNil())
		})

		It("should not be able to register invalid participant", func() {
			err := SigAgClient.Register("1", "127.0.0.1", "3", "4") // same user
			Expect(err).ToNot(BeNil())

			err = SigAgClient.Register("3", "127.1", "3", "4") // invalid ip
			Expect(err).ToNot(BeNil())
		})

		It("should be able to get participant list", func() {
			participants, err := SigAgClient.GetParticipants()
			Expect(err).To(BeNil())
			Expect(participants).ToNot(BeNil())
			Expect(len(participants)).To(Equal(2))

			Expect(participants["1"]).To(Equal("127.0.0.1:3:4"))
			Expect(participants["2"]).To(Equal("127.0.0.1:5:6"))
		})
	})
})
