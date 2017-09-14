package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	_ "github.com/lib/pq"
	util "github.com/woanware/goutil"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// ##### Types ###############################################################

// Encapsulates a Processor object and its properties
type Processor struct {
	id     int
	config *Config
	db     *runner.DB
}

//
type ImportTask struct {
	Domain string
	Host   string
	Data   string
}

// ##### Methods #############################################################

// Constructor/Initialiser for the Processor struct
func NewProcessor(id int, config *Config, db *runner.DB) *Processor {

	p := Processor{
		id:     id,
		config: config,
		db:     db,
	}

	return &p
}

// Process an individual set of host data
func (p *Processor) Process(it ImportTask) {

	p.correctXMLHeader(&it)

	tx, err := p.db.Begin()
	if err != nil {
		logger.Errorf("Error beginning transaction: %v", err)
		return
	}

	previousInstanceId := p.getPreviousInstanceId(it)
	currentInstanceId := p.getCurrentInstanceId(it)

	instance := p.insertInstance(it)
	if instance.Id < 1 {
		tx.Rollback()
		return
	}

	logger.Infof("Processing data: Domain: %s, Host: %s, Previous ID: %d, Current ID: %d, New ID: %d", it.Domain, it.Host, previousInstanceId, currentInstanceId, instance.Id)

	if currentInstanceId > -1 {

		// Move the current autorun data into the previous table
		if p.moveCurrentAutorunData(tx, currentInstanceId) == false {
			tx.Rollback()
			return
		}

		// Now delete the data from "current_autorun" where instance == previousInstanceId
		if p.deleteStaleAutoruns(currentInstanceId, true) == false {
			tx.Rollback()
			return
		}

		// Now delete the data from "current_autorun" where instance == currentInstanceId
		if p.deleteStaleAutoruns(previousInstanceId, false) == false {
			tx.Rollback()
			return
		}

		// Now delete the data from "instance" where id == previousInstanceId
		if p.deleteOldInstance(it.Domain, it.Host, currentInstanceId, instance.Id) == false {
			tx.Rollback()
			return
		}
	}

	p.insertAutoRunData(instance.Id, it)

	// Now commit the data as we have reached a good point
	tx.Commit()

	p.analyseData(instance, currentInstanceId)

	if len(config.ArchiveDir) > 0 {
		p.archiveData(it)
	}
}

// Autoruns.exe produces UTF16 XML but we convert to UTF8 in the client, but the XML header is wrong.
// Update the XML header so that the unmarshalling process knows the file is UTF8 now
func (p *Processor) correctXMLHeader(it *ImportTask) {
	it.Data = strings.Replace(string(it.Data), "utf-16", "UTF-8", -1)
}

// The server has a configurable **archive** option that stores a compressed archive of the autorun data.
// The data is stored in the directory specified by the **archive_dir** configuration value.
// The archives are stored in sub-directories in the form "domain-host". Each time a new set of data
// is received, the XML is compressed as a zip file, using the timestamp as the file name.
func (p *Processor) archiveData(it ImportTask) {

	domainHost := strings.ToLower(it.Domain) + "-" + strings.ToLower(it.Host)
	archiveDir := path.Join(config.ArchiveDir, domainHost)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	tempFile := p.writeZipArchive(domainHost, timestamp, it.Data)
	if len(tempFile) == 0 {
		return
	}

	defer func() {
		if util.DoesFileExist(tempFile) == true {
			err := os.Remove(tempFile)
			if err != nil {
				logger.Errorf("Error deleting temporary archive file: %v (%s)", err, tempFile)
			}
		}
	}()

	//
	md5, err := util.Md5File(tempFile)
	if err != nil {
		logger.Errorf("Error generate archive MD5: %v (%s)", err, tempFile)
		return
	}

	// Locate the last archive for the domain-host-user
	// If exists then read MD5, if MD5 is different then
	// move the archive we have just created, else delete it

	if util.DoesDirectoryExist(archiveDir) == true {

		// Retrieve a list of files for the particular "domain-host"
		files, _ := ioutil.ReadDir(archiveDir)
		// The file naming convention means that the last file is the last file written

		if len(files) == 0 {
			p.moveArchiveFileToArchiveDir(archiveDir, timestamp, tempFile, md5, false)
		} else {
			lastFile := files[len(files)-1].Name()

			oldMd5, err := util.ReadTextFromFile(path.Join(archiveDir, lastFile))
			if err != nil {
				logger.Errorf("Error reading old archive MD5 file: %v (%s)", err, path.Join(archiveDir, lastFile))
				return
			}

			//
			if md5 != oldMd5 {
				p.moveArchiveFileToArchiveDir(archiveDir, timestamp, tempFile, md5, false)
			}
		}
	} else {
		p.moveArchiveFileToArchiveDir(archiveDir, timestamp, tempFile, md5, true)
	}
}

