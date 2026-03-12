package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = ""
	port = 2222
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := initialModel(s)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func StartSSHServer() {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"), // It'll create this for you
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Printf("Could not start server: %v\n", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Starting SSH server on %s:%d\n", host, port)

	go func() {
		if err = s.ListenAndServe(); err != nil {
			fmt.Printf("Server failed: %v\n", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("Could not stop server: %v\n", err)
	}
}
