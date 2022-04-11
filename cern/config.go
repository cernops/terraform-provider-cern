package cern

type CernConfig interface {
	GetLandbClient() (*LandbClient, error)
}

//Config stores the information needed by the provider to work
type config struct {
	LdapServer    string
	LandbEndpoint string
	LandbUsername string
	LandbPassword string
	TeigiClient   *Teigi
	RogerClient   *Roger
	CertMgrClient *CertMgr
}

func (c config) GetLandbClient() (*LandbClient, error) {
	// This LanDB client is initialised with a token that should be valid for
	// a few hours. A renovation mechanism has not been implemented yet.
	return NewLandbClient(c.LandbEndpoint, c.LandbUsername, c.LandbPassword)
}
