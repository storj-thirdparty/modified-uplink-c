#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

typedef struct Uplink_Handle {
    size_t _handle;
} Uplink_Handle;

typedef struct Uplink_Access {
    size_t _handle;
} Uplink_Access;
typedef struct Uplink_Project {
    size_t _handle;
} Uplink_Project;
typedef struct Uplink_Download {
    size_t _handle;
} Uplink_Download;
typedef struct Uplink_Upload {
    size_t _handle;
} Uplink_Upload;

typedef struct Uplink_EncryptionKey {
    size_t _handle;
} Uplink_EncryptionKey;

typedef struct Uplink_Config {
    char *user_agent;

    int32_t dial_timeout_milliseconds;

    // temp_directory specifies where to save data during downloads to use less memory.
    char *temp_directory;
} Uplink_Config;

typedef struct Uplink_Bucket {
    char *name;
    int64_t created;
} Uplink_Bucket;

typedef struct Uplink_SystemMetadata {
    int64_t created;
    int64_t expires;
    int64_t content_length;
} Uplink_SystemMetadata;

typedef struct Uplink_CustomMetadataEntry {
    char *key;
    size_t key_length;

    char *value;
    size_t value_length;
} Uplink_CustomMetadataEntry;

typedef struct Uplink_CustomMetadata {
    Uplink_CustomMetadataEntry *entries;
    size_t count;
} Uplink_CustomMetadata;

typedef struct Uplink_Object {
    char *key;
    bool is_prefix;
    Uplink_SystemMetadata system;
    Uplink_CustomMetadata custom;
} Uplink_Object;

typedef struct Uplink_UploadOptions {
    // When expires is 0 or negative, it means no expiration.
    int64_t expires;
} Uplink_UploadOptions;

typedef struct Uplink_DownloadOptions {
    int64_t offset;
    // When length is negative, it will read until the end of the blob.
    int64_t length;
} Uplink_DownloadOptions;

typedef struct Uplink_ListObjectsOptions {
    char *prefix;
    char *cursor;
    bool recursive;

    bool system;
    bool custom;
} Uplink_ListObjectsOptions;

typedef struct Uplink_ListBucketsOptions {
    char *cursor;
} Uplink_ListBucketsOptions;

typedef struct Uplink_ObjectIterator {
    size_t _handle;
} Uplink_ObjectIterator;
typedef struct Uplink_BucketIterator {
    size_t _handle;
} Uplink_BucketIterator;

typedef struct Uplink_Permission {
    bool allow_download;
    bool allow_upload;
    bool allow_list;
    bool allow_delete;

    // unix time in seconds when the permission becomes valid.
    // disabled when 0.
    int64_t not_before;
    // unix time in seconds when the permission becomes invalid.
    // disabled when 0.
    int64_t not_after;
} Uplink_Permission;

typedef struct Uplink_SharePrefix {
    char *bucket;
    // prefix is the prefix of the shared object keys.
    char *prefix;
} Uplink_SharePrefix;

typedef struct Uplink_Error {
    int32_t code;
    char *message;
} Uplink_Error;

enum {
    UPLINK_ERROR_INTERNAL = 0x02,
    UPLINK_ERROR_CANCELED = 0x03,
    UPLINK_ERROR_INVALID_HANDLE = 0x04,
    UPLINK_ERROR_TOO_MANY_REQUESTS = 0x05,
    UPLINK_ERROR_BANDWIDTH_LIMIT_EXCEEDED = 0x06,

    UPLINK_ERROR_BUCKET_NAME_INVALID = 0x10,
    UPLINK_ERROR_BUCKET_ALREADY_EXISTS = 0x11,
    UPLINK_ERROR_BUCKET_NOT_EMPTY = 0x12,
    UPLINK_ERROR_BUCKET_NOT_FOUND = 0x13,

    UPLINK_ERROR_OBJECT_KEY_INVALID = 0x20,
    UPLINK_ERROR_OBJECT_NOT_FOUND = 0x21,
    UPLINK_ERROR_UPLOAD_DONE = 0x22
};

typedef struct Uplink_AccessResult {
    Uplink_Access *access;
    Uplink_Error *error;
} Uplink_AccessResult;

typedef struct Uplink_ProjectResult {
    Uplink_Project *project;
    Uplink_Error *error;
} Uplink_ProjectResult;

typedef struct Uplink_BucketResult {
    Uplink_Bucket *bucket;
    Uplink_Error *error;
} Uplink_BucketResult;

typedef struct Uplink_ObjectResult {
    Uplink_Object *object;
    Uplink_Error *error;
} Uplink_ObjectResult;

typedef struct Uplink_UploadResult {
    Uplink_Upload *upload;
    Uplink_Error *error;
} Uplink_UploadResult;

typedef struct Uplink_DownloadResult {
    Uplink_Download *download;
    Uplink_Error *error;
} Uplink_DownloadResult;

typedef struct Uplink_WriteResult {
    size_t bytes_written;
    Uplink_Error *error;
} Uplink_WriteResult;

typedef struct Uplink_ReadResult {
    size_t bytes_read;
    Uplink_Error *error;
} Uplink_ReadResult;

typedef struct Uplink_StringResult {
    char *string;
    Uplink_Error *error;
} Uplink_StringResult;

typedef struct Uplink_EncryptionKeyResult {
    Uplink_EncryptionKey *encryption_key;
    Uplink_Error *error;
} Uplink_EncryptionKeyResult;