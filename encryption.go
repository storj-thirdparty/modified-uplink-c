// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"reflect"
	"unsafe"

	"storj.io/uplink"
)

// EncryptionKey represents a key for encrypting and decrypting data.
type EncryptionKey struct {
	*uplink.EncryptionKey
}

//export uplink_derive_encryption_key
// uplink_derive_encryption_key derives a salted encryption key for passphrase using the
// salt.
//
// This function is useful for deriving a salted encryption key for users when
// implementing multitenancy in a single app bucket.
//export MAKE_CONST=1,2
func uplink_derive_encryption_key(passphrase *C.char, salt unsafe.Pointer, length C.size_t) C.Uplink_EncryptionKeyResult {
	if passphrase == nil {
		return C.Uplink_EncryptionKeyResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.Uplink_EncryptionKeyResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var goSalt []byte
	*(*reflect.SliceHeader)(unsafe.Pointer(&goSalt)) = reflect.SliceHeader{
		Data: uintptr(salt),
		Len:  ilength,
		Cap:  ilength,
	}

	encKey, err := uplink.DeriveEncryptionKey(C.GoString(passphrase), goSalt)
	if err != nil {
		return C.Uplink_EncryptionKeyResult{
			error: mallocError(err),
		}
	}

	return C.Uplink_EncryptionKeyResult{
		encryption_key: (*C.Uplink_EncryptionKey)(mallocHandle(universe.Add(&EncryptionKey{encKey}))),
	}
}

//export uplink_free_encryption_key_result
// uplink_free_encryption_key_result frees the resources associated with encryption key.
func uplink_free_encryption_key_result(result C.Uplink_EncryptionKeyResult) {
	uplink_free_error(result.error)
	freeEncryptionKey(result.encryption_key)
}

func freeEncryptionKey(encryptionKey *C.Uplink_EncryptionKey) {
	if encryptionKey == nil {
		return
	}
	defer C.free(unsafe.Pointer(encryptionKey))
	defer universe.Del(encryptionKey._handle)
}
