package cronx

import "github.com/robfig/cron/v3"

type CronManager struct {
	cron *cron.Cron
}

func NewCronManager() *CronManager {
	return &CronManager{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (m *CronManager) AddFunc(spec string, cmd func()) {
	if _, err := m.cron.AddFunc(spec, cmd); err != nil {
		panic(err)
	}
}

func (m *CronManager) Start() {
	m.cron.Start()
}

func (m *CronManager) Stop() {
	m.cron.Stop()
}
