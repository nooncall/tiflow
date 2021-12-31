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

import (
	"github.com/pingcap/tiflow/dm/wasm/common"
)

func RegisterImports(instance common.WasmInstance) {
	_ = instance.RegisterFunc("env", "proxy_log", ProxyLog)

	_ = instance.RegisterFunc("env", "proxy_set_effective_context", ProxySetEffectiveContext)

	//_ = instance.RegisterFunc("env", "proxy_get_property", ProxyGetProperty)
	//_ = instance.RegisterFunc("env", "proxy_set_property", ProxySetProperty)

	_ = instance.RegisterFunc("env", "proxy_get_buffer_bytes", ProxyGetBufferBytes)
	_ = instance.RegisterFunc("env", "proxy_set_buffer_bytes", ProxySetBufferBytes)

	_ = instance.RegisterFunc("env", "proxy_done", ProxyDone)

	_ = instance.RegisterFunc("env", "proxy_call_foreign_function", ProxyCallForeignFunction)
}

func ProxyLog(instance common.WasmInstance, level int32, logDataPtr int32, logDataSize int32) int32 {
	logContent, err := instance.GetMemory(uint64(logDataPtr), uint64(logDataSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	callback := getImportHandler(instance)

	return callback.Log(LogLevel(level), string(logContent)).Int32()
}

func ProxySetEffectiveContext(instance common.WasmInstance, contextID int32) int32 {
	ctx := getImportHandler(instance)
	return ctx.SetEffectiveContextID(contextID).Int32()
}

func ProxySetTickPeriodMilliseconds(instance common.WasmInstance, tickPeriodMilliseconds int32) int32 {
	ctx := getImportHandler(instance)
	return ctx.SetTickPeriodMilliseconds(tickPeriodMilliseconds).Int32()
}

func ProxyGetCurrentTimeNanoseconds(instance common.WasmInstance, resultUint64Ptr int32) int32 {
	ctx := getImportHandler(instance)

	nano, res := ctx.GetCurrentTimeNanoseconds()
	if res != WasmResultOk {
		return res.Int32()
	}

	err := instance.PutUint32(uint64(resultUint64Ptr), uint32(nano))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	return WasmResultOk.Int32()
}

func ProxyDone(instance common.WasmInstance) int32 {
	ctx := getImportHandler(instance)
	return ctx.Done().Int32()
}

func ProxyCallForeignFunction(instance common.WasmInstance, funcNamePtr int32, funcNameSize int32,
	paramPtr int32, paramSize int32, returnData int32, returnSize int32) int32 {

	funcName, err := instance.GetMemory(uint64(funcNamePtr), uint64(funcNameSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	param, err := instance.GetMemory(uint64(paramPtr), uint64(paramSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	ctx := getImportHandler(instance)

	ret, res := ctx.CallForeignFunction(string(funcName), param)
	if res != WasmResultOk {
		return res.Int32()
	}

	return copyBytesIntoInstance(instance, ret, returnData, returnSize).Int32()
}
