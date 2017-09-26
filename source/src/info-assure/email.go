package main

import (
	"io/ioutil"
	"os"
	"crypto/tls"
	"fmt"
	"bytes"
	"encoding/csv"
	"time"
	"strconv"
	"strings"
	util "github.com/woanware/goutil"
	gomail "gopkg.in/gomail.v2"
)

//
func sendAlerts() {

	logger.Info("Checking for alerts to send via SMTP")

	rows, err := db.DB.Queryx(fmt.Sprintf("SELECT * FROM alert WHERE id > %d", lastAlertId))
	if err != nil {
		logger.Errorf("Error performing the SELECT for alert summary: %v", err)
		return
	}
	defer rows.Close()

	// Initialise the CSV writer and output header line
	bAll := &bytes.Buffer{}
	cwAll := csv.NewWriter(bAll)
	cwAll.Write([]string{"Domain", "Host", "Profile", "Timestamp", "Location", "ItemName", "Enabled",
		"LaunchString", "Description", "Company", "Signer", "VersionNumber", "FilePath", "FileName",
		"FileDirectory", "SHA256", "MD5", "Verified"})

	bUnverified := &bytes.Buffer{}
	cwUnverified := csv.NewWriter(bUnverified)
	cwUnverified.Write([]string{"Domain", "Host", "Profile", "Timestamp", "Location", "ItemName", "Enabled",
		"LaunchString", "Description", "Company", "Signer", "VersionNumber", "FilePath", "FileName",
		"FileDirectory", "SHA256", "MD5", "Verified"})

	bMd5 := &bytes.Buffer{}
	cwMd5 := csv.NewWriter(bMd5)
	cwMd5.Write([]string{"MD5"})

	bSha256 := &bytes.Buffer{}
	cwSha256 := csv.NewWriter(bSha256)
	cwSha256.Write([]string{"SHA256"})

	a := Alert{}
	alertCount := 0
	for rows.Next() {
		err = rows.StructScan(&a)
		if err != nil {
			logger.Errorf("Error performing struct scan for alert summary: %v", err)
			continue
		}

		alertCount++

		// Ignore autoruns where the binary doesn't exist
		if strings.Contains(strings.ToLower(a.FilePath), "file not found") == true {
			continue
		}

		cwAll.Write([]string{a.Domain, a.Host, a.Profile, a.Time.Format(time.RFC3339), a.Location, a.ItemName,
			strconv.FormatBool(a.Enabled), a.LaunchString, a.Description, a.Company, a.Signer, a.VersionNumber,
			a.FilePath, a.FileName, a.FileDirectory, a.Sha256, a.Md5, getVerifiedString(a.Verified)})

		cwMd5.Write([]string{a.Md5})
		cwSha256.Write([]string{a.Sha256})

		if a.Verified == VERIFIED_FALSE {
			cwUnverified.Write([]string{a.Domain, a.Host, a.Profile, a.Time.Format(time.RFC3339), a.Location, a.ItemName,
				strconv.FormatBool(a.Enabled), a.LaunchString, a.Description, a.Company, a.Signer, a.VersionNumber,
				a.FilePath, a.FileName, a.FileDirectory, a.Sha256, a.Md5, getVerifiedString(a.Verified)})
		}

		lastAlertId = a.Id
	}

	if alertCount == 0 {
		logger.Info("No alerts found")
		return
	}

	logger.Infof("Found %d alerts to send via SMTP", alertCount)

	cwAll.Flush()
	cwUnverified.Flush()
	cwMd5.Flush()
	cwSha256.Flush()

	tmpFileAll := writeTempCsvFile("all", APP_NAME + "-all.csv", bAll)
	tmpFileUnverified := writeTempCsvFile("unverified", APP_NAME + "-unverified.csv", bUnverified)
	tmpFileMd5 := writeTempCsvFile("MD5", APP_NAME + "-md5.csv", bMd5)
	tmpFileSha256 :=  writeTempCsvFile("SHA256", APP_NAME + "-sha256.csv", bSha256)

	sendEmail(tmpFileAll, tmpFileUnverified, tmpFileMd5, tmpFileSha256)
}

//
func writeTempCsvFile(alertType string, fileName string, buffer *bytes.Buffer) string {

	tmpFile, err := ioutil.TempFile(config.TempDir, APP_NAME)
	if err != nil {
		logger.Errorf("Error creating CSV file for %s alerts email: %v", alertType, err)
		return ""
	}

	tmpFileName := tmpFile.Name()

	err = util.WriteBytesToFile(tmpFile.Name(), buffer.Bytes(), false)
	if err != nil {
		tmpFile.Close()
		logger.Errorf("Error writing CSV file for %s alerts email: %v", alertType, err)
		return ""
	}
	tmpFile.Close()

	err = os.Rename(tmpFileName, fileName)
	if err != nil {
		logger.Errorf("Error renaming temporary %s alerts file: %v", alertType, err)
		return ""
	}

	return fileName
}

//
func sendEmail(attachmentPathAll string, attachmentPathUnverified string, attachmentPathMd5 string, attachmentPathSha256 string) {

	defer os.Remove(attachmentPathAll)
	defer os.Remove(attachmentPathUnverified)
	defer os.Remove(attachmentPathMd5)
	defer os.Remove(attachmentPathSha256)

	dialer := gomail.NewDialer(config.SmtpServer, config.SmtpPort, config.SmtpUser, config.SmtpPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.SmtpSender)
	msg.SetHeader("To", config.SmtpReceiver)
	msg.SetHeader("Subject", EMAIL_ALERT_SUBJECT)
	msg.SetBody("text/html", "")

	if len(attachmentPathAll) > 0 {
		msg.Attach(attachmentPathAll)
	}

	if len(attachmentPathUnverified) > 0 {
		msg.Attach(attachmentPathUnverified)
	}

	if len(attachmentPathMd5) > 0 {
		msg.Attach(attachmentPathMd5)
	}

	if len(attachmentPathSha256) > 0 {
		msg.Attach(attachmentPathSha256)
	}

	if err := dialer.DialAndSend(msg); err != nil {
		logger.Errorf("Error sending alert email: %v", err)
	}
}
