// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package client

import (
	"time"

	"github.com/m3db/m3x/pool"
	xretry "github.com/m3db/m3x/retry"
	xtime "github.com/m3db/m3x/time"
)

var writeAttemptArgsZeroed writeAttemptArgs

type writeAttempt struct {
	args writeAttemptArgs

	session *session

	attemptFn xretry.Fn
}

type writeAttemptArgs struct {
	namespace  string
	id         string
	t          time.Time
	value      float64
	unit       xtime.Unit
	annotation []byte
}

func (w *writeAttempt) reset() {
	w.args = writeAttemptArgsZeroed
}

func (w *writeAttempt) perform() error {
	return w.session.writeAttempt(w.args.namespace, w.args.id,
		w.args.t, w.args.value, w.args.unit, w.args.annotation)
}

type writeAttemptPool struct {
	pool    pool.ObjectPool
	session *session
}

func newWriteAttemptPool(
	session *session,
	opts pool.ObjectPoolOptions,
) *writeAttemptPool {
	p := pool.NewObjectPool(opts)
	return &writeAttemptPool{pool: p, session: session}
}

func (p *writeAttemptPool) Init() {
	p.pool.Init(func() interface{} {
		w := &writeAttempt{session: p.session}
		// NB(r): Bind attemptFn once to avoid creating receiver
		// and function method pointer over and over again
		w.attemptFn = w.perform
		w.reset()
		return w
	})
}

func (p *writeAttemptPool) Get() *writeAttempt {
	return p.pool.Get().(*writeAttempt)
}

func (p *writeAttemptPool) Put(w *writeAttempt) {
	w.reset()
	p.pool.Put(w)
}