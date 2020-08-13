// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export uplink_stat_bucket
// uplink_stat_bucket returns information about a bucket.
//export MAKE_CONST=1,2
func uplink_stat_bucket(project *C.Uplink_Project, bucket_name *C.char) C.Uplink_BucketResult { //nolint:golint
	if project == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.StatBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.Uplink_BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_create_bucket
// uplink_create_bucket creates a new bucket.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
//export MAKE_CONST=1,2
func uplink_create_bucket(project *C.Uplink_Project, bucket_name *C.char) C.Uplink_BucketResult { //nolint:golint
	if project == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.CreateBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.Uplink_BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_ensure_bucket
// uplink_ensure_bucket creates a new bucket and ignores the error when it already exists.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
//export MAKE_CONST=1,2
func uplink_ensure_bucket(project *C.Uplink_Project, bucket_name *C.char) C.Uplink_BucketResult { //nolint:golint
	if project == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.EnsureBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.Uplink_BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_delete_bucket
// uplink_delete_bucket deletes a bucket.
//
// When bucket is not empty it returns ErrBucketNotEmpty.
//export MAKE_CONST=1,2
func uplink_delete_bucket(project *C.Uplink_Project, bucket_name *C.char) C.Uplink_BucketResult { //nolint:golint
	if project == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.Uplink_BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.Uplink_BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteBucket(proj.scope.ctx, C.GoString(bucket_name))
	return C.Uplink_BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(deleted),
	}
}

func mallocBucket(bucket *uplink.Bucket) *C.Uplink_Bucket {
	if bucket == nil {
		return nil
	}

	cbucket := (*C.Uplink_Bucket)(C.calloc(C.sizeof_Uplink_Bucket, 1))
	cbucket.name = C.CString(bucket.Name)
	cbucket.created = timeToUnix(bucket.Created)

	return cbucket
}

//export uplink_free_bucket_result
// uplink_free_bucket_result frees memory associated with the BucketResult.
func uplink_free_bucket_result(result C.Uplink_BucketResult) {
	uplink_free_error(result.error)
	uplink_free_bucket(result.bucket)
}

//export uplink_free_bucket
// uplink_free_bucket frees memory associated with the bucket.
//export MAKE_CONST=1
func uplink_free_bucket(bucket *C.Uplink_Bucket) {
	if bucket == nil {
		return
	}
	defer C.free(unsafe.Pointer(bucket))

	if bucket.name != nil {
		C.free(unsafe.Pointer(bucket.name))
	}
}