// Moves the archive file to the "domain-host" specific directory
func (p *Processor) moveArchiveFileToArchiveDir(archiveDir string, fileName string, tempFile string, md5 string, makeDir bool) {

	if makeDir == true {
		// Create "domain-host-user" directory
		err := os.Mkdir(archiveDir, os.ModePerm)
		if err != nil {
			logger.Errorf("Error creating archive directory: %v (%s)", err, archiveDir)
			return
		}
	}

	// Move archive file to archive directory
	err := os.Rename(tempFile, path.Join(archiveDir, fileName+".zip"))
	if err != nil {
		logger.Errorf("Error moving file to archive directory: %v (%s)", err, archiveDir)
		return
	}

	// Write MD5 value to a file
	err = util.WriteBytesToFile(path.Join(archiveDir, fileName+".zip.md5"), []byte(md5), false)
	if err != nil {
		logger.Errorf("Error writing MD5 file: %v (%s)", err, archiveDir)
		return
	}
}

// Create zip file containing the XML autoruns output
func (p *Processor) writeZipArchive(domainHost string, timestamp string, data string) string {

	tf, err := ioutil.TempFile(config.TempDir, "arl-")
	if err != nil {
		logger.Errorf("Error creating temp file: %v", err)
		return ""
	}
	defer tf.Close()

	zw := zip.NewWriter(tf)
	defer zw.Close()

	zf, err := zw.Create(domainHost + ".xml")
	if err != nil {
		logger.Errorf("Error creating zip writer file: %v", err)
		return ""
	}

	_, err = zf.Write([]byte(data))
	if err != nil {
		logger.Errorf("Error creating zip writer file: %v", err)
		return ""
	}

	err = zw.Close()
	if err != nil {
		logger.Errorf("Error creating zip writer file: %v", err)
		return ""
	}

	return tf.Name()
}

// Moves the current "domain-host" combinations autorun data from
// the "current_autoruns" table to the "previous_autoruns" table
func (p *Processor) moveCurrentAutorunData(tx *runner.Tx, instanceId int64) bool {

	var data []*Autorun

	err := db.
		Select(`*`).
		From("current_autoruns").
		Where("instance = $1", instanceId).
		QueryStructs(&data)

	if err != nil {
		logger.Error(err)
		return false
	}

	//var a *Autorun
	for _, v := range data {
		err = tx.
			InsertInto("previous_autoruns").
			Columns("instance", "location", "item_name", "enabled", "profile", "launch_string", "description", "company",
				"signer", "version_number", "file_path", "file_name", "file_directory", "time", "sha256", "md5", "verified").
			Values(v.Instance, v.Location, v.ItemName, v.Enabled, v.Profile,
				v.LaunchString, v.Description, v.Company, v.Signer, v.VersionNumber,
				v.FilePath, v.FileName, v.FileDirectory, v.Time, v.Sha256, v.Md5, v.Verified).
			QueryStruct(&v)

		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") == false {
				logger.Errorf("Error moving Autorun: %v", err)
				return false
			}
		}
	}

	return true
}

