package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatarURL = errors.New("chat: cannot get avatar URL")

type Avatar interface {
	// returns specified avatar URL.
	// something wrong, returns error,
	// especially ErrNoAvatarURL in the case if no avatar URL is got.
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar
type AuthAvatar struct{}
type GravatarAvatar struct{}
type FileSystemAvatar struct{}

var UseAuthAvatar AuthAvatar
var UseGravatarAvatar GravatarAvatar
var UseFileSystemAvatar FileSystemAvatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir(UploadDirectory); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return UploadDirectory + "/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
