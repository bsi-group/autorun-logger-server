package main

import (
	"github.com/gin-gonic/gin"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"fmt"
)

// Provides a HTTP response to a GET request for https://x.x.x.x:XX/
func index (c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION))
}

// Provides a HTTP response to a GET request for https://x.x.x.x:XX/domain/host/user
func receive (c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION))
}

// Receives the HTTP POST data from the clients
func receiveData (c *gin.Context) {

	// Get the URL values
	host := c.Param("host")
	domain := c.Param("domain")

	reader, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		logger.Errorf("Error creating GZIP reader: %v", err)
		c.String(http.StatusInternalServerError, "")
		return
	}
	defer reader.Close()

	// Read all of the HTTP POST body
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Errorf("Error reading POST body: %v", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	// Create an ImportTask struct and populate with the HTTP POST data
	it := ImportTask{domain, host, string(data)}
	go func() {
		workQueue <- it
	}()

	c.String(http.StatusOK, fmt.Sprintf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION))
	return
}