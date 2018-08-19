package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	useridFormValue = "userid"
	uploadFormFile  = "avatarFile"
	UploadDirectory = "avatars"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue(useridFormValue)
	file, header, err := req.FormFile(uploadFormFile)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join(UploadDirectory, userId+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "success!")
}