// Deletes autorun data for a specific instance ID. Can work on either the current or previous autorun data tables
func (p *Processor) deleteStaleAutoruns(instanceId int64, currentTable bool) bool {

	tableName := "current_autoruns"
	if currentTable == false {
		tableName = "previous_autoruns"
	}

	_, err := p.db.
		DeleteFrom(tableName).
		Where("instance = $1", instanceId).
		Exec()

	if err != nil {
		logger.Errorf("Error deleting stale autoruns: %v (%d, %s)", err, instanceId, tableName)
		return false
	}

	return true
}

// Deletes any old instance records e.g. where the instance is not the current one nor the previous
func (p *Processor) deleteOldInstance(domain string, host string, previousId int64, currentId int64) bool {

	_, err := p.db.
		DeleteFrom("instance").
		Where("domain = $1 AND host = $2 AND (id <> $3) AND (id <> $4)", domain, host, previousId, currentId).
		Exec()

	if err != nil {
		logger.Errorf("Error deleting old instances: %v (Domain: %s, Host: %s, Previous ID: %d, Current ID: %d)", err, domain, host, previousId, currentId)
		return false
	}

	return true
}

// Inserts a new "instance" record into the database
func (p *Processor) insertInstance(it ImportTask) Instance {

	i := Instance{}
	err := p.db.
		InsertInto("instance").
		Columns("domain", "host", "timestamp").
		Values(it.Domain, it.Host, time.Now().UTC()).
		Returning("*").
		QueryStruct(&i)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Instance record: %v", err)
			return i
		}
	}

	return i
}

// Inserts a new "alert" record into the database
func (p *Processor) insertAlert(a *Autorun, i Instance, previousInstanceId int64) {

	alert := Alert{}
	err := p.db.
		InsertInto("alert").
		Columns("instance", "domain", "host", "timestamp", "autorun_id", "location",
			"item_name", "enabled", "profile", "launch_string", "description", "company",
			"signer", "version_number", "file_path", "file_name", "file_directory",
			"time", "sha256", "md5", "verified", "text", "linked").
		Values(i.Id, i.Domain, i.Host, i.Timestamp, a.Id, a.Location, a.ItemName, a.Enabled,
			a.Profile, a.LaunchString, a.Description, a.Company, a.Signer, a.VersionNumber,
			a.FilePath, a.FileName, a.FileDirectory, a.Time, a.Sha256, a.Md5, a.Verified, p.getAlertText(a),
			p.getLinkedAutoruns(previousInstanceId, a.FilePath, a.Sha256, a.Location, a.ItemName)).
		QueryStruct(&alert)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Alert record: %v", err)
		}
	}
}

//
func (p *Processor) getAlertText(a *Autorun) string {

	return fmt.Sprintf(
		`<strong>Item Name:</strong> %s<br>
		<strong>Location:</strong> %s<br>
		<strong>File Path:</strong> %s<br>
		<strong>Launch String:</strong> %s<br>
		<strong>Enabled:</strong> %t<br>
		<strong>Description:</strong> %s<br>
		<strong>Company:</strong> %s<br>
		<strong>Signer:</strong> %s<br>
		<strong>Verified:</strong> %s<br>
		<strong>Version:</strong> %s<br>
		<strong>Time:</strong> %s<br>
		<strong>SHA256:</strong> %s<br>
		<strong>MD5:</strong> %s<br>`,
		a.ItemName, a.Location, a.FilePath, a.LaunchString, a.Enabled,
		a.Description, a.Company, a.Signer, p.getVerifiedText(a.Verified), a.VersionNumber,
		a.Time.Format("15:04:05 02/01/2006"), a.Sha256, a.Md5)
}

//
func (p *Processor) getVerifiedText(verified int8) string {

	switch verified {
	case VERIFIED_FALSE:
		return "False"
	case VERIFIED_TRUE:
		return "True"
	case VERIFIED_MS:
		return "True (Microsoft)"
	default:
		return ""
	}
}

