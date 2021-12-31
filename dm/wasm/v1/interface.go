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

package v1

import "github.com/pingcap/tiflow/dm/wasm/common"

// Exports contains ABI that exported by wasm module.
type Exports interface {
}

type ImportsHandler interface {
	// utils
	Log(level LogLevel, msg string) WasmResult
	GetRootContextID() int32
	SetEffectiveContextID(contextID int32) WasmResult
	SetTickPeriodMilliseconds(tickPeriodMilliseconds int32) WasmResult
	GetCurrentTimeNanoseconds() (int32, WasmResult)
	Done() WasmResult

	// property
	GetProperty(key string) (string, WasmResult)
	SetProperty(key string, value string) WasmResult

	// foreign
	CallForeignFunction(funcName string, param []byte) ([]byte, WasmResult)
	GetFuncCallData() common.IoBuffer
}
