# Configuration

There are various configuration options for the server application. This document details the configuration options and the required values:

- database_server: The host name or IP address of the PostgreSQL database server
- database_name: Name of the database (arl)
- database_user:  Database server user name (postgres)
- database_password: Password for the database user
- http_ip: IP address of the interface for the HTTPS API. Use 0.0.0.0 to access on all interfaces
- http_port: Port for the HTTPS server. Use the **arl-setbind.sh** file to allow lower port access such as port 80
- debug: Show each HTTPS request in the logs (true/false)
- processor_threads. The number of processor threads. Use the value 0 to auto configure
- archive_dir: Directory used to store the unique archives of autorun data. If no value is supplied then the archiving is not performed
- temp_dir: Directory used as a temporary working area. The reason for using a directory other than /temp, is that programmatically moving files from one drive to another caused issues
- summary_dir: Directory used to store the automatically generated summary files
- server_pem: Full path to the server PEM file (server.pem)
- server_key: Full path to the server key file (server.key)
- max_data_age_days: Maximum age of alert, autorun data in days. Set to 0 or -1 to keep all data
- alert_duration_hours: The number of alerts between alerts. Must be between 0 & 23. A value of 0 will disable the alerts
