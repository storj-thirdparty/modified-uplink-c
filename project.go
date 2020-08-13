// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// Project provides access to managing buckets.
type Project struct {
	scope
	*uplink.Project
}

//export uplink_open_project
// uplink_open_project opens project using access grant.
//export MAKE_CONST=1
func uplink_open_project(access *C.Uplink_Access) C.Uplink_ProjectResult {
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

	scope := rootScope("")
	config := uplink.Config{}

	proj, err := config.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.Uplink_ProjectResult{
			error: mallocError(err),
		}
	}

	return C.Uplink_ProjectResult{
		project: (*C.Uplink_Project)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

//export uplink_close_project
// uplink_close_project closes the project.
//export MAKE_CONST=1
func uplink_close_project(project *C.Uplink_Project) *C.Uplink_Error {
	if project == nil {
		return nil
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	proj.cancel()
	return mallocError(proj.Close())
}

//export uplink_free_project_result
// uplink_free_project_result frees any associated resources.
func uplink_free_project_result(result C.Uplink_ProjectResult) {
	uplink_free_error(result.error)
	freeProject(result.project)
}

// freeProject closes the project and frees any associated resources.
func freeProject(project *C.Uplink_Project) {
	if project == nil {
		return
	}
	defer C.free(unsafe.Pointer(project))
	defer universe.Del(project._handle)

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return
	}

	proj.cancel()
	// in case we haven't already closed the project
	_ = proj.Close()
	// TODO: log error when we didn't close manually and the close returns an error
}
