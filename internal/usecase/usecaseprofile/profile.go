package usecaseprofile

import (
	"fmt"
	"strings"

	"github.com/Vai3soh/goovpn/entity"
)

type ProfileUseCase struct {
	fileSetters      FileSetters
	fileToolsManager FileToolsManager
	fileReader       FileReader
	cfgSetters       ConfigSetters
	cfgCheck         ConfigChecker
	cfgMerg          ConfigMerger
	cfgRemover       ConfigRemover
	cfgTools         ConfigTools
	cfgBody          ConfigBody
	profileRepo      ProfileRepository
}

func NewProfileUseCase(
	fileSetters FileSetters,
	fileToolsManager FileToolsManager,
	fileReader FileReader,
	cfgSetters ConfigSetters,
	cfgCheck ConfigChecker,
	cfgMerg ConfigMerger,
	cfgRemover ConfigRemover,
	cfgTools ConfigTools,
	cfgBody ConfigBody,
	profileRepo ProfileRepository,

) (obj *ProfileUseCase, err error) {
	obj = &ProfileUseCase{
		fileSetters:      fileSetters,
		fileToolsManager: fileToolsManager,
		fileReader:       fileReader,
		cfgSetters:       cfgSetters,
		cfgCheck:         cfgCheck,
		cfgMerg:          cfgMerg,
		cfgRemover:       cfgRemover,
		cfgTools:         cfgTools,
		cfgBody:          cfgBody,
		profileRepo:      profileRepo,
	}
	return
}

func (p *ProfileUseCase) SearchFileAbsolutePath(file string) (*string, error) {
	p.fileSetters.SetPath(file)
	fileAbs, err := p.fileToolsManager.AbsolutePath()
	if err != nil {
		return nil, fmt.Errorf("not found abs path file [%w]", err)
	}
	return fileAbs, nil
}

func (p *ProfileUseCase) ReadFile() ([]byte, error) {
	body, err := p.fileReader.ReadFileAsByte()
	if err != nil {
		return nil, fmt.Errorf("don't read file: [%w]", err)
	}
	return body, nil
}

func (p *ProfileUseCase) GetMergedStringCfg(
	fileAbs, key, profileCurrent string,
) (*string, error) {
	p.fileSetters.SetPath(fileAbs)
	b, err := p.ReadFile()
	if err != nil {
		return nil, err
	}
	p.cfgSetters.SetBody(string(b))
	merged := p.cfgMerg.MergeCertsAndKeys(key)
	p.cfgSetters.SetBody(profileCurrent)
	p.cfgRemover.RemoveCommentLines()
	p.cfgRemover.RemoveCertsAndKeys()
	return &merged, nil
}

func (p *ProfileUseCase) OpenFileAndAddToConfig() (*string, error) {

	infile, err := p.fileToolsManager.FileOpen()
	if err != nil {
		return nil, fmt.Errorf("file open err: [%w]", err)
	}
	s := p.cfgTools.AddStringToConfig(infile)
	return &s, nil
}

func (p *ProfileUseCase) GetMapWithFileInConfig(profileBody string) map[string]string {
	p.cfgSetters.SetBody(profileBody)
	p.cfgRemover.RemoveNotCertsAndKeys()
	filesMap := p.cfgTools.SearchFilesPaths()
	return filesMap
}

func (p *ProfileUseCase) CleanConfigFromUseFiles(cfg string) string {
	p.cfgSetters.SetBody(cfg)
	p.cfgRemover.RemoveCertsAndKeys()
	return p.cfgBody.GetBody()
}

func (p *ProfileUseCase) SaveProfileWithBody(profile *entity.Profile) {
	p.profileRepo.Delete(*profile)
	p.profileRepo.Store(*profile)
}

func (p *ProfileUseCase) SaveProfileFromFile(cfg string) error {
	body, err := p.ReadFile()
	if err != nil {
		return err
	}
	p.fileSetters.SetBody(body)
	p.cfgSetters.SetBody(string(body))
	p.cfgRemover.RemoveSpaceLines()
	p.cfgRemover.RemoveCommentLines()
	p.cfgRemover.RemoveEmptyString()
	profile := entity.Profile{
		Body: p.cfgBody.GetBody(),
		Path: cfg,
	}
	p.profileRepo.Store(profile)
	return nil
}

func (p *ProfileUseCase) GetProfileFromCache(cfg string) entity.Profile {
	return p.profileRepo.Find(cfg)
}

func (p *ProfileUseCase) CheckFileExists(f string) bool {
	p.fileSetters.SetPath(f)
	return p.fileToolsManager.CheckFileExists()
}

func (p *ProfileUseCase) CheckUseCfgFile() bool {
	return p.cfgCheck.CheckConfigUseFiles()
}

func (p *ProfileUseCase) SaveProfileWithoutCfgFile(cfg string) error {
	err := p.SaveProfile(cfg)
	if err != nil {
		return fmt.Errorf("don't save profile: [%w]", err)
	}
	return nil
}

func (p *ProfileUseCase) SaveProfile(cfg string) error {
	profile := p.GetProfileFromCache(cfg)

	if profile.Body == "" {
		for _, f := range [2]string{cfg + ".ovpn", cfg + ".conf"} {
			if p.CheckFileExists(f) {
				break
			}
		}
		err := p.SaveProfileFromFile(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProfileUseCase) iterMapAndGetCurrentProfile(
	profileBody string, filesMap map[string]string,
) (*entity.Profile, error) {
	profileCurrent := profileBody
	for key, file := range filesMap {
		file = strings.Trim(file, `'`)
		file = strings.Trim(file, `"`)
		fileAbs, err := p.SearchFileAbsolutePath(file)
		if err != nil {
			return nil, err
		}
		merged, err := p.GetMergedStringCfg(*fileAbs, key, profileCurrent)
		if err != nil {
			return nil, err
		}
		profileCurrent = strings.Trim(profileCurrent, "\n") + *merged
	}
	pr := &entity.Profile{
		Body: profileCurrent,
	}
	return pr, nil
}

func (p *ProfileUseCase) getProfileFromDisk() (*entity.Profile, error) {
	profileBody, err := p.OpenFileAndAddToConfig()
	if err != nil {
		return nil, fmt.Errorf("add string to config err: [%w]", err)
	}
	filesMap := p.GetMapWithFileInConfig(*profileBody)
	profile, err := p.iterMapAndGetCurrentProfile(
		p.CleanConfigFromUseFiles(*profileBody), filesMap,
	)
	if err != nil {
		return nil, fmt.Errorf("don't iter map: [%w]", err)
	}
	return profile, nil
}

func (p *ProfileUseCase) SaveProfileWithCfgFile(path string) error {
	profile, err := p.getProfileFromDisk()
	if err != nil {
		return fmt.Errorf("don't get profile: [%w]", err)
	}
	profile.Path = path
	p.SaveProfileWithBody(profile)
	return nil
}
