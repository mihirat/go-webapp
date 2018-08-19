package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

const (
	testUrl       = "http://url-to-avatar/"
	testURLPrefix = "//www.gravatar.com/avatar/"
	testUserID    = "0bc83cb571cd1c50ba6f3e8a78ef1346"
	testFile      = "abc.jpg"
	testID        = "abc"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("in case of nonavailable value, AuthAvatar.GetAvatarURL must return ErrNoAvatarURL")
	}

	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("in case of available value, AuthAvatar.GetAvatarURL must not return error")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL must return correct URL")
		}
	}
}

func TestGravatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: testID}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("in case of available value, GravatarAvatar.GetAvatarURL must not return error")
	} else {
		if url != testURLPrefix+testID {
			t.Error("GravatarAvatar.GetAvatarURL returned wrong URL: " + url)
		}
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join(UploadDirectory, "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: testID}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("in case of available value, FilesystemAvatar.GetAvatarURL must not return error")
	}
	if url != (UploadDirectory + "/" + testFile) {
		t.Error("GravatarAvatar.GetAvatarURL returned wrong URL: " + url)
	}
}
