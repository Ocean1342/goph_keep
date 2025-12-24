package application

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goph_keeper/cli/auth"
	"goph_keeper/cli/config"
	"goph_keeper/cli/crypter"
	"goph_keeper/cli/printer"
	"goph_keeper/cli/remote"
	"goph_keeper/cli/storage"
	"goph_keeper/cli/syncer"
	"os"
	"time"
)

type Application struct {
	Auth    *auth.Auth
	Crypter *crypter.Crypter
	Storage storage.Repository
	Syncer  *syncer.Syncer
	Printer *printer.Printer
	Remote  remote.Requester
	cancel  context.CancelFunc
}

var (
	App        Application
	ConfigData []byte
)

func Init(cmd *cobra.Command, _ []string) error {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	cfg, err := config.New(ConfigData)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := resty.New()
	crypt, _ := crypter.New()
	repository, inErr := storage.New()
	remoteService, _ := remote.New(client, cfg)
	authService, _ := auth.New(repository, remoteService)
	printerService := printer.New(cmd)
	if inErr != nil {
		err = inErr
	}
	syncerService := syncer.New(repository, remoteService)
	go syncerService.Sync(ctx, time.Duration(cfg.SyncInterval)*time.Second)
	App = Application{
		Auth:    authService,
		Crypter: crypt,
		Storage: repository,
		Syncer:  syncerService,
		Printer: printerService,
		Remote:  remoteService,
		cancel:  cancel,
	}
	if err != nil {
		return fmt.Errorf("init application failed. err: %v", err)
	}
	return nil
}

func ShutDown(cmd *cobra.Command, _ []string) {
	App.Printer.Println("Stopping all jobs... 3.2.1. Bye")
	time.Sleep(5 * time.Second)
	App.cancel()
	//время на завершение всех горутин
	time.Sleep(1 * time.Second)
}
