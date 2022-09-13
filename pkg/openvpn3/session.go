/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package openvpn3

/*

#cgo CFLAGS: -I${SRCDIR}/bridge

#cgo LDFLAGS: -L${SRCDIR}/bridge
//main lib link
//TODO reuse GOOS somehow?
#cgo darwin,amd64 LDFLAGS: -lopenvpn3_darwin_amd64
#cgo ios,arm64 LDFLAGS: -lopenvpn3_ios_arm64
#cgo linux,!android,amd64 LDFLAGS: -lopenvpn3_linux_amd64
#cgo windows LDFLAGS: -lopenvpn3_windows_amd64
#cgo android,arm64 LDFLAGS: -lopenvpn3_android_arm64
#cgo android,amd64 LDFLAGS: -lopenvpn3_android_amd64
#cgo android,386 LDFLAGS: -lopenvpn3_android_x86
#cgo android,arm LDFLAGS: -lopenvpn3_android_armeabi-v7a
//TODO copied from openvpnv3 lib build tool - do we really need all of this?
#cgo darwin,amd64 LDFLAGS: -framework Security -framework CoreFoundation -framework SystemConfiguration -framework IOKit -framework ApplicationServices -mmacosx-version-min=10.8 -stdlib=libc++
//iOS frameworks
#cgo ios,arm64 LDFLAGS: -fobjc-arc -framework UIKit
#cgo windows LDFLAGS: -lws2_32 -liphlpapi

#include <library.h>

*/
import "C"
import (
	"context"
	"errors"
	"unsafe"

	"golang.org/x/sync/errgroup"
)

// Session represents the openvpn session
type Session struct {
	config          Config
	userCredentials *UserCredentials
	callbacks       interface{}
	tunnelSetup     TunnelSetup
	g               *errgroup.Group
	resError        error
	errorChan       chan error
	ptSession       unsafe.Pointer
}

// NewSession creates a new session given the callbacks
func NewSession(
	config Config, userCredentials UserCredentials, callbacks interface{},

) *Session {

	return &Session{
		config:          config,
		userCredentials: &userCredentials,
		callbacks:       callbacks,
		tunnelSetup:     &NoOpTunnelSetup{},
		g:               &errgroup.Group{},
		resError:        nil,
		ptSession:       nil,
		errorChan:       make(chan error),
	}
}

// MobileSessionCallbacks are the callbacks required for a mobile session
type MobileSessionCallbacks interface {
	EventConsumer
	Logger
	StatsConsumer
}

// NewMobileSession creates a new mobile session provided the required callbacks and tunnel setup
func NewMobileSession(
	config Config, userCredentials UserCredentials,
	callbacks MobileSessionCallbacks, tunSetup TunnelSetup,
) *Session {

	return &Session{
		config:          config,
		userCredentials: &userCredentials,
		callbacks:       callbacks,
		tunnelSetup:     tunSetup,
		g:               &errgroup.Group{},
		resError:        nil,
		ptSession:       nil,
		errorChan:       make(chan error),
	}
}

// ErrInitFailed is the error we return when openvpn3 initiation fails
var ErrInitFailed = errors.New("openvpn3 init failed")

// ErrConnectFailed is the error we return when openvpn3 fails to connect
var ErrConnectFailed = errors.New("openvpn3 connect failed")

func (session *Session) controlCancelContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if session.ptSession != nil {
				C.stop_session(session.ptSession)
				return
			}
		default:
			continue
		}
	}
}

func (session *Session) getConfig() (*C.struct___3, func()) {
	cConfig, cConfigUnregister := session.config.toPtr()
	return &cConfig, cConfigUnregister
}

func (session *Session) getCread() (*C.struct___4, func()) {
	cCredentials, cCredentialsUnregister := session.userCredentials.toPtr()
	return &cCredentials, cCredentialsUnregister
}

func (session *Session) getCallbackDelegate() (expCallbacks, func()) {
	callbacksDelegate, removeCallback := registerCallbackDelegate(session.callbacks)
	return callbacksDelegate, removeCallback
}

func (session *Session) getTunnelBuilderCallbacks() (*C.struct___8, func()) {
	tunBuilderCallbacks, removeTunCallbacks := registerTunnelSetupDelegate(session.tunnelSetup)
	return &tunBuilderCallbacks, removeTunCallbacks
}

func (session *Session) getSessionPtr(
	cConfig *C.struct___3, cCredentials *C.struct___4,
	callbacksDelegate expCallbacks, tunBuilderCallbacks *C.struct___8,
) unsafe.Pointer {
	sessionPtr, _ := C.new_session(
		*cConfig,
		*cCredentials,
		C.callbacks_delegate(callbacksDelegate),
		C.tun_builder_callbacks(*tunBuilderCallbacks),
	)
	return sessionPtr
}

func (session *Session) Start(ctx context.Context) {

	cConfig, cConfigUnregister := session.getConfig()
	cCredentials, cCredentialsUnregister := session.getCread()
	callbacksDelegate, removeCallback := session.getCallbackDelegate()
	tunBuilderCallbacks, removeTunCallbacks := session.getTunnelBuilderCallbacks()

	sessionPtr := session.getSessionPtr(
		cConfig, cCredentials,
		callbacksDelegate, tunBuilderCallbacks,
	)

	session.g.Go(func() error {

		defer cConfigUnregister()
		defer cCredentialsUnregister()
		defer removeCallback()
		defer removeTunCallbacks()

		if sessionPtr == nil {
			session.resError = ErrInitFailed
			return session.resError
		}

		session.ptSession = sessionPtr
		go session.controlCancelContext(ctx)

		err := session.run()
		if err != nil {
			return err
		}
		C.cleanup_session(sessionPtr)
		session.ptSession = nil
		return nil
	})

	go session.getErrorFromSession()
}

func (session *Session) getErrorFromSession() {
	session.errorChan <- session.g.Wait()
}

func (session *Session) CallbackError() error {
	for {
		select {
		case err := <-session.errorChan:
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
}

func (session *Session) run() error {
	if session.ptSession != nil {
		res, _ := C.start_session(session.ptSession)
		if res != 0 {
			session.resError = ErrConnectFailed
		}
	}
	return session.resError
}

func (session *Session) Reconnect(seconds int) {
	if session.ptSession != nil {
		C.reconnect_session(session.ptSession, C.int(seconds))
	}
}

func (session *Session) Stop() {
	if session.ptSession != nil {
		C.stop_session(session.ptSession)
	}
}
