// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export uplink_stat_object
// uplink_stat_object returns information about an object at the specific key.
//export MAKE_CONST=1,2,3
func uplink_stat_object(project *C.Uplink_Project, bucket_name, object_key *C.char) C.Uplink_ObjectResult { //nolint:golint
	if project == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	object, err := proj.StatObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.Uplink_ObjectResult{
		error:  mallocError(err),
		object: mallocObject(object),
	}
}

//export uplink_delete_object
// uplink_delete_object deletes an object.
//export MAKE_CONST=1,2,3
func uplink_delete_object(project *C.Uplink_Project, bucket_name, object_key *C.char) C.Uplink_ObjectResult { //nolint:golint
	if project == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_ObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.Uplink_ObjectResult{
		error:  mallocError(err),
		object: mallocObject(deleted),
	}
}

func mallocObject(object *uplink.Object) *C.Uplink_Object {
	if object == nil {
		return nil
	}

	cobject := (*C.Uplink_Object)(C.calloc(C.sizeof_Uplink_Object, 1))
	*cobject = objectToC(object)
	return cobject
}

func objectToC(object *uplink.Object) C.Uplink_Object {
	if object == nil {
		return C.Uplink_Object{}
	}
	return C.Uplink_Object{
		key:       C.CString(object.Key),
		is_prefix: C.bool(object.IsPrefix),
		system: C.Uplink_SystemMetadata{
			created:        timeToUnix(object.System.Created),
			expires:        timeToUnix(object.System.Expires),
			content_length: C.int64_t(object.System.ContentLength),
		},
		custom: customMetadataToC(object.Custom),
	}
}

//export uplink_free_object_result
// uplink_free_object_result frees memory associated with the ObjectResult.
func uplink_free_object_result(obj C.Uplink_ObjectResult) {
	uplink_free_error(obj.error)
	uplink_free_object(obj.object)
}

//export uplink_free_object
// uplink_free_object frees memory associated with the Object.
//export MAKE_CONST=1
func uplink_free_object(obj *C.Uplink_Object) {
	if obj == nil {
		return
	}
	defer C.free(unsafe.Pointer(obj))

	if obj.key != nil {
		C.free(unsafe.Pointer(obj.key))
	}

	freeSystemMetadata(&obj.system)
	freeCustomMetadataData(&obj.custom)
}

func freeSystemMetadata(system *C.Uplink_SystemMetadata) {
}