// Attempts to identify other autoruns that are linked either by file path or SHA256
func (p *Processor) getLinkedAutoruns(previousInstanceId int64, filePath string, sha256 string, location string, itemName string) string {

	var autoruns []*Autorun

	err := p.db.
		Select("*").
		From("previous_autoruns").
		Where("instance = $1 AND (file_path = $2 AND location = $4 AND item_name = $5) OR (sha256 = $3) AND sha256 <> ''", previousInstanceId, filePath, sha256, location, itemName).
		QueryStructs(&autoruns)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error retrieving linked autorun records: %v", err)
			return ""
		}
	}

	linked := make([]string, 0)
	for _, a := range autoruns {
		linked = append(linked, p.getAlertText(a))
	}

	return strings.Join(linked, "<br>")
}

// Parses the autorun XML data and inserts each entry as a record in the database
func (p *Processor) insertAutoRunData(instanceId int64, it ImportTask) {

	var autoruns XmlAutoruns
	err := xml.Unmarshal([]byte(it.Data), &autoruns)
	if err != nil {
		logger.Errorf("Error unmarshalling Autorun data: %v (%s)", err, "")
		return
	}

	var autorun *Autorun
	var a XmlAutorun
	var filePath string
	var fileName string
	var fileDirectory string
	var signer string

	tx, err := p.db.Begin()
	defer tx.AutoRollback()
	for _, a = range autoruns.Autoruns {

		filePath = util.RemoveQuotes(a.ImagePath)
		fileName, fileDirectory = util.SplitPath(filePath)

		autorun = new(Autorun)
		autorun.Instance = instanceId
		autorun.Company = util.RemoveQuotes(a.Company)
		autorun.Description = util.RemoveQuotes(a.Description)
		autorun.Enabled = util.ParseBoolean(a.Enabled, "Enabled", true)
		autorun.Location = util.RemoveQuotes(a.Location)
		autorun.FilePath = filePath
		autorun.FileName = fileName
		autorun.FileDirectory = fileDirectory
		autorun.ItemName = util.RemoveQuotes(a.ItemName)
		autorun.LaunchString = util.RemoveQuotes(a.LaunchString)
		autorun.Profile = util.RemoveQuotes(a.Profile)
		autorun.VersionNumber = util.RemoveQuotes(a.Version)
		autorun.Time = util.ParseTimestamp(LAYOUT_AUTORUNS, util.RemoveQuotes(a.Time))
		autorun.Sha256 = util.RemoveQuotes(a.Sha256)
		autorun.Md5 = util.RemoveQuotes(a.Md5)

		if strings.Contains(strings.ToLower(a.Signer), "(verified)") == true {

			// This is belt and braces e.g. if autoruns changes the format of its "signer" output
			signer = strings.Replace(util.RemoveQuotes(a.Signer), "(Verified) ", "", -1)
			signer = strings.Replace(signer, "(verified) ", "", -1)

			if IsFalsePositive(true, signer) == true {
				autorun.Verified = VERIFIED_MS
			} else {
				autorun.Verified = VERIFIED_TRUE
			}

		} else {
			autorun.Verified = VERIFIED_FALSE
			signer = strings.Replace(util.RemoveQuotes(a.Signer), "(Not verified) ", "", -1)
			// This is belt and braces e.g. if autoruns changes the format of its "signer" output
			signer = strings.Replace(signer, "(Not Verified) ", "", -1)
			signer = strings.Replace(signer, "(not verified) ", "", -1)
		}

		autorun.Signer = signer

		err = tx.
			InsertInto("current_autoruns").
			Columns("instance", "location", "item_name", "enabled", "profile", "launch_string", "description", "company",
				"signer", "version_number", "file_path", "file_name", "file_directory", "time", "sha256", "md5", "verified").
			Values(autorun.Instance, autorun.Location, autorun.ItemName, autorun.Enabled, autorun.Profile,
				autorun.LaunchString, autorun.Description, autorun.Company, autorun.Signer, autorun.VersionNumber,
				autorun.FilePath, autorun.FileName, autorun.FileDirectory, autorun.Time, autorun.Sha256, autorun.Md5, autorun.Verified).
			QueryStruct(&autorun)

		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") == false {
				logger.Errorf("Error inserting Autorun record: %v", err)
				return
			}
		}
	}
	tx.Commit()
}

