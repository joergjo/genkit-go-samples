package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	getUserMessage := func(ctx context.Context) (string, bool) {
		scanC := make(chan string, 1)
		errC := make(chan bool, 1)

		go func() {
			if scanner.Scan() {
				scanC <- scanner.Text()
			} else {
				errC <- true
			}
		}()

		select {
		case <-ctx.Done():
			return "", false
		case text := <-scanC:
			return text, true
		case <-errC:
			return "", false
		}
	}

	ctx := context.Background()
	agent := NewAgent(ctx, getUserMessage)
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	err := agent.Run(notifyCtx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
