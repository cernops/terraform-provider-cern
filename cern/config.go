package cern

//Config stores the information needed by the provider to work
type Config struct {
	LdapServer  string
	LandbClient *LandbClient
	TeigiClient *Teigi
}