// Retrieves the instance id of the previous domain/host/user autorun data
func (p *Processor) getPreviousInstanceId(it ImportTask) int64 {

	var i Instance

	err := p.db.
		Select("id").
		From("instance").
		Where("domain = $1 AND host = $2", it.Domain, it.Host).
		OrderBy("timestamp ASC").
		Limit(1).
		QueryStruct(&i)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error retrieving previous instance record: %v", err)
		} else {
			logger.Errorf("No previous instance record: (Domain: %s, Host: %s)", it.Domain, it.Host)
		}

		return -1
	}

	return i.Id
}

// Retrieves the instance id of the previous domain/host/user autorun data
func (p *Processor) getCurrentInstanceId(it ImportTask) int64 {

	var i Instance

	err := p.db.
		Select("id").
		From("instance").
		Where("domain = $1 AND host = $2", it.Domain, it.Host).
		OrderBy("timestamp DESC").
		Limit(1).
		QueryStruct(&i)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error retrieving current instance record: %v", err)
		} else {
			logger.Errorf("No current instance record: (Domain: %s, Host: %s)", it.Domain, it.Host)
		}

		return -1
	}

	return i.Id
}

// Retrieves a specific set of autorun data
func (p *Processor) getAutorunData(instanceId int64, currentTable bool) []*Autorun {

	tableName := "current_autoruns"
	if currentTable == false {
		tableName = "previous_autoruns"
	}

	var autoruns []*Autorun

	err := p.db.
		Select("*").
		From(tableName).
		Where("instance = $1", instanceId).
		QueryStructs(&autoruns)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error retrieving autorun data record: %v", err)
		} else {
			logger.Errorf("No instance autoruns data: %v (Instance: %d)", err, instanceId)
		}

		return []*Autorun{}
	}

	return autoruns
}

// Identifies new/changed/deleted autorun entries
func (p *Processor) analyseData(i Instance, previousInstanceId int64) {

	previous := p.getAutorunData(previousInstanceId, false)
	if len(previous) == 0 {
		logger.Errorf("No previous data to analyse: instance.id=%d", previousInstanceId)
		return
	}

	current := p.getAutorunData(i.Id, true)
	if len(current) == 0 {
		logger.Errorf("No current data to analyse: instance.id=%d", i.Id)
		return
	}

	var curr *Autorun
	var prev *Autorun
	located := false
	count := 0

	for _, curr = range current {
		located = false

		for _, prev = range previous {

			if strings.ToLower(curr.ItemName) == strings.ToLower(prev.ItemName) &&
				strings.ToLower(curr.Location) == strings.ToLower(prev.Location) &&
				strings.ToLower(curr.Profile) == strings.ToLower(prev.Profile) &&
				strings.ToLower(curr.FilePath) == strings.ToLower(prev.FilePath) &&
				strings.ToLower(curr.LaunchString) == strings.ToLower(prev.LaunchString) &&
				strings.ToLower(curr.Sha256) == strings.ToLower(prev.Sha256) {

				located = true
				break
			}
		}

		if located == false {
			located	= p.checkBasicProperties(curr, previous)
		}

		if located == false {
			count++
			p.insertAlert(curr, i, previousInstanceId)
		}
	}

	if count > 0 {
		logger.Infof("Added %d alerts: (Domain: %s, Host: %s)", count, i.Domain, i.Host)
	}
}

// Re-run the checks again but miss the Item Name as this seems to create a lot of false positives
func (p *Processor) checkBasicProperties (a *Autorun, previous []*Autorun) bool {

	located := false

	for _, prev := range previous {

		if strings.ToLower(a.Location) == strings.ToLower(prev.Location) &&
			strings.ToLower(a.Profile) == strings.ToLower(prev.Profile) &&
			strings.ToLower(a.FilePath) == strings.ToLower(prev.FilePath) &&
			strings.ToLower(a.LaunchString) == strings.ToLower(prev.LaunchString) &&
			strings.ToLower(a.Sha256) == strings.ToLower(prev.Sha256) {

			located = true
			break
		}
	}

	return located
}
