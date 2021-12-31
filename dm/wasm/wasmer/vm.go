/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wasmer

import (
	"github.com/pingcap/tiflow/dm/wasm/common"
	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
)

type VM struct {
	engine *wasmerGo.Engine
	store  *wasmerGo.Store
}

func NewWasmerVM() common.WasmVM {
	vm := &VM{}
	vm.Init()

	return vm
}

func (w *VM) Name() string {
	return "wasmer"
}

func (w *VM) Init() {
	w.engine = wasmerGo.NewEngine()
	w.store = wasmerGo.NewStore(w.engine)
}

func (w *VM) NewModule(wasmBytes []byte) common.WasmModule {
	if len(wasmBytes) == 0 {
		return nil
	}

	m, err := wasmerGo.NewModule(w.store, wasmBytes)
	if err != nil {
		return nil
	}

	return NewWasmerModule(w, m, wasmBytes)
}
