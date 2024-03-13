package rpc_test

import (
	"context"
	"frost/internal/sigag"
	client "frost/internal/sigag/sigagclient"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var SigAgClient client.SigAgClient

var _ = Describe("Rpc", Ordered, func() {
	BeforeAll(func() {
		go func() {
			sigAg := sigag.New(sigag.Options{
				Logger: zap.NewExample(),
				Port:   "8080",
			})
			sigAg.StartSignatureAggregator(context.Background(), 10*time.Second)
		}()
		time.Sleep(5 * time.Second) // await until server is up

		SigAgClient = client.New("localhost", "8080", "/")
	})

	Context("While Party Client Interaction with SigAg Rpc", func() {
		It("should be able to check uptime", func() {
			isAlive, err := SigAgClient.CheckUptime()
			Expect(err).To(BeNil())
			Expect(isAlive).To(BeTrue())
		})
	})
})
