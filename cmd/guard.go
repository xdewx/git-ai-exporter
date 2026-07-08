package cmd

import (
	"fmt"
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

func runGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Run()
}

func stopGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Stop()
}

func restartGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	if err := s.Stop(); err != nil {
		log.Printf("stop guard before restart: %v", err)
	}
	return s.Start()
}

func doInstallGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	s.Stop()
	s.Uninstall()
	if err := s.Install(); err != nil {
		return fmt.Errorf("install guard service: %w", err)
	}
	if logger, err := s.Logger(nil); err == nil {
		logger.Info("Guard service installed")
	}
	fmt.Fprintf(log.Writer(), "Guard service installed: %s\n", s.String())
	if err := s.Start(); err != nil {
		fmt.Fprintf(log.Writer(), "Warning: failed to start guard service: %v\n", err)
	} else {
		fmt.Fprintln(log.Writer(), "Guard service started")
	}
	return nil
}

func doUninstallGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	s.Stop()
	if err := s.Uninstall(); err != nil {
		return fmt.Errorf("uninstall guard service: %w", err)
	}
	if logger, err := s.Logger(nil); err == nil {
		logger.Info("Guard service uninstalled")
	}
	fmt.Fprintln(log.Writer(), "Guard service uninstalled")
	return nil
}

func newService() (service.Service, error) {
	cfg := &service.Config{
		Name:        guardServiceName(),
		DisplayName: "Git AI Exporter Guard",
		Description: "Keeps git-ai daemon alive",
		Arguments:   []string{"--guard"},
		Option: service.KeyValue{
			"UserService":             true,
			"StartType":               "automatic",
			"OnFailure":               "restart",
			"OnFailureDelayDuration":  "5s",
			"OnFailureResetPeriod":    60,
		},
	}
	return service.New(&guardProgram{}, cfg)
}
