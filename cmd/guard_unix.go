//go:build !windows

package cmd

import (
	"fmt"
	"log"

	"github.com/kardianos/service"
)

func runGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Run()
}

func doInstallGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	if err := s.Install(); err != nil {
		return fmt.Errorf("install guard service: %w", err)
	}
	if logger, err := s.Logger(nil); err == nil {
		logger.Info("Guard service installed")
	}
	fmt.Fprintf(log.Writer(), "Guard service installed: %s\n", s.String())
	if !service.Interactive() {
		if err := s.Start(); err != nil {
			fmt.Fprintf(log.Writer(), "Warning: failed to start guard service: %v\n", err)
		} else {
			fmt.Fprintln(log.Writer(), "Guard service started")
		}
	}
	return nil
}

func doUninstallGuard() error {
	s, err := newService()
	if err != nil {
		return err
	}
	if !service.Interactive() {
		s.Stop()
	}
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
			"UserService": true,
		},
	}
	return service.New(&guardProgram{}, cfg)
}
