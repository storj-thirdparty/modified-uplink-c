// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// ObjectIterator is an iterator over objects.
type ObjectIterator struct {
	scope
	iterator *uplink.ObjectIterator

	initialError error
}

//export uplink_list_objects
// uplink_list_objects lists objects.
//export MAKE_CONST=1,2,3
func uplink_list_objects(project *C.Uplink_Project, bucket_name *C.char, options *C.Uplink_ListObjectsOptions) *C.Uplink_ObjectIterator { //nolint:golint
	if project == nil {
		return (*C.Uplink_ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	if bucket_name == nil {
		return (*C.Uplink_ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("bucket_name"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.Uplink_ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrInvalidHandle.New("project"),
		})))
	}

	opts := &uplink.ListObjectsOptions{}
	if options != nil {
		opts.Prefix = C.GoString(options.prefix)
		opts.Cursor = C.GoString(options.cursor)
		opts.Recursive = bool(options.recursive)

		opts.System = bool(options.system)
		opts.Custom = bool(options.custom)
	}

	scope := proj.scope.child()
	iterator := proj.ListObjects(scope.ctx, C.GoString(bucket_name), opts)

	return (*C.Uplink_ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

//export uplink_object_iterator_next
// upink_object_iterator_next prepares next Object for reading.
//
// It returns false if the end of the iteration is reached and there are no more objects, or if there is an error.
//export MAKE_CONST=1
func uplink_object_iterator_next(iterator *C.Uplink_ObjectIterator) C.bool {
	if iterator == nil {
		return C.bool(false)
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return C.bool(false)
	}
	if iter.initialError != nil {
		return C.bool(false)
	}

	return C.bool(iter.iterator.Next())
}

//export uplink_object_iterator_err
// uplink_object_iterator_err returns error, if one happened during iteration.
//export MAKE_CONST=1
func uplink_object_iterator_err(iterator *C.Uplink_ObjectIterator) *C.Uplink_Error {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}
	if iter.initialError != nil {
		return mallocError(iter.initialError)
	}

	return mallocError(iter.iterator.Err())
}

//export uplink_object_iterator_item
// uplink_object_iterator_item returns the current object in the iterator.
//export MAKE_CONST=1
func uplink_object_iterator_item(iterator *C.Uplink_ObjectIterator) *C.Uplink_Object {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return nil
	}

	return mallocObject(iter.iterator.Item())
}

//export uplink_free_object_iterator
// uplink_free_object_iterator frees memory associated with the ObjectIterator.
//export MAKE_CONST=1
func uplink_free_object_iterator(iterator *C.Uplink_ObjectIterator) {
	if iterator == nil {
		return
	}
	defer C.free(unsafe.Pointer(iterator))
	defer universe.Del(iterator._handle)

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if ok {
		if iter.scope.cancel != nil {
			iter.scope.cancel()
		}
	}
}
