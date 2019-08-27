package types

type Dbconf struct {
	Type     string
	Host     string
	Port     uint16
	User     string
	Pass     string
	Database string
}

type Runtime struct {
	AppName        string
	Version        string
	BuildId        string
	Stage          string
	ConfigLocation string
	Dbconf         Dbconf
	Modloc         string
	Libloc         string
}
