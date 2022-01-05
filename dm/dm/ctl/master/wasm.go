// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package master

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/pingcap/errors"
	"github.com/pingcap/tiflow/dm/dm/ctl/common"
	"github.com/pingcap/tiflow/dm/dm/pb"
	"github.com/spf13/cobra"
	"github.com/wasmerio/wasmer-go/wasmer"
)

// NewWasmCmd creates a wasm command.
func NewWasmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm <command>",
		Short: "manage wasm modules",
	}
	cmd.AddCommand(
		newWasmShowCmd(),
		newWasmLoadCmd(),
		newWasmUnloadCmd(),
	)

	return cmd
}

func newWasmShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show wasm modules",
		RunE:  showWasmModulesFunc,
	}
	return cmd
}

func newWasmLoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load <path> [--override]",
		Short: "load wasm module",
		RunE:  loadWasmModulesFunc,
	}
	cmd.Flags().BoolP("override", "", false, "override if wasm module with the same name already exists")
	return cmd
}

func newWasmUnloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unload <name>",
		Short: "unload wasm module",
		RunE:  unloadWasmModulesFunc,
	}
	return cmd
}

func showWasmModulesFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp := &pb.QueryWasmModulesResponse{}
	err := common.SendRequest(
		ctx,
		"QueryWasmModules",
		&pb.QueryWasmModulesRequest{},
		&resp,
	)
	if err != nil {
		return err
	}

	for i := 0; i < len(resp.Modules); i++ {
		// wasm 模块可能很多，全都打印出来会影响展示效果。
		resp.Modules[i].Content = nil
	}
	common.PrettyPrintResponse(resp)
	return nil
}

func validateWasmModule(wasmBytes []byte) error {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	module, err := wasmer.NewModule(store, wasmBytes)
	_ = module
	if err != nil {
		return errors.Wrap(err, "invalid wasm module")
	}
	return nil
}

func loadWasmModulesFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}

	wasmPath := args[0]
	content, err := os.ReadFile(wasmPath)
	if err != nil {
		return err
	}
	if err := validateWasmModule(content); err != nil {
		return err
	}

	md5Hash := md5.Sum(content)
	md5HashHex := hex.EncodeToString(md5Hash[:])
	override, err := cmd.Flags().GetBool("override")
	if err != nil {
		return err
	}
	req := &pb.LoadWasmModuleRequest{
		Override: override,
		Module: &pb.WasmModule{
			Name:    filepath.Base(wasmPath),
			Md5:     md5HashHex,
			Content: content,
		},
	}
	resp := &pb.LoadWasmModuleResponse{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = common.SendRequest(
		ctx,
		"LoadWasmModule",
		req,
		&resp,
	)
	if err != nil {
		return err
	}
	common.PrettyPrintResponse(resp)
	return nil
}

func unloadWasmModulesFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}

	req := &pb.UnLoadWasmModuleRequest{Name: args[0]}
	resp := &pb.UnLoadWasmModuleResponse{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := common.SendRequest(
		ctx,
		"UnloadWasmModule",
		req,
		&resp,
	)
	if err != nil {
		return err
	}
	common.PrettyPrintResponse(resp)
	return nil
}
