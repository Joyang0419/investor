package crawler

import (
	"time"
)

type SleepParam struct {
	MinSec int
	MaxSec int
}

type ICrawler[
	startParam,
	crawlParam,
	crawledData,
	errData any,
] interface {
	Start(
		startParam startParam,
		timeout time.Duration,
		randomSleep SleepParam,
	) error
	Validate(
		startParam startParam,
	) error
	Crawl(
		crawlParam crawlParam,
		timeout time.Duration,
		randomSleep SleepParam,
	) (
		response crawledData,
		err error,
	)
	HandleCrawledData(data crawledData) error
	HandleCrawlFailed(data errData) error
}
