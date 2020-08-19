package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omegaatt36/logaggd/cliflag"

	"github.com/urfave/cli"
)

// App is cli wrapper that do some common operation and creates signal handler.
type App struct {
	Flags []cli.Flag
	Main  func(ctx context.Context, c *cli.Context)
}

func (a *App) before(c *cli.Context) (err error) {
	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
			cli.ShowAppHelp(c)
			err = errors.New("init failed")
		}
	}()

	cliflag.Initialize(c)
	return
}

func (a *App) after(c *cli.Context) error {
	cliflag.Finalize()
	return nil
}

func (a *App) wrapMain(c *cli.Context) {

	// Signal handling.
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("\nReceives signal: %v\n", sig)
		cancel()
	}()

	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()

	a.Main(ctx, c)
	time.Sleep(3 * time.Second)
	log.Println("terminated")
}

// Run setups everything and runs Main.
func (a *App) Run() {
	app := cli.NewApp()
	app.Flags = a.Flags
	app.Flags = append(app.Flags, cliflag.Globals()...)
	app.Before = a.before
	app.After = a.after
	app.Action = a.wrapMain

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
