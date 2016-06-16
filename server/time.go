// Copyright 2016 Olivier Wulveryck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"log"
	"time"
)

// JSONTime is present because The JSON time format expected is not the default used by go
// Let's implement the interface Marshaler and override the format
type JSONTime time.Time

// MarshalJSON method for my custom type
func (t JSONTime) MarshalJSON() ([]byte, error) {
	const layout = "2006-01-02T15:04:05.000Z"
	log.Println(time.Time(t).Format(layout))
	return []byte(time.Time(t).Format(layout)), nil
}
