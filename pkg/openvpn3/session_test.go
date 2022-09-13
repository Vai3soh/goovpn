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

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fmtLogger struct {
}

func (l *fmtLogger) Log(text string) {
	fmt.Println(text)
}

func TestSessionStartStopDoesNotPanic(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	session := NewSession(Config{}, UserCredentials{}, &fmtLogger{})
	session.Start(ctx)
	session.Stop()
}

func TestSessionInitFailsForInvalidProfile(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	session := NewSession(Config{}, UserCredentials{}, &fmtLogger{})
	session.Start(ctx)
	err := session.CallbackError()
	assert.Equal(t, ErrInitFailed, err)
}

func TestSessionConnectFailsForInvalidRemote(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	cfg := NewConfig(`remote localhost 1111`)
	session := NewSession(cfg, UserCredentials{}, &fmtLogger{})
	session.Start(ctx)
	err := session.CallbackError()
	assert.Equal(t, ErrConnectFailed, err)
}
