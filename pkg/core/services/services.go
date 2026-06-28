package services

import (
	"github.com/projectdiscovery/fastdialer/fastdialer"
	"github.com/xevonlive-dev/xevon/pkg/core/hosterrors"
	hostlimit "github.com/xevonlive-dev/xevon/pkg/core/ratelimit"
	"github.com/xevonlive-dev/xevon/pkg/dedup"
	"github.com/xevonlive-dev/xevon/pkg/notify"
	"github.com/xevonlive-dev/xevon/pkg/types"
)

// Services contains runtime services used across the application.
type Services struct {
	// Options contains CLI configuration
	Options *types.Options

	// HostLimiter limits concurrent requests per hostname
	HostLimiter *hostlimit.HostRateLimiter

	// HostErrors tracks host failures for circuit breaking
	HostErrors *hosterrors.Cache

	// Notifier sends notifications (Telegram, Discord, etc.)
	Notifier *notify.Manager

	// Dialer is the fastdialer instance for DNS resolution
	Dialer *fastdialer.Dialer

	// DedupManager manages deduplication for modules
	DedupManager *dedup.Manager
}
