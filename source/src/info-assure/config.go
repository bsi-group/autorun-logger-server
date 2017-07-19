package main

// Holds the various objects/structs that are used in the system that don't warrant their own individual file

// ##### Structs #############################################################

// Stores the YAML config file data
type Config struct {
	DatabaseServer     string `yaml:"database_server"`
	DatabaseName       string `yaml:"database_name"`
	DatabaseUser       string `yaml:"database_user"`
	DatabasePassword   string `yaml:"database_password"`
	HttpIp             string `yaml:"http_ip"`
	HttpPort           int16  `yaml:"http_port"`
	ProcessorThreads   int    `yaml:"processor_threads"`
	Debug              bool   `yaml:"debug"`
	ArchiveDir         string `yaml:"archive_dir"`
	TempDir            string `yaml:"temp_dir"`
	ExportDir          string `yaml:"export_dir"`
	ServerPem          string `yaml:"server_pem"`
	ServerKey          string `yaml:"server_key"`
	MaxDataAgeDays     int    `yaml:"max_data_age_days"`
	AlertDurationHours int    `yaml:"alert_duration_hours"`
	SmtpServer         string `yaml:"smtp_server"`
	SmtpUser           string `yaml:"smtp_user"`
	SmtpPassword       string `yaml:"smtp_password"`
	SmtpPort           int    `yaml:"smtp_port"`
	SmtpSender         string `yaml:"smtp_sender"`
	SmtpFrom           string `yaml:"smtp_from"`
	SmtpReceiver       string `yaml:"smtp_receiver"`
	SmtpIsEncrypted    bool   `yaml:"smtp_is_encrypted"`
}
