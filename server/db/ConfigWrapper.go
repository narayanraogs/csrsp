package db

import (
	"context"
	"csrsp/server/db/sqlc"
)

func GetAcquisitionModes(mode string) ([]string, error) {
	m := sqlc.AcquisitionmodeAcqtype(mode)
	return global.GetAcquisitionModes(context.Background(), m)
}

// Get all configurations in the database
func GetAllConfigNames() ([]string, error) {
	return global.GetAllConfigurations(context.Background())
}

// Gets all the ConfigNames for the provided Acq Mode
func GetConfigNameForSetAcqMode(acqMode string) ([]string, error) {
	return global.GetConfigNameForAcqMode(context.Background(), acqMode)
}

// Gets all the Payload names
func GetAllPayloadNames() ([]string, error) {
	return global.GetPayloads(context.Background())
}

// Gets all the testphases in the database
func GetTestPhases() ([]string, error) {
	return global.GetTestPhases(context.Background())
}

// Gets the currently selected test phase from Database
func GetSelectedTestPhase() (string, error) {
	return global.GetSelectedTestPhase(context.Background())
}

// Gets the Remote Parameter Details
func GetRemoteParamDetails(paramID int) (sqlc.Remoteparameter, error) {
	return global.GetRemoteParameterDetails(context.Background(), int64(paramID))
}

// Gets the HK Files Details
func GetHKFiles() ([]sqlc.Hkdetail, error) {
	return global.GetHKDetails(context.Background())
}

// Create a new testphase in the database and select it
func CreateTestPhase(testPhase string) error {
	err := global.DeselectTestPhase(context.Background())
	if err != nil {
		return err
	}
	return global.CreateTestPhase(context.Background(), testPhase)
}

// Change the selected Testphase in the database
func ChangeTestPhase(testPhase string) error {
	err := global.DeselectTestPhase(context.Background())
	if err != nil {
		return err
	}
	return global.SelectTestPhase(context.Background(), testPhase)
}

// Change the  IP of the specified Data Acquisition System
func ChangeDASIP(name string, ip string) error {
	var arg sqlc.ChangeDASIPParams
	arg.Dasname = name
	arg.Ipaddress = ip
	return global.ChangeDASIP(context.Background(), arg)
}
