package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/wire"
	"github.com/neo532/gofr_layout/internal/config"
	kitLog "github.com/neo532/gofr_layout/kit/log"
	fp "github.com/neo532/gokit/filepath"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/logger/slog"
	"github.com/neo532/gokit/logger/writer"
	"github.com/neo532/gokit/logger/writer/lumberjack"
	"github.com/neo532/gokit/logger/writer/stdout"
)

var (
	ProviderSet = wire.NewSet(BootContext, InitConfig, InitLogger)
	ConfigPath  = "./configs/"
	Entry       = "-"
	Version     string
)

const (
	EnvDev = "dev"
)

func InitConfig(c context.Context) (cfg *config.Config, cleanup func(), err error) {

	// Resolve relative path to absolute by walking up to go.mod.
	if !filepath.IsAbs(ConfigPath) {
		if dir, e := os.Getwd(); e == nil {
			for {
				if _, e := os.Stat(filepath.Join(dir, "go.mod")); e == nil {
					ConfigPath = filepath.Join(dir, ConfigPath)
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
	watcher := fp.New(ConfigPath)
	cleanup = func() {
		watcher.Close()
	}
	if err = watcher.Watch(c, func(fileName string, data []byte) (err error) {
		config.Load(config.Cfg, fileName, data)
		return
	}); err != nil {
		return
	}

	ip, _ := os.Hostname()
	config.Cfg.General.Ip.Store(ip)

	if Version != "" {
		config.Cfg.General.Version.Store(Version)
	}

	cfg = config.Cfg
	return
}

func InitLogger(cfg *config.Config) (log logger.Logger, cleanup func(), err error) {

	cfg.General.Entry.Store(Entry)

	var wit writer.Writer
	wit = lumberjack.New(
		lumberjack.WithFilename(
			strings.NewReplacer(
				"{entry}", cfg.General.Entry.Load().(string),
			).Replace(cfg.ConfigGeneral.General.Logger.Filename.Load().(string)),
		),
		lumberjack.WithMaxBackups(int(cfg.ConfigGeneral.General.Logger.MaxBackup.Load())),
		lumberjack.WithMaxSize(int(cfg.ConfigGeneral.General.Logger.MaxSize.Load())),
	)

	if cfg.General.Env.Load().(string) == EnvDev {
		wit = stdout.New()
	}

	opts := []slog.Option{
		slog.WithWriter(wit),
		slog.WithGlobalParam(
			"env", cfg.ConfigGeneral.General.Env.Load().(string),
			"ip", cfg.ConfigGeneral.General.Ip.Load().(string),
			"name", cfg.ConfigGeneral.General.Name.Load().(string),
			"version", cfg.ConfigGeneral.General.Version.Load().(string),
			"entry", string(cfg.General.Entry.Load().(string)),
		),
		slog.WithLevel(cfg.ConfigGeneral.General.Logger.Level.Load().(string)),
		slog.WithContextParam(kitLog.ProtocolFromContext),
	}

	if cfg.General.Env.Load().(string) == EnvDev {
		opts = append(opts, slog.WithHandler(slog.NewPrettyHandler()))
	}
	l := slog.New(opts...)
	if err = l.Error(); err != nil {
		return
	}

	return logger.NewDefaultLogger(l), func() { l.Close() }, nil
}

func BootContext() (c context.Context) {
	c = context.Background()
	return
}
