package usecase

type VpnUseCase struct {
	LogRepo     LoggerInteractor
	GlueRepo    GlueConfigInteractor
	CloseRepo   CloseAppInteractor
	SessionRepo SessionInteractor
	UiRepo      UiInteractor
	TrayRepo    SysTrayInteractor
	FileRepo    FileInteractor
	CmdRepo     CommandInteractor
	DnsRepo     DnsInteractor
	MemoryRepo  MemoryInteractor
}

func New(
	l LoggerInteractor,
	g GlueConfigInteractor,
	cl CloseAppInteractor,
	sess SessionInteractor,
	u UiInteractor,
	st SysTrayInteractor,
	f FileInteractor,
	c CommandInteractor,
	d DnsInteractor,
	m MemoryInteractor,

) *VpnUseCase {

	return &VpnUseCase{
		LogRepo:     l,
		GlueRepo:    g,
		CloseRepo:   cl,
		SessionRepo: sess,
		UiRepo:      u,
		TrayRepo:    st,
		FileRepo:    f,
		CmdRepo:     c,
		DnsRepo:     d,
		MemoryRepo:  m,
	}
}
