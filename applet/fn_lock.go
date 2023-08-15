package applet

type KeyboardFnLockStatus int16

const (
	KeyboardFnLockStatusUnavailable KeyboardFnLockStatus = -1
	KeyboardFnLockStatusDisabled    KeyboardFnLockStatus = 0
	KeyboardFnLockStatusEnabled     KeyboardFnLockStatus = 1
)

func GetKeyboardFnLockStatus() (KeyboardFnLockStatus, error) {
	obj := conn.Object("com.retrixe.LenovoCtrl.v0", "/com/retrixe/LenovoCtrl/v0")
	call := obj.Call("com.retrixe.LenovoCtrl.v0.GetKeyboardFnLockStatus", 0)
	if call.Err != nil {
		return KeyboardFnLockStatusUnavailable, call.Err
	}
	var fnLockStatus int16
	err := call.Store(&fnLockStatus)
	if err != nil {
		return KeyboardFnLockStatusUnavailable, err
	}
	return KeyboardFnLockStatus(fnLockStatus), nil
}

func SetKeyboardFnLock(enabled bool) error {
	obj := conn.Object("com.retrixe.LenovoCtrl.v0", "/com/retrixe/LenovoCtrl/v0")
	call := obj.Call("com.retrixe.LenovoCtrl.v0.SetKeyboardFnLock", 0, enabled)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
