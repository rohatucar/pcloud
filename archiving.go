package pcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// savezip
// extractarchive
// extractarchiveprogress
// savezipprogress

// GetZip; https://docs.pcloud.com/methods/archiving/getzip.html
func (c *PCloudClient) GetZip(forceDownload int, filename string, timeOffset string, isEU bool) (io.Reader, error) {
	values := url.Values{
		"auth": {*c.Auth},
	}

	if forceDownload > 0 {
		values.Add("forcedownload", strconv.Itoa(forceDownload))
	}
	if filename != "" {
		values.Add("filename", filename)
	}
	if timeOffset != "" {
		values.Add("timeoffset", timeOffset)
	}

	resp, err := c.Client.Get(urlBuilder("getzip", values, isEU))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, checkResult(resp, err)
	}

	return resp.Body, nil
}

// GetZipLink; https://docs.pcloud.com/methods/archiving/getziplink.html
func (c *PCloudClient) GetZipLink(maxspeed int, forceDownload int, filename string, timeOffset string, isEU bool) ([]string, error) {
	var links []string

	values := url.Values{
		"auth": {*c.Auth},
	}

	if maxspeed > 0 {
		values.Add("maxspeed", strconv.Itoa(maxspeed))
	}
	if forceDownload > 0 {
		values.Add("forcedownload", strconv.Itoa(forceDownload))
	}
	if filename != "" {
		values.Add("filename", filename)
	}
	if timeOffset != "" {
		values.Add("timeoffset", timeOffset)
	}

	resp, err := c.Client.Get(urlBuilder("getziplink", values, isEU))
	if err != nil {
		return links, err
	}

	defer resp.Body.Close()
	result := struct {
		Result int      `json:"result"`
		Error  string   `json:"error"`
		Path   string   `json:"path"`
		Hosts  []string `json:"hosts"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return links, err
	}

	if result.Result != 0 {
		return links, errors.New(result.Error)
	}

	for _, host := range result.Hosts {
		links = append(links, fmt.Sprintf("https://%s%s", host, result.Path))
	}
	return links, nil
}
