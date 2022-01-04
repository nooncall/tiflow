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

package worker

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	toolutils "github.com/pingcap/tidb-tools/pkg/utils"
	"github.com/pingcap/tiflow/dm/dm/pb"
	"github.com/pingcap/tiflow/dm/pkg/log"
	"github.com/pingcap/tiflow/dm/pkg/terror"
	"github.com/pingcap/tiflow/dm/pkg/utils"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (s *Server) WatchWasm() error {
	err := os.MkdirAll(s.cfg.WasmDir, 0o777)
	if err != nil {
		log.L().Error("create wasm cache dir failed", zap.Stringp("wasm_dir", &s.cfg.WasmDir))
		return err
	}
	ctx := context.Background()

	ch := make(chan struct{}, 10)
	ch <- struct{}{}
	go func() {
		watchCh := s.etcdClient.Watch(ctx, "/wasm", clientv3.WithPrefix())
		for {
			res := <-watchCh
			log.L().Info("etcd watch response", zap.Any("watch_response", res))
			ch <- struct{}{}
		}
	}()

	go func() {
		for {
			<-ch
			err := s.reloadWasmModules()
			if err != nil {
				log.L().Error("reloadWasmModules() failed", zap.Error(err))
			}
		}
	}()

	return nil
}

func (s *Server) newMasterClient() (pb.MasterClient, error) {
	tls, err := toolutils.NewTLS(s.cfg.Security.SSLCA, s.cfg.Security.SSLCert, s.cfg.Security.SSLKey, "", s.cfg.Security.CertAllowedCN)
	if err != nil {
		return nil, err
	}

	endpoints := s.etcdClient.Endpoints()
	for _, endpoint := range endpoints {
		//nolint:staticcheck
		conn, err := grpc.Dial(utils.UnwrapScheme(endpoint), tls.ToGRPCDialOption(), grpc.WithBackoffMaxDelay(3*time.Second), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
		if err == nil {
			masterClient := pb.NewMasterClient(conn)
			return masterClient, nil
		}
	}

	return nil, terror.ErrCtlGRPCCreateConn.AnnotateDelegate(err, "can't connect to %s", strings.Join(endpoints, ","))
}

func (s *Server) reloadWasmModules() error {
	masterClient, err := s.newMasterClient()
	if err != nil {
		return err
	}

	resp, err := masterClient.QueryWasmModules(context.Background(), &pb.QueryWasmModulesRequest{})
	if err != nil {
		return err
	}

	newModulesMap := map[string]*pb.WasmModule{}
	for _, module := range resp.Modules {
		newModulesMap[module.Name] = module
	}

	localModules, err := s.WasmCacheListModules()
	if err != nil {
		return err
	}
	localModulesMap := map[string]*pb.WasmModule{}
	for _, module := range localModules {
		localModulesMap[module.Name] = module
	}

	// calc differences
	toDeleteModuleNames := []string{}
	toAddModules := []*pb.WasmModule{}
	toUpdateModules := []*pb.WasmModule{}
	for name := range localModulesMap {
		if _, ok := newModulesMap[name]; !ok {
			toDeleteModuleNames = append(toDeleteModuleNames, name)
		}
	}
	for name, newModule := range newModulesMap {
		localModule, ok := localModulesMap[name]
		if !ok {
			toAddModules = append(toAddModules, newModule)
		} else {
			if bytes.Compare(newModule.Content, localModule.Content) != 0 {
				toUpdateModules = append(toUpdateModules, newModule)
			}
		}
	}

	for _, name := range toDeleteModuleNames {
		if err := s.WasmCacheDeleteModule(name); err != nil {
			return err
		}
	}
	for _, module := range toAddModules {
		if err := s.WasmCacheWriteModule(module); err != nil {
			return err
		}
	}
	for _, module := range toUpdateModules {
		if err := s.WasmCacheWriteModule(module); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) WasmCacheListModules() ([]*pb.WasmModule, error) {
	entries, err := ioutil.ReadDir(s.cfg.WasmDir)
	if err != nil {
		return nil, err
	}

	modules := []*pb.WasmModule{}
	for _, entry := range entries {
		content, err := ioutil.ReadFile(path.Join(s.cfg.WasmDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		module := &pb.WasmModule{
			Name:    entry.Name(),
			Content: content,
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func (s *Server) WasmCacheWriteModule(module *pb.WasmModule) error {
	f, err := os.OpenFile(path.Join(s.cfg.WasmDir, module.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(module.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) WasmCacheDeleteModule(moduleName string) error {
	return os.Remove(path.Join(s.cfg.WasmDir, moduleName))
}
