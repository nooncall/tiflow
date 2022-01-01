package v1

func (a *ABIContext) CallWasmFunction(funcName string, args ...interface{}) (interface{}, Action, error) {
	ff, err := a.Instance.GetExportsFunc(funcName)
	if err != nil {
		return nil, ActionContinue, err
	}

	res, err := ff.Call(args...)
	if err != nil {
		a.Instance.HandleError(err)
		return nil, ActionContinue, err
	}

	// if we have sync call, e.g. HttpCall, then unlock the wasm instance and wait until it resp
	action := a.Imports.Wait()

	return res, action, nil
}

func (a *ABIContext) ProxyOnContextCreate(contextID int32, parentContextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_context_create", contextID, parentContextID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ABIContext) ProxyOnDone(contextID int32) (int32, error) {
	res, _, err := a.CallWasmFunction("proxy_on_done", contextID)
	if err != nil {
		return 0, err
	}
	return res.(int32), nil
}

func (a *ABIContext) ProxyOnLog(contextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_log", contextID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ABIContext) ProxyOnDelete(contextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_delete", contextID)
	if err != nil {
		return err
	}
	return nil
}
