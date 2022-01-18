// Copyright 2022 ByteDance Inc.
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

package gctuner

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testState struct {
	count int32
}

func TestFinalizer(t *testing.T) {
	maxCount := int32(16)
	is := assert.New(t)
	state := &testState{}
	f := newFinalizer(func() {
		n := atomic.AddInt32(&state.count, 1)
		if n > maxCount {
			t.Fatalf("cannot exec finalizer callback after f has been gc")
		}
	})
	for i := int32(1); i <= maxCount; i++ {
		runtime.GC()
		is.Equal(atomic.LoadInt32(&state.count), i)
	}
	is.Nil(f.ref)

	f.stop()
	is.Equal(atomic.LoadInt32(&state.count), maxCount)
	runtime.GC()
	is.Equal(atomic.LoadInt32(&state.count), maxCount)
	runtime.GC()
	is.Equal(atomic.LoadInt32(&state.count), maxCount)
}
