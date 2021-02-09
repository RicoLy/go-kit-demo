package client

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"go-kit-demo/demo03/v3/user_agent/pb"
	"os"
	"testing"
	"time"
)

func TestUserAgent_NewGrpcClient(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	client, err := NewUserAgentClient([]string{"127.0.0.1:2379"}, logger)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 6; i++ {
		time.Sleep(time.Second)
		userAgent, err := client.UserAgentClient()
		if err != nil {
			t.Error(err)
			return
		}
		ack, err := userAgent.Login(context.Background(), &pb.Login{
			Account: "ryan",
			Password: "123",
		})
		if err != nil {
			fmt.Println(err)
		} else {
			t.Log(ack.Token)
		}
	}
}
