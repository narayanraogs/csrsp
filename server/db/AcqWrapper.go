package db

import (
	"context"
	"csrspServer/db/sqlc"
)

// Gets the DAS ID for the provided DAS Name
func GetDASID(name string) (int32, error) {
	return global.GetDASID(context.Background(), name)
}

// Get the data acquistion system details provided the das id
func GetDASDetails(id int) (sqlc.Dasdetail, error) {
	return global.GetDASDetails(context.Background(), int32(id))
}

// Gets the details of chains to be acquired for the particular acquisition mode
func GetAcquisitionChainDetails(acqMode string) ([]sqlc.Daspath, error) {
	return global.GetChainDetails(context.Background(), acqMode)
}

// Gets all the DAS from which BER data is to be acquired
func GetUniqueDASForBER(acqMode string) ([]int32, error) {
	return global.GetUniqueDASForBER(context.Background(), acqMode)
}

// Gets all the BER Logging details for the provied Acqusition Mode
func GetBERLoggingDetails(acqMode string) ([]sqlc.Berlogging, error) {
	return global.GetBERLoggingDetails(context.Background(), acqMode)
}

// Get all the details of the DAS systems in data base
func GetAllDASDetails() ([]sqlc.Dasdetail, error) {
	return global.GetAllDASDetails(context.Background())
}

// Get the acquisition Id for the provided acquistion mode
func GetAcqID(acqMode string) (int32, error) {
	return global.GetAcqID(context.Background(), acqMode)
}

// GetUniqueCortexAndDPUID ... Gets the unique Cortex And DPU ID
func GetUniqueCortexAndDPUID() ([]sqlc.GetDistinctDASIDAndDPUNumberRow, error) {
	return global.GetDistinctDASIDAndDPUNumber(context.Background())
}

// Gets the unique DAS ID for Acq Mode
func GetUniqueDASIDForAcqMode(acqMode string) ([]int32, error) {
	return global.GetDistinctDASID(context.Background(), acqMode)
}

// Get DAS Component Details By Component ID
func GetDASComponentDeatilsByID(compID int) (sqlc.Dascomponentdetail, error) {
	return global.GetDASComponentDetails(context.Background(), int32(compID))
}

// Get DAS Component Details By Component Code
func GetDASComponentDeatilsByCompCode(compCode int) (sqlc.Dascomponentdetail, error) {
	return global.GetDASComponentDetails(context.Background(), int32(compCode))
}

// Get all DAS Component Details which are enabled
func GetEnabledDASComponents() ([]sqlc.Dascomponentdetail, error) {
	return global.GetEnabledDASComponents(context.Background())
}

// Get all DAS Component Details which are Dependent parent
func GetDependentDASComponents() ([]sqlc.Dascomponentdetail, error) {
	return global.GetDependentDASComponents(context.Background())
}

// Gets the list of latest n acquisitions
func GetLatestAcquisitionsList(getAll bool, noOfAcqs int) ([]sqlc.Acquisition, error) {
	acq, err := global.GetAllAcquisitions(context.Background())
	if err != nil {
		return nil, err
	}
	if getAll || len(acq) < noOfAcqs {
		return acq, nil
	}
	return acq[:noOfAcqs], nil
}

// Gets the acqusition details by Acqusition ID and system name
func GetAcquisitionByAcqIDSystemName(acqID int, systemName string) (sqlc.Acquisition, error) {
	var arg sqlc.GetAcquistionDetailParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	return global.GetAcquistionDetail(context.Background(), arg)
}

// Gets the acquisition details by date and time
func GetAcquisitionByTime(acqDate string, acqTime string) (sqlc.Acquisition, error) {
	var arg sqlc.GetAcquistionDetailByTimeParams
	arg.Date = acqDate
	arg.Time = acqTime
	return global.GetAcquistionDetailByTime(context.Background(), arg)
}

// Archives the acquisition
func AddNewAcquisition(systemName, satName, testPhase, acqMode, acqDate, acqTime, configName, remarks string) (int64, error) {
	var arg sqlc.CreateAcquisitionParams
	arg.Systemname = systemName
	arg.Satname = satName
	arg.Testphase = testPhase
	arg.Acqmode = acqMode
	arg.Date = acqDate
	arg.Time = acqTime
	arg.Configname = configName
	arg.Remark = remarks
	acq, err := global.CreateAcquisition(context.Background(), arg)
	if err != nil {
		return -1, err
	}
	return acq.LastInsertId()
}

// ChangeRemarks ... Change the remarks for a particular dataset
func ChangeRemarks(remark string, acqDate string, acqTime string) error {
	var arg sqlc.ChangeRemarkParams
	arg.Remark = remark
	arg.Date = acqDate
	arg.Time = acqTime
	return global.ChangeRemark(context.Background(), arg)
}

// GetDASSystemsByAcquisitionMode retrieves all unique DAS systems for a given acq mode.
func GetDASSystemsByAcquisitionMode(ctx context.Context, acqMode string) ([]sqlc.Dasdetail, error) {
	return global.GetDASSystemsByAcquisitionMode(ctx, acqMode)
}

// GetDASConfigurations retrieves all configurations for a given DAS and acquisition mode.
func GetDASConfigurations(ctx context.Context, params sqlc.GetDASConfigurationsParams) ([]sqlc.Dasconfiguration, error) {
	return global.GetDASConfigurations(ctx, params)
}
