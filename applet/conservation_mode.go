package applet

type ConservationModeStatus int16

const (
	ConservationModeStatusUnavailable ConservationModeStatus = -1
	ConservationModeStatusDisabled    ConservationModeStatus = 0
	ConservationModeStatusEnabled     ConservationModeStatus = 1
)

func GetConservationModeStatus() (ConservationModeStatus, error) {
	obj := conn.Object("com.retrixe.LenovoCtrl.v0", "/com/retrixe/LenovoCtrl/v0")
	call := obj.Call("com.retrixe.LenovoCtrl.v0.GetConservationModeStatus", 0)
	if call.Err != nil {
		return ConservationModeStatusUnavailable, call.Err
	}
	var conservationModeStatus int16
	err := call.Store(&conservationModeStatus)
	if err != nil {
		return ConservationModeStatusUnavailable, err
	}
	return ConservationModeStatus(conservationModeStatus), nil
}

func SetConservationMode(enabled bool) error {
	obj := conn.Object("com.retrixe.LenovoCtrl.v0", "/com/retrixe/LenovoCtrl/v0")
	call := obj.Call("com.retrixe.LenovoCtrl.v0.SetConservationMode", 0, enabled)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
