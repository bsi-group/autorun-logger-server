package main

import (
	util "github.com/woanware/goutil"
	"io/ioutil"
	"time"
	"os"
	"path"
	"strings"
)

// ##### Constants #############################################################

const SQL_MD5 string = "SELECT DISTINCT md5 FROM current_autoruns WHERE md5 <> '' ORDER BY md5"
const SQL_SHA256 string = "SELECT DISTINCT sha256 FROM current_autoruns WHERE sha256 <> '' ORDER BY sha256"
const SQL_DOMAIN string = "SELECT DISTINCT domain FROM instance WHERE domain <> '' ORDER BY domain"
const SQL_HOST string = "SELECT DISTINCT host FROM instance WHERE host <> '' ORDER BY host"

const PREFIX_EXPORT_MD5 string = "export-md5-"
const PREFIX_EXPORT_SHA256 string = "export-sha256-"
const PREFIX_EXPORT_DOMAIN string = "export-domain-"
const PREFIX_EXPORT_HOST string = "export-host-"

// ##### Methods #############################################################

//
func exportAutorunData(sql string, typeName string, dataType int, prefix string) {

	rows, err := db.DB.Query(sql)
	if err != nil {
		logger.Errorf("Error querying for %s export: %v", typeName, err)
		return
	}
	defer rows.Close()

	tf, err:= ioutil.TempFile(config.TempDir, "arl-summary-")
	if err != nil {
		logger.Errorf("Error creating temp file for %s export: %v", typeName, err)
		return
	}
	defer tf.Close()

	defer func() {
		if util.DoesFileExist(path.Join(config.TempDir, tf.Name())) == true {
			err := os.Remove(path.Join(config.TempDir, tf.Name()))
			if err != nil {
				logger.Errorf("Error deleting temporary %s summary file: %v (%s)", typeName, err, tf.Name)
			}
		}
	}()

	data := ""
	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
			return
		}

		tf.WriteString(data)
		tf.WriteString("\n")
	}

	timestamp := time.Now().UTC()
	fileName := prefix + timestamp.Format(LAYOUT_DAILY_SUMMARY) + ".csv"

	// Move the file
	err = os.Rename(tf.Name(), path.Join(config.ExportDir, fileName))
	if err != nil {
		logger.Errorf("Error moving file to summary directory: %v (%s)", err, fileName)
		return
	}

	// Insert the summary record
	setExportRecord(dataType, fileName)
}

func exportInstanceData(sql string, typeName string, dataType int, prefix string) {

	rows, err := db.DB.Query(sql)
	if err != nil {
		logger.Errorf("Error querying for %s export: %v", typeName, err)
		return
	}
	defer rows.Close()

	tf, err:= ioutil.TempFile(config.TempDir, "arl-summary-")
	if err != nil {
		logger.Errorf("Error creating temp file for %s export: %v", typeName, err)
		return
	}
	defer tf.Close()

	defer func() {
		if util.DoesFileExist(path.Join(config.TempDir, tf.Name())) == true {
			err := os.Remove(path.Join(config.TempDir, tf.Name()))
			if err != nil {
				logger.Errorf("Error deleting temporary %s summary file: %v (%s)", typeName, err, tf.Name)
			}
		}
	}()

	data := ""
	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
			return
		}

		tf.WriteString(data)
		tf.WriteString("\n")
	}

	timestamp := time.Now().UTC()
	fileName := prefix + timestamp.Format(LAYOUT_DAILY_SUMMARY) + ".csv"


	// Move the file
	err = os.Rename(tf.Name(), path.Join(config.ExportDir, fileName))
	if err != nil {
		logger.Errorf("Error moving file to summary directory: %v (%s)", err, fileName)
		return
	}

	// Insert/Update the export record
	setExportRecord(dataType, fileName)
}

//
func setExportRecord(dataType int, fileName string) {

	var e Export

	err := db.
		Select("id, data_type, file_name, updated").
		From("export").
		Where("data_type = $1 and file_name = $2", dataType, fileName).
		OrderBy("id ASC").
		QueryStruct(&e)

	if e.Id > 0 {
		err = db.
			Update("export").
			Set("updated", time.Now().UTC()).
			Where("id = $1", e.Id).
			QueryScalar(&e.Updated)
		return
	}

	err = db.
		InsertInto("export").
		Columns("data_type", "file_name", "updated").
		Values(dataType, fileName, time.Now().UTC()).
		QueryStruct(&e)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting export record: %v (%s, %s, %s)", err, dataType, fileName)
		}
	}
}
