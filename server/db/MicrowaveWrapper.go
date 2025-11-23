package db

import (
	"context"
	"csrspServer/db/sqlc"
)

// Gets All the Microwave Processing IDs from the Database
func GetAllSARProcessingIDs() ([]int32, error) {
	return global.GetMicrowaveProcessingID(context.Background())

}

// Gets all the Polarization values in the database
func GetMicrowavePols() ([]string, error) {
	polPID, err := global.GetMicrowavePolarization(context.Background())
	if err != nil {
		return nil, err
	}
	status, err := GetStatusMap(polPID)
	if err != nil {
		return nil, err
	}
	toBeReturned := make([]string, 0)
	for _, v := range status.Status {
		toBeReturned = append(toBeReturned, v)
	}
	return toBeReturned, nil
}

// Gets all the Timing values in the database
func GetMicrowaveTimings() ([]string, error) {
	polPID, err := global.GetMicrowaveTimingState(context.Background())
	if err != nil {
		return nil, err
	}
	status, err := GetStatusMap(polPID)
	if err != nil {
		return nil, err
	}
	toBeReturned := make([]string, 0)
	for _, v := range status.Status {
		toBeReturned = append(toBeReturned, v)
	}
	return toBeReturned, nil
}

// Gets the details for sar processing for given Process ID
func GetSARProcessingDetails(processID int) (sqlc.Sarvideodataprocessing, error) {
	return global.GetMicrowaveProcessingDetails(context.Background(), int32(processID))
}
