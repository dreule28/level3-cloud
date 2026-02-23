package model

type Instance struct {
	ID		string `json:"id"`
	Status	string `json:"status"`
}

type ConnectionInfo struct {
	Host		string	`json:"host"`
	Port		int		`json:"port"`
	Database	string	`json:"database"`
	User		string	`json:"user"`
	Password	string	`json:"password"`
	Endpoint	string	`json:"endpoint"`
}

type InstanceDetails struct {
	ID			string			`json:"id"`
	Status		string			`json:"status"`
	Connection	*ConnectionInfo	`json:"connection,omitempty"`
}

type CreateInstanceRequest struct {
	ID			string	`json:"id"`
	Instances	int		`json:"instances"`
	StorageGi	int		`json:"storageGi"`
}
