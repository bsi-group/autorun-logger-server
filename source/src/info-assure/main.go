package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"github.com/voxelbrain/goptions"
	util "github.com/woanware/goutil"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
	"time"
)

// ##### Constants  ###########################################################

const APP_TITLE string = "AutoRun Logger Server"
const APP_NAME string = "arl-server"
const APP_VERSION string = "1.0.5"

// ##### Variables ###########################################################

var (
	logger    *logging.Logger
	config    *Config
	workQueue chan ImportTask
	db        *runner.DB
	cronner   *cron.Cron
)

// ##### Methods #############################################################

// Application entry point
func main() {
	fmt.Printf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION)

	initialiseLogging()

	opt := struct {
		ConfigFile string        `goptions:"-c, --config, description='Config file path'"`
		Help       goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{ // Default values
		ConfigFile: "./" + APP_NAME + ".config",
	}

	goptions.ParseAndFail(&opt)

	// Load the applications configuration such as database credentials
	config = loadConfig(opt.ConfigFile)

	initialiseDatabase()
	createProcessors()

	// Debug
	//p := NewProcessor(1, config, db)
	//i := Instance{}
	//i.Id = 277
	//p.analyseData(i, 273)
	//return

	cronner = cron.New()
	cronner.AddFunc("1 * * * * *", performHourlyTasks)
	cronner.AddFunc("* * * * 1 *", performDataPurge)
	//cronner.AddFunc("@hourly", performHourlyTasks)
	cronner.Start()

	var r *gin.Engine
	if config.Debug == false {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	} else {
		// DEBUG
		r = gin.Default()
	}

	r.GET("/", index)
	r.GET("/:domain/:host", receive)
	r.POST("/:domain/:host", receiveData)
	r.RunTLS(config.HttpIp+":"+fmt.Sprintf("%d", config.HttpPort), config.ServerPem, config.ServerKey)
}

// Initialises the database connection
func initialiseDatabase() {
	// create a normal database connection through database/sql
	tempDb, err := sql.Open("postgres",
		fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
			config.DatabaseServer, config.DatabaseName, config.DatabaseUser, config.DatabasePassword))

	if err != nil {
		logger.Fatalf("Unable to open database connection: %v", err)
		return
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(tempDb)

	// set to reasonable values for production
	tempDb.SetMaxIdleConns(4)
	tempDb.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 50 * time.Millisecond

	db = runner.NewDB(tempDb, "postgres")
}

// Initialise the channels for the cross process comms and then start the workers
func createProcessors() {
	processorCount := runtime.NumCPU()
	if config.ProcessorThreads > 0 {
		processorCount = config.ProcessorThreads
	}

	workQueue = make(chan ImportTask, 100)

	// Create the workers that perform the actual processing
	for i := 0; i < processorCount; i++ {
		logger.Infof("Initialising processor: %d", i+1)
		p := NewProcessor(i, config, db)
		go func(p *Processor) {
			for j := range workQueue {
				p.Process(j)
			}
		}(p)
	}
}

// Sets up the logging infrastructure e.g. Stdout and /var/log
func initialiseLogging() {
	// Setup the actual loggers
	logger = logging.MustGetLogger(APP_NAME)

	// Check that we have a "nca" sub directory in /var/log
	if _, err := os.Stat("/var/log/" + APP_NAME); os.IsNotExist(err) {
		logger.Fatal("The /var/log/" + APP_NAME + " directory does not exist")
	}

	// Check that we have permission to write to the /var/log/APP_NAME directory
	f, err := os.Create("/var/log/" + APP_NAME + "/test.txt")
	if err != nil {
		logger.Fatal("Unable to write to /var/log/" + APP_NAME)
	}

	// Clear up our tests
	os.Remove("/var/log/" + APP_NAME + "/test.txt")
	f.Close()

	// Define the /var/log file
	logFile, err := os.OpenFile("/var/log/"+APP_NAME+"/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("Error opening the log file: %v", err)
	}

	// Define the StdOut loggingDatabaser
	backendStdOut := logging.NewLogBackend(os.Stdout, "", 0)
	formatStdOut := logging.MustStringFormatter(
		"%{color}%{time:2006-01-02T15:04:05.999} %{color:reset} %{message}")
	formatterStdOut := logging.NewBackendFormatter(backendStdOut, formatStdOut)

	// Define the /var/log logging
	backendFile := logging.NewLogBackend(logFile, "", 0)
	formatFile := logging.MustStringFormatter(
		"%{time:2006-01-02T15:04:05.999} %{level:.4s} %{message}")
	formatterFile := logging.NewBackendFormatter(backendFile, formatFile)

	logging.SetBackend(formatterStdOut, formatterFile)
}

// Loads the applications config file contents (yaml) and marshals to a struct
func loadConfig(configPath string) *Config {
	c := new(Config)
	data, err := util.ReadTextFromFile(configPath)
	if err != nil {
		logger.Fatal("Error reading the config file: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		logger.Fatal("Error unmarshalling the config file: %v", err)
	}

	if len(c.DatabaseServer) == 0 {
		logger.Fatal("Database server not set in config file")
	}

	if len(c.DatabaseName) == 0 {
		logger.Fatal("Database name not set in config file")
	}

	if len(c.DatabaseUser) == 0 {
		logger.Fatal("Database user not set in config file")
	}

	if len(c.DatabasePassword) == 0 {
		logger.Fatal("Database password not set in config file")
	}

	if len(c.HttpIp) == 0 {
		logger.Fatal("HTTP IP not set in config file")
	}

	if len(c.ServerPem) == 0 {
		logger.Fatal("Server PEM file not set in config file")
	}

	if len(c.ServerKey) == 0 {
		logger.Fatal("Server key file not set in config file")
	}

	return c
}

//
func performHourlyTasks() {
	exportAutorunData(SQL_SHA256, "SHA256", EXPORT_TYPE_SHA256, PREFIX_EXPORT_SHA256)
	exportAutorunData(SQL_MD5, "MD5", EXPORT_TYPE_MD5, PREFIX_EXPORT_MD5)
	exportInstanceData(SQL_DOMAIN, "Domain", EXPORT_TYPE_DOMAIN, PREFIX_EXPORT_DOMAIN)
	exportInstanceData(SQL_HOST, "Host", EXPORT_TYPE_HOST, PREFIX_EXPORT_HOST)
}

//
func performDataPurge() {

	if config.MaxDataAgeDays == 0 || config.MaxDataAgeDays == -1 {
		return
	}

	// Use the config file value to determine what is classed as an old job
	staleTimestamp := time.Now().UTC().Add(-time.Duration(24*config.MaxDataAgeDays) * time.Hour)

	_, err := db.
		DeleteFrom("alerts").
		Where("timestamp < $1", staleTimestamp).
		Exec()

	if err != nil {
		logger.Errorf("Error deleting stale alerts: %v", err)
	}

	// Retrieve instance.Id's where timestamp is less
	var ids []int64

	err = db.
		Select(`id`).
		From("instance").
		Where("timestamp < $1", staleTimestamp).
		QueryStructs(&ids)

	for i := range ids {
		_, err = db.
			DeleteFrom("previous_autoruns").
			Where("instance = $1", i).
			Exec()

		logger.Errorf("Error deleting stale previous_autoruns records: %v (%d)", err, i)
	}

	for i := range ids {
		_, err = db.
			DeleteFrom("current_autoruns").
			Where("instance = $1", i).
			Exec()

		logger.Errorf("Error deleting stale current_autoruns records: %v (%d)", err, i)
	}
}
