package main

import (
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/asdine/storm"
	"io"
)

const (
	socks5Version   = uint8(5)
	noAcceptable    = uint8(255)
	userAuthVersion = uint8(1)
	authSuccess     = uint8(0)
	authFailure     = uint8(1)
)

type DatabaseAuthenticator struct {
	DB *storm.DB
}

func (a DatabaseAuthenticator) GetCode() uint8 {
	return socks5.UserPassAuth
}

func (a DatabaseAuthenticator) Authenticate(reader io.Reader, writer io.Writer) (*socks5.AuthContext, error) {
	// Tell the client to use user/pass auth
	if _, err := writer.Write([]byte{socks5Version, socks5.UserPassAuth}); err != nil {
		return nil, err
	}

	// Get the version and username length
	header := []byte{0, 0}
	if _, err := io.ReadAtLeast(reader, header, 2); err != nil {
		return nil, err
	}

	// Ensure we are compatible
	if header[0] != userAuthVersion {
		return nil, fmt.Errorf("Unsupported auth version: %v", header[0])
	}

	// Get the username
	usernameLen := int(header[1])
	username := make([]byte, usernameLen)
	if _, err := io.ReadAtLeast(reader, username, usernameLen); err != nil {
		return nil, err
	}

	// Get the password length
	if _, err := reader.Read(header[:1]); err != nil {
		return nil, err
	}

	// Get the password
	passwordLen := int(header[0])
	password := make([]byte, passwordLen)
	if _, err := io.ReadAtLeast(reader, password, passwordLen); err != nil {
		return nil, err
	}

	// Get users from db
	var user User
	err := a.DB.One("Username", string(username), &user)

	// User not found, auth failure
	if err == storm.ErrNotFound {
		if _, err := writer.Write([]byte{userAuthVersion, authFailure}); err != nil {
			return nil, err
		}
	}

	// If another error...
	if err != nil {
		return nil, err
	}

	// Verify the password
	if user.Password == string(password) {
		if _, err := writer.Write([]byte{userAuthVersion, authSuccess}); err != nil {
			return nil, err
		}
	} else {
		if _, err := writer.Write([]byte{userAuthVersion, authFailure}); err != nil {
			return nil, err
		}
		return nil, socks5.UserAuthFailed
	}

	// Done
	return &socks5.AuthContext{socks5.UserPassAuth, map[string]string{"Username": string(username)}}, nil
}
