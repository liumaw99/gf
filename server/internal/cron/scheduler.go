package cron

import "github.com/robfig/cron/v3"

// Scheduler 定时任务调度器
type Scheduler struct {
	c *cron.Cron
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		c: cron.New(),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	// TODO: 注册定时任务
	// 每日提醒：0 9 * * * （每天上午 9 点）
	// 晚间反思：0 21 * * * （每天晚上 9 点）
	// 周报生成：0 10 * * 1 （每周一上午 10 点）

	s.c.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.c.Stop()
}
