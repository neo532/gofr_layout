package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/wire"
	"github.com/neo532/gofr_layout/internal/config"
	fp "github.com/neo532/gokit/filepath"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/logger/slog"
	"github.com/neo532/gokit/logger/writer/lumberjack"
)

type Entry string

const (
	EntryApi      = "api"
	EntryScript   = "script"
	EntryConsumer = "consumer"
)

var (
	InitUnitTestSet = wire.NewSet(BootContext, InitConfig, InitLogger)
)

// func ConfBootstap() (bs *conf.Bootstrap) {
// 	rootPath := "."

// 	if pwd, err := os.Getwd(); err == nil {

// 		pn := strings.Split(reflect.TypeOf(pkg{}).PkgPath(), "/")
// 		pkgName := "/" + pn[len(pn)-2]

// 		tmp := strings.SplitN(pwd, pkgName, 2)
// 		if len(tmp) > 0 {
// 			rootPath = tmp[0] + pkgName
// 		}
// 	}

// 	var err error
// 	if bs, err = InitConfig(rootPath + "/configs/config.yaml"); err != nil {
// 		panic(err)
// 	}
// 	return bs
// }

// func ConfLogger(bs *conf.Bootstrap) klog.Logger {
// 	return InitLogger(bs, middleware.EntryTest, bs.General.Logger.FilenameTest)
// }

func InitConfig(c context.Context, path string) (cfg *config.Config, cleanup func(), err error) {

	// Resolve relative path to absolute by walking up to go.mod.
	if !filepath.IsAbs(path) {
		if dir, e := os.Getwd(); e == nil {
			for {
				if _, e := os.Stat(filepath.Join(dir, "go.mod")); e == nil {
					path = filepath.Join(dir, path)
					break
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					break
				}
				dir = parent
			}
		}
	}

	config.Cfg = &config.Config{}
	watcher := fp.New(path)
	cleanup = func() {
		watcher.Close()
	}
	if err = watcher.Watch(context.Background(), func(fileName string, data []byte) (err error) {
		config.Load(config.Cfg, fileName, data)
		return
	}); err != nil {
		return
	}

	ip, _ := os.Hostname()
	config.Cfg.General.Ip.Store(ip)

	cfg = config.Cfg
	return
}

func InitLogger(cfg *config.Config, entry Entry) (log logger.Logger, cleanup func(), err error) {
	l := slog.New(
		slog.WithWriter(
			lumberjack.New(
				lumberjack.WithFilename(
					strings.NewReplacer("{entry}", string(entry)).Replace(cfg.ConfigGeneral.General.Logger.Filename.Load().(string)),
				),
				lumberjack.WithMaxBackups(int(cfg.ConfigGeneral.General.Logger.MaxBackup.Load())),
				lumberjack.WithMaxSize(int(cfg.ConfigGeneral.General.Logger.MaxSize.Load())),
			),
		),
		slog.WithGlobalParam(
			"env", cfg.ConfigGeneral.General.Env.Load().(string),
			"name", cfg.ConfigGeneral.General.Name.Load().(string),
			"version", cfg.ConfigGeneral.General.Version.Load().(string),
			"entry", string(entry),
		),
		slog.WithLevel(cfg.ConfigGeneral.General.Logger.Level.Load().(string)),
		// slog.WithContextParam(cp, sp),
		// slog.WithHandler(NewPrettyHandler()),
	)
	if err = l.Error(); err != nil {
		return
	}
	log = logger.NewDefaultLogger(l)
	cleanup = func() {
		l.Close()
	}
	return
}

func BootContext() (c context.Context) {
	c = context.Background()
	return
}
