// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// BucketIterator is an iterator over buckets.
type BucketIterator struct {
	scope
	iterator *uplink.BucketIterator

	initialError error
}

//export uplink_list_buckets
// uplink_list_buckets lists buckets.
func uplink_list_buckets(project *C.Uplink_Project, options *C.Uplink_ListBucketsOptions) *C.Uplink_BucketIterator {
	if project == nil {
		return (*C.Uplink_BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.Uplink_BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrInvalidHandle.New("project"),
		})))
	}

	opts := &uplink.ListBucketsOptions{}
	if options != nil {
		opts.Cursor = C.GoString(options.cursor)
	}

	scope := proj.scope.child()
	iterator := proj.ListBuckets(scope.ctx, opts)
	return (*C.Uplink_BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

//export uplink_bucket_iterator_next
// uplink_bucket_iterator_next prepares next Bucket for reading.
//
// It returns false if the end of the iteration is reached and there are no more buckets, or if there is an error.
func uplink_bucket_iterator_next(iterator *C.Uplink_BucketIterator) C.bool {
	if iterator == nil {
		return C.bool(false)
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return C.bool(false)
	}
	if iter.initialError != nil {
		return C.bool(false)
	}

	return C.bool(iter.iterator.Next())
}

//export uplink_bucket_iterator_err
// uplink_bucket_iterator_err returns error, if one happened during iteration.
func uplink_bucket_iterator_err(iterator *C.Uplink_BucketIterator) *C.Uplink_Error {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}
	if iter.initialError != nil {
		return mallocError(iter.initialError)
	}

	return mallocError(iter.iterator.Err())
}

//export uplink_bucket_iterator_item
// uplink_bucket_iterator_item returns the current bucket in the iterator.
func uplink_bucket_iterator_item(iterator *C.Uplink_BucketIterator) *C.Uplink_Bucket {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return nil
	}

	return mallocBucket(iter.iterator.Item())
}

//export uplink_free_bucket_iterator
// uplink_free_bucket_iterator frees memory associated with the BucketIterator.
func uplink_free_bucket_iterator(iterator *C.Uplink_BucketIterator) {
	if iterator == nil {
		return
	}
	defer C.free(unsafe.Pointer(iterator))
	defer universe.Del(iterator._handle)

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if ok {
		if iter.scope.cancel != nil {
			iter.scope.cancel()
		}
	}
}
