package pcloud

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// GetThumbLink; https://docs.pcloud.com/methods/thumbnails/getthumblink.html
func (c *PCloudClient) GetThumbLink(fileID int, path, size string, isEU bool) (string, error) {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case fileID >= 0:
		values.Add("fileid", strconv.Itoa(fileID))
	default:
		return "", errors.New("bad params")
	}

	if size != "" {
		values.Add("size", size)
	} else {
		return "", errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("getthumblink", values, isEU))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result := struct {
		Link   string `json:"link"`
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Result > 0 {
		return "", errors.New(result.Error)
	}
	return result.Link, nil
}

// SaveThumb; https://docs.pcloud.com/methods/thumbnails/savethumb.html
func (c *PCloudClient) SaveThumb(fileID int, path, size, toPath string, toFolderID int, toName string, isEU bool) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values.Add("fileid", strconv.Itoa(fileID))
	case path != "":
		values.Add("path", path)
	default:
		return errors.New("bad params")
	}

	if size != "" {
		values.Add("size", size)
	} else {
		return errors.New("bad params")
	}

	switch {
	case toFolderID > 0 && toName != "":
		values.Add("tofolderid", strconv.Itoa(toFolderID))
		values.Add("toname", toName)
	case toPath != "":
		values.Add("topath", toPath)
	default:
		return errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("savethumb", values, isEU))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Result > 0 {
		return errors.New(result.Error)
	}
	return nil
}
