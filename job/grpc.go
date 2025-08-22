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

package job

import (
	v1 "github.com/cybergarage/go-job/job/api/gen/go/v1"
)

func newQueryFromGrpcQuery(query *v1.Query) (Query, error) {
	queryOpts := []QueryOption{}
	queryKind := query.Kind
	if queryKind != nil && 0 < len(*queryKind) {
		queryOpts = append(queryOpts, WithQueryKind(*queryKind))
	}
	queryUUID := query.Uuid
	if queryUUID != nil && 0 < len(*queryUUID) {
		uuid, err := NewUUIDFrom(*queryUUID)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, WithQueryUUID(uuid))
	}
	queryState := query.State
	if queryState != nil {
		state, err := newStateFrom(*queryState)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, WithQueryState(state))
	}

	return NewQuery(queryOpts...), nil
}

func newGrpcQueryFromQuery(query Query) *v1.Query {
	pbQuery := &v1.Query{
		Kind:  nil,
		Uuid:  nil,
		State: nil,
	}
	kind, ok := query.Kind()
	if ok {
		pbQuery.Kind = &kind
	}
	id, ok := query.UUID()
	if ok {
		idStr := id.String()
		pbQuery.Uuid = &idStr
	}
	state, ok := query.State()
	if ok {
		pbState, err := state.protoState()
		if err != nil {
			return nil
		}
		pbQuery.State = &pbState
	}
	return pbQuery
}
