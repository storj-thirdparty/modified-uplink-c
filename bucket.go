// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export stat_bucket
// stat_bucket returns information about a bucket.
func stat_bucket(project *C.Project, bucket_name *C.char) C.BucketResult {
	if bucket_name == nil {
		return C.BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.StatBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export create_bucket
// create_bucket creates a new bucket.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func create_bucket(project *C.Project, bucket_name *C.char) C.BucketResult {
	if bucket_name == nil {
		return C.BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.CreateBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export ensure_bucket
// ensure_bucket creates a new bucket and ignores the error when it already exists.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func ensure_bucket(project *C.Project, bucket_name *C.char) C.BucketResult {
	if bucket_name == nil {
		return C.BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.EnsureBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export delete_bucket
// delete_bucket deletes a bucket.
//
// When bucket is not empty it returns ErrBucketNotEmpty.
func delete_bucket(project *C.Project, bucket_name *C.char) C.BucketResult {
	if bucket_name == nil {
		return C.BucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.BucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteBucket(proj.scope.ctx, C.GoString(bucket_name))
	return C.BucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(deleted),
	}
}

func mallocBucket(bucket *uplink.Bucket) *C.Bucket {
	if bucket == nil {
		return nil
	}

	cbucket := (*C.Bucket)(C.malloc(C.sizeof_Bucket))
	cbucket.name = C.CString(bucket.Name)
	cbucket.created = timeToUnix(bucket.Created)

	return cbucket
}

//export free_bucket_result
// free_bucket_result frees memory associated with the BucketResult.
func free_bucket_result(result C.BucketResult) {
	free_error(result.error)
	free_bucket(result.bucket)
}

//export free_bucket
// free_bucket frees memory associated with the bucket.
func free_bucket(bucket *C.Bucket) {
	if bucket == nil {
		return
	}

	if bucket.name != nil {
		C.free(unsafe.Pointer(bucket.name))
		bucket.name = nil
	}
	C.free(unsafe.Pointer(bucket))
}