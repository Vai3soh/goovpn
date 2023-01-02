package app

import (
	"context"
	"embed"
	"os"
	"runtime"

	"github.com/Vai3soh/goovpn/entity"
	"github.com/Vai3soh/goovpn/internal/adapters/db/memory"
	"github.com/Vai3soh/goovpn/internal/cli"
	"github.com/Vai3soh/goovpn/internal/config"
	"github.com/Vai3soh/goovpn/internal/dns"
	"github.com/Vai3soh/goovpn/internal/fileextended"
	"github.com/Vai3soh/goovpn/internal/gui"
	"github.com/Vai3soh/goovpn/internal/parser"
	"github.com/Vai3soh/goovpn/internal/session"
	transport "github.com/Vai3soh/goovpn/internal/transport/openvpn"
	"github.com/Vai3soh/goovpn/internal/usecase"
	"github.com/Vai3soh/goovpn/internal/usecase/usecasedns"
	"github.com/Vai3soh/goovpn/internal/usecase/usecaseprofile"
	"github.com/Vai3soh/goovpn/pkg/boltdb"
	"github.com/Vai3soh/goovpn/pkg/logger"
	"github.com/Vai3soh/ovpncli"
	"github.com/wailsapp/wails/v2"
	lg "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func SaveHardCodedParams(
	bdb *boltdb.BoltDB, l *logger.Logger,
	bucketName string, values []entity.Message,
) {
	mustOpen(bdb, l)
	err := bdb.CreateBucket(bucketName)
	if err != nil {
		l.Fatalf("don't create storage [%s]\n", err)
	}

	mustOpen(bdb, l)
	err = bdb.StoreBulk(values)
	if err != nil {
		l.Fatalf("don't bulk store [%s]\n", err)
	}
}

func mustOpen(bdb *boltdb.BoltDB, l *logger.Logger) {
	err := bdb.ReOpen()
	if err != nil {
		l.Fatalf("don't open [%w]\n", err)
	}
}

func Run_wails(assets embed.FS, icon []byte, logLev, dbPath string) {

	l := logger.NewLogger(
		logger.WithLogTextFormatter(), logger.WithLogLevel(&logLev),
	)

	file := fileextended.NewFile()

	file.SetPath(dbPath)
	absDb, err := file.AbsolutePath()
	if err != nil {
		l.Fatal(err)
	}
	file.SetPath("")
	bdb, err := boltdb.NewBoltDB(*absDb)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	useSystemd := false
	configsPath := `~/ovpnconfigs`
	if runtime.GOOS == "windows" {
		userPath := os.Getenv(`USERPROFILE`)
		configsPath = userPath + `\` + `Desktop` + `\` + `ovpnconfigs` + `\`
	}
	countReccon := "3"

	mustOpen(bdb, l)
	bdb.SetNameBucket(`general_configure`)

	if !bdb.BucketIsCreate() {

		values := []entity.Message{
			{AtrId: "#ssl", Value: "0"},
			{AtrId: "#cmp", Value: "yes"},
		}
		SaveHardCodedParams(bdb, l, "ssl_cmp", values)

		values = []entity.Message{
			{AtrId: "tun_persist", Value: "checkbox"},
			{AtrId: "legacy_algo", Value: "checkbox"},
			{AtrId: "preferred_dc_algo", Value: "checkbox"},
		}

		if runtime.GOOS == `windows` {
			for _, e := range values {
				if e.AtrId == "legacy_algo" {
					e.Value = ""
				}
			}
		}

		SaveHardCodedParams(bdb, l, "general_openvpn_library", values)

		values = []entity.Message{
			{AtrId: "config_dir_path", Value: configsPath},
			{AtrId: "reconn_count", Value: countReccon},
		}
		SaveHardCodedParams(bdb, l, "general_configure", values)

		values = []entity.Message{
			{AtrId: "with_conn_timeout", Value: "0"},
		}
		SaveHardCodedParams(bdb, l, "other_options", values)
	}

	param := config.NewParamsDefault(bdb, l)
	param = param.GetParamIfStoreInDb()

	if runtime.GOOS == `windows` {
		param.SetLegacyAlgo(false)
	}
	sessOvpn := session.NewOpenvpnClient(
		context.TODO(),
		ovpncli.WithConnTimeout(param.ConnTimeout()),
		ovpncli.WithCompressionMode(param.CompressionMode()),
		ovpncli.WithDisableClientCert(param.DisableCert()),
		ovpncli.WithTunPersist(param.TunPersist()),
		ovpncli.WithLegacyAlgorithms(param.LegacyAlgorithms()),
		ovpncli.WithNonPreferredDCAlgorithms(param.NonPreferredDCAlgorithms()),
	)
	useSystemd = param.UseSystemd()

	gl := parser.NewConfig()

	ocl := sessOvpn.GetOverwriteClient()
	ocl.ClientAPI_OpenVPNClient = sessOvpn.Client

	memory := memory.NewDb(memory.WithMapMemory(
		make(map[string]entity.Profile)),
	)

	profileUseCase, err := usecaseprofile.NewProfileUseCase(
		file, file, file, gl, gl, gl, gl, gl, gl, memory,
	)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	command := cli.NewCmd()
	cmdResolver := cli.NewResolver(cli.WithCliResolver(command))

	dnsUseCase, err := usecasedns.NewDnsUseCase(
		command, command, cmdResolver, command, cmdResolver,
	)

	if err != nil {
		l.Fatalf("don't get constructor [%s]\n", err)
	}

	fileExt := fileextended.NewFile()

	vpnUseCase, _ := usecase.NewVpnUseCase(
		sessOvpn, sessOvpn, gl, gl, gl, gl, gl, gl,
		sessOvpn.GetOverwriteClient(), file, file, cmdResolver,
	)

	sys, err := dns.NewSystem(runtime.GOOS, dnsUseCase, dnsUseCase, useSystemd)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	names, err := dns.NewNames(dnsUseCase)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}
	names.SetGoos(*sys)

	var tr *transport.TransportOvpnClient
	if runtime.GOOS != "windows" {
		tr = transport.New(
			configsPath+"/", vpnUseCase,
			profileUseCase, names, l, bdb, sessOvpn,
		)
	} else {

		tr = transport.New(
			configsPath, vpnUseCase,
			profileUseCase, names, l, bdb, sessOvpn,
		)
	}
	gui := gui.NewGui(fileExt, sessOvpn, bdb, l, tr)

	err = wails.Run(&options.App{
		Title:             "Goovpn",
		Width:             450,
		Height:            850,
		MinWidth:          450,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          lg.DEBUG,
		OnStartup: func(ctx context.Context) {
			sessOvpn.SetContext(ctx)
			gui.Startup(ctx)
			tr.SetContext(ctx)
		},
		OnDomReady:       gui.DomReady,
		OnBeforeClose:    gui.BeforeClose,
		OnShutdown:       gui.Shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			gui, tr,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Goovpn",
				Message: "",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		l.Fatal(err)
	}

}
