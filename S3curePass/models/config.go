package models

type Config struct {
	Database   Database   `json:"database"`
	Token      Token      `json:"token"`
	Encryption Encryption `json:"encryption"`
}

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"Password"`
}

type Token struct {
	Secret string `json:"secret"`
	Header string `json:"header"`
}

type Encryption struct {
	Secret string `json:"secret"`
}
