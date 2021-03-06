// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iterator

import (
	"context"
)

// Repository interface.
type Repository interface {
	// GetData - get rows from table.
	GetData(ctx context.Context, table, key string, fields []string,
		offset, limit int) ([]map[string]interface{}, error)
	// CreateStream - create stream.
	CreateStream(ctx context.Context, stream, table string) error
	// GetTrackingData - get rows from tracking table.
	GetTrackingData(ctx context.Context, stream, trackingTable string, fields []string,
		offset, limit int,
	) ([]map[string]interface{}, error)
	// CreateTrackingTable - create tracking table
	CreateTrackingTable(ctx context.Context, trackingTable, table string) error
	// Close - shutdown repository.
	Close() error
}
