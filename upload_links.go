package pcloud

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// CreateUploadLink; https://docs.pcloud.com/methods/upload_links/createuploadlink.html
func (c *PCloudClient) CreateUploadLink(path, comment string, folderID int, isEU bool) (link, code string, err error) {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	default:
		return "", "", errors.New("bad params")
	}
	switch {
	case comment != "":
		values.Add("comment", comment)
	default:
		return "", "", errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("createuploadlink", values, isEU))

	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	result := struct {
		Link   string `json:"link"`
		Code   string `json:"code"`
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	if result.Result > 0 {
		return "", "", errors.New(result.Error)
	}
	return result.Link, result.Code, nil
}
