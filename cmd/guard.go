package cmd

import (
	"log"
	"time"

	"github.com/kardianos/service"
	"github.com/xdewx/git-ai-exporter/internal/git"
)

type guardProgram struct{}

func (p *guardProgram) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *guardProgram) Stop(s service.Service) error {
	return nil
}

func (p *guardProgram) run() {
	log.Println("git-ai-exporter guard started")
	interval := 5 * time.Second
	for {
		if !git.ProcessRunning("git-ai") {
			log.Println("git-ai daemon not running, starting...")
			git.RunCmd(".", "git-ai", "bg", "restart")
			time.Sleep(time.Second)
			if !git.ProcessRunning("git-ai") {
				git.RunCmd(".", "git-ai", "bg", "start")
			}
			if git.ProcessRunning("git-ai") {
				log.Println("git-ai daemon started")
			} else {
				log.Println("failed to start git-ai daemon, retrying later")
			}
		}
		time.Sleep(interval)
	}
}

func guardServiceName() string {
	return "git-ai-exporter-guard"
}
