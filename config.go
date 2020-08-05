// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"time"

	"storj.io/uplink"
)

//export uplink_config_request_access_with_passphrase
// uplink_config_request_access_with_passphrase requests satellite for a new access grant using a passhprase.
func uplink_config_request_access_with_passphrase(config C.Uplink_Config, satellite_address, api_key, passphrase *C.char) C.Uplink_AccessResult { //nolint:golint
	if satellite_address == nil {
		return C.Uplink_AccessResult{
			error: mallocError(ErrNull.New("satellite_address")),
		}
	}
	if api_key == nil {
		return C.Uplink_AccessResult{
			error: mallocError(ErrNull.New("api_key")),
		}
	}
	if passphrase == nil {
		return C.Uplink_AccessResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ctx := context.Background()

	cfg := uplinkConfig(config)

	access, err := cfg.RequestAccessWithPassphrase(ctx, C.GoString(satellite_address), C.GoString(api_key), C.GoString(passphrase))
	if err != nil {
		return C.Uplink_AccessResult{
			error: mallocError(err),
		}
	}

	return C.Uplink_AccessResult{
		access: (*C.Uplink_Access)(mallocHandle(universe.Add(&Access{access}))),
	}
}

//export uplink_config_open_project
// uplink_config_open_project opens project using access grant.
func uplink_config_open_project(config C.Uplink_Config, access *C.Uplink_Access) C.Uplink_ProjectResult {
	if access == nil {
		return C.Uplink_ProjectResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.Uplink_ProjectResult{
			error: mallocError(ErrInvalidHandle.New("Access")),
		}
	}

	scope := rootScope(C.GoString(config.temp_directory))

	cfg := uplinkConfig(config)
	proj, err := cfg.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.Uplink_ProjectResult{
			error: mallocError(err),
		}
	}

	return C.Uplink_ProjectResult{
		project: (*C.Uplink_Project)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

func uplinkConfig(config C.Uplink_Config) uplink.Config {
	return uplink.Config{
		UserAgent:   C.GoString(config.user_agent),
		DialTimeout: time.Duration(config.dial_timeout_milliseconds) * time.Millisecond,
	}
}
