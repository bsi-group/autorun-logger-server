package main

// Holds the various objects/structs that are used in the system that don't warrant their own individual file

import (

)

// ##### Structs #############################################################

// Stores the YAML config file data
type Config struct {
	DatabaseServer		string 	`yaml:"database_server"`
	DatabaseName		string 	`yaml:"database_name"`
	DatabaseUser		string 	`yaml:"database_user"`
	DatabasePassword    string 	`yaml:"database_password"`
	HttpIp				string	`yaml:"http_ip"`
	HttpPort			int16	`yaml:"http_port"`
	ProcessorThreads	int		`yaml:"processor_threads"`
	Debug				bool	`yaml:"debug"`
	ArchiveDir			string	`yaml:"archive_dir"`
	TempDir				string	`yaml:"temp_dir"`
	ExportDir			string	`yaml:"export_dir"`
	ServerPem			string	`yaml:"server_pem"`
	ServerKey			string	`yaml:"server_key"`
}