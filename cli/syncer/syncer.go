package syncer

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"goph_keeper/cli/remote"
	"goph_keeper/cli/storage"
	"runtime"
	"sync"
	"time"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

// Syncer - флоу работы - команда
// демон читает из БД и отправляет данные на сервер (что делать, если данные попытаются изменить в момент отправки?)
// если отправить не удалось, то в очередь возвращается дата и ретрай таймс инкрементиться
type Syncer struct {
	store  storage.Repository
	sender remote.Requester
}

func New(store storage.Repository, sender remote.Requester) *Syncer {
	return &Syncer{
		store:  store,
		sender: sender,
	}
}

// Sync - демон синхронизации данных
func (s Syncer) Sync(ctx context.Context, syncInterval time.Duration) {
	t := time.NewTicker(syncInterval)
	defer t.Stop()
	s.run(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.run(ctx)
		}
	}
}

func (s Syncer) run(ctx context.Context) {
	//сходить в бд взять то, что ещё не отправлено на сервер
	items, err := s.store.GetNotSyncedItems(ctx)
	if err != nil {
		logrus.Errorf("get not-synced items failed. err %v", err)
		return
	}
	if len(items) == 0 {
		logrus.Info("no not-synced items")
		return
	}
	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		wg         sync.WaitGroup
	)

	for _, item := range items {
		if semErr := sem.Acquire(ctx, 1); semErr != nil {
			logrus.Errorf("Failed to acquire semaphore: %v", semErr)
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sem.Release(1)
			switch item.DataType {
			case string(storage.LogoPass):
				s.syncLogoPass(ctx, item)
			case string(storage.ArbitraryBinary):
				s.syncBin(ctx, item)
			default:
				logrus.Warnf("not support data type: %s data_id:%s", item.DataType, item.DataID)
			}
		}()
	}
	wg.Wait()
}

func (s Syncer) syncLogoPass(ctx context.Context, item storage.CatalogItem) {
	logo, pass, meta, err := s.store.GetLogoPass(ctx, item.UserID, item.DataName)
	if err != nil {
		logrus.Errorf("get logo-pass failed. err %v", err)
		return
	}
	err = s.sender.SendLogoPass(ctx, remote.StoreLogoPassRequest{
		UserLogin:    item.UserLogin,
		DataName:     item.DataName,
		DataLogin:    logo,
		Token:        item.UserToken,
		DataPassword: pass,
		DataMeta:     meta,
	})
	if err != nil {
		logrus.Errorf("send logo-pass failed. err %v", err)
		return
	}
	err = s.store.UpdateCatalogItem(ctx, item.UserID, item.DataID)
	if err != nil {
		logrus.Errorf("update catalog failed. err %v", err)
	}
}

func (s Syncer) syncBin(ctx context.Context, item storage.CatalogItem) {
	bin, meta, err := s.store.GetBin(ctx, item.UserID, item.DataName)
	if err != nil {
		logrus.Errorf("get bin failed. err %v", err)
		return
	}
	err = s.sender.SendBinData(ctx, remote.StoreBin{
		UserLogin: item.UserLogin,
		DataName:  item.DataName,
		Token:     item.UserToken,
		DataMeta:  meta,
		BinData:   bin,
	})
	if err != nil {
		logrus.Errorf("send bin failed. err %v", err)
		return
	}
	err = s.store.UpdateCatalogItem(ctx, item.UserID, item.DataID)
	if err != nil {
		logrus.Errorf("update catalog failed. err %v", err)
	}
}
