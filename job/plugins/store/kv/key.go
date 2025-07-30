// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kv

import (
	"bytes"

	"github.com/google/uuid"
)

// Key represents a key in the key-value store.
type Key string

// KeyTypePrefix is a prefix for key types in the key-value store.
type KeyTypePrefix = string

const (
	instancePrefix      KeyTypePrefix = "i:"
	instanceStatePrefix KeyTypePrefix = "s:"
	instanceLogPrefix   KeyTypePrefix = "l:"
)

func newKeyFromUUID(prefix string, uuid uuid.UUID, suffixes ...string) Key {
	key := Key(prefix + uuid.String())
	for _, suffix := range suffixes {
		key += Key(":" + suffix)
	}
	return key
}

// UUID returns the UUID representation of the key.
func (k Key) UUID() (uuid.UUID, error) {
	return uuid.Parse(string(k)[len(instancePrefix):])
}

// Equal checks if two keys are equal.
func (k Key) Equal(other Key) bool {
	return bytes.Equal(k.Bytes(), other.Bytes())
}

// Bytes returns the byte representation of the key.
func (k Key) Bytes() []byte {
	return []byte(k)
}

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}
