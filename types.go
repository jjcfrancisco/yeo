package main

type databases struct {
	Dbs []database `json:"databases"`
}

type database struct {
	Name     string `json:"name"`
	IsLocal  bool   `json:"isLocal"`
	Database string `json:"database"`
	Port     string `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}
