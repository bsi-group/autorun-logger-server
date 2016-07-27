package main

import (

)

type XmlAutorun struct {
	Location 		string	`xml:"location"`
	ItemName 		string	`xml:"itemname"`
	Enabled 		string  `xml:"enabled"`
	Profile 		string  `xml:"profile"`
	LaunchString 	string	`xml:"launchstring"`
	Description 	string	`xml:"description"`
	Company 		string  `xml:"company"`
	Signer 			string  `xml:"signer"`
	Version 		string  `xml:"version"`
	ImagePath 		string  `xml:"imagepath"`
	Time 			string  `xml:"time"`
	Sha256			string  `xml:"sha256hash"`
	Md5				string  `xml:"md5hash"`
}

// Encapsulates the Autoruns parent (from CrowdResponse)
type XmlAutoruns struct {
	Autoruns	[]XmlAutorun	`xml:"item"`
}