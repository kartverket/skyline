package sender

import (
	"context"
	"github.com/kartverket/skyline/pkg/email"
)

type Sender interface {
	Send(ctx context.Context, email *email.SkylineEmail) error
}
