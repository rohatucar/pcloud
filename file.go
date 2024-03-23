package pcloud

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

// uploadprogress
// downloadfile
// checksumfile

// DownloadFile; https://docs.pcloud.com/methods/file/downloadfile.html
func (c *PCloudClient) DownloadFile(urlStr string, path string, folderid int, target string, isEU bool) error {
	values := url.Values{
		"url":  {urlStr},
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderid >= 0:
		values.Add("folderid", strconv.Itoa(folderid))
	}

	if target != "" {
		values.Add("target", target)
	}

	return checkResult(c.Client.Get(urlBuilder("downloadfile", values, isEU)))
}

// UploadFile; https://docs.pcloud.com/methods/file/uploadfile.html
func (c *PCloudClient) UploadFile(reader io.Reader, path string, folderID int, filename string, noPartial int, progressHash string, renameIfExists int, isEU bool) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	default:
		return errors.New("bad params")
	}

	if filename == "" {
		return errors.New("bad params")
	}

	if noPartial > 0 {
		values.Add("nopartial", strconv.Itoa(noPartial))
	}
	if progressHash != "" {
		values.Add("progresshash", progressHash)
	}
	if renameIfExists > 0 {
		values.Add("renameifexists", strconv.Itoa(renameIfExists))
	}

	fw, err := w.CreateFormFile(filename, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", urlBuilder("uploadfile", values, isEU), &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return checkResult(c.Client.Do(req))
}

// CopyFile; https://docs.pcloud.com/methods/file/copyfile.html
func (c *PCloudClient) CopyFile(fileID int, path string, toFolderID int, toName string, toPath string, isEU bool) error {
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

	switch {
	case toFolderID > 0 && toName != "":
		values.Add("tofolderid", strconv.Itoa(toFolderID))
		values.Add("toname", toName)
	case toPath != "":
		values.Add("topath", toPath)
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("copyfile", values, isEU)))
}

// DeleteFile; https://docs.pcloud.com/methods/file/deletefile.html
func (c *PCloudClient) DeleteFile(fileID int, path string, isEU bool) error {
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

	return checkResult(c.Client.Get(urlBuilder("deletefile", values, isEU)))
}

// RenameFile; https://docs.pcloud.com/methods/file/renamefile.html
func (c *PCloudClient) RenameFile(fileID int, path string, toPath string, toFolderID int, toName string, isEU bool) error {
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

	switch {
	case toPath != "":
		values["topath"] = []string{toPath}
	case toFolderID > 0 && toName != "":
		values["toname"] = []string{toName}
		values["tofolderid"] = []string{strconv.Itoa(toFolderID)}
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("renamefile", values, isEU)))
}
