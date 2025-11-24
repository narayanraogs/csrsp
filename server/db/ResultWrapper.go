package db

import (
	"context"
	"csrsp/server/db/sqlc"
	"csrsp/server/utils/slice"
	"database/sql"
	"fmt"
	"slices"
)

// Gets all the satellite names in the acquisitions table
func GetAcquiredSatelliteNames() ([]string, error) {
	return global.GetAcquiredSatelliteNames(context.Background())
}

// Gets all the test phases in the acquisitions table
func GetAcquiredTestPhases(satName string) ([]string, error) {
	tp, err := global.GetAcquiredTestPhases(context.Background(), satName)
	if err != nil {
		return nil, err
	}
	tp = slice.Unique(tp)
	return tp, nil
}

// Get the result profile id for the provided for the result name
func GetResultProfileID(profileName string) (int32, error) {
	return global.GetResultProfileID(context.Background(), profileName)
}

// Get saved files where FrameID is in the given frame id list
func GetFilteredSavedFiles(inputAcqIDs []int, inputSystemNames []string, inputFrameIDs []int) ([]sqlc.Savedfile, error) {
	var tbr = make([]sqlc.Savedfile, 0)
	var fids = make([]int32, 0)
	for i := 0; i < len(inputAcqIDs); i++ {
		fids = append(fids, int32(inputAcqIDs[i]))
	}
	for i := 0; i < len(inputAcqIDs); i++ {
		var filterArgs sqlc.FilterSavedFilesParams
		filterArgs.Acquisitionid = int32(inputAcqIDs[i])
		filterArgs.Systemname = inputSystemNames[i]
		filterArgs.Frameids = fids
		savedFiles, err := global.FilterSavedFiles(context.Background(), filterArgs)
		if err != nil {
			return nil, err
		}
		tbr = append(tbr, savedFiles...)
	}
	return tbr, nil
}

// Get SavedFiles for the provided FrameIDs
func FilterSavedFilesByFrameIDs(inputFrameIDs []int) ([]sqlc.Savedfile, error) {
	var fids = make([]int32, 0)
	for i := 0; i < len(inputFrameIDs); i++ {
		fids = append(fids, int32(inputFrameIDs[i]))
	}
	return global.GetSavedFilesForFrameIDs(context.Background(), fids)
}

// Gets the list of acquisitions where acqID in the given list of AcqIDs
func GetAcquisitionList(acqIDs []int, systemNames []string, testPhase string) ([]sqlc.Acquisition, error) {
	acqList := make([]sqlc.Acquisition, 0)
	for i := 0; i < len(acqIDs); i++ {
		var arg sqlc.GetAcqusitionParams
		arg.Acquisitionid = int32(acqIDs[i])
		arg.Systemname = systemNames[i]
		arg.Testphase = testPhase
		acq, err := global.GetAcqusition(context.Background(), arg)
		if err != nil {
			return acqList, err
		}
		acqList = append(acqList, acq)
	}
	return acqList, nil
}

// Gets the list of all acquisitions having specified criteria
func GetOfflineProcessingResultSet(satName string, testPhase string, filterParams []string, filterValues []string) ([]sqlc.Acquisition, error) {
	var arg sqlc.GetOfflineProcessingResultsParams
	arg.Satname = satName
	arg.Testphase = testPhase
	arg.Configname = getDBParam("Config Name", filterParams, filterValues)
	arg.Date = getDBParam("Date", filterParams, filterValues)
	arg.Time = getDBParam("Time", filterParams, filterValues)
	arg.Acqmode = getDBParam("Acq Mode", filterParams, filterValues)
	arg.Remark = getDBParam("Remarks", filterParams, filterValues)

	return global.GetOfflineProcessingResults(context.Background(), arg)
}

func getDBParam(filter string, filterParams []string, filterValues []string) string {
	index := slices.Index(filterParams, filter)
	if index == -1 {
		return "%"
	}
	return filterValues[index]
}

// Gets the list of all acquisitions filtered by Configuration
func GetOfflineProcessingResultSetConfig(satName string, configName string) ([]sqlc.Acquisition, error) {
	var arg sqlc.GetOfflineProcessingResultsParams
	arg.Satname = satName
	arg.Testphase = "%"
	arg.Configname = configName
	arg.Date = "%"
	arg.Time = "%"
	arg.Acqmode = "%"
	arg.Remark = "%"

	return global.GetOfflineProcessingResults(context.Background(), arg)
}

type SavedFile struct {
	FileName  string
	FrameType string
	FrameID   int
}

// Returns FileNames, FrameTypes and FrameIDs for that acquisition
func GetAllArchivedFilesForAcquisition(acqDate string, acqTime string) ([]SavedFile, error) {
	var arg sqlc.GetOfflineProcessingResultsParams
	arg.Satname = "%"
	arg.Testphase = "%"
	arg.Configname = "%"
	arg.Date = acqDate
	arg.Time = acqTime
	arg.Acqmode = "%"
	arg.Remark = "%"

	acqs, err := global.GetOfflineProcessingResults(context.Background(), arg)
	if err != nil {
		return nil, err
	}
	if len(acqs) == 0 {
		return nil, fmt.Errorf("no acquistions for the provided date and time")
	}

	frameTypes, err := global.GetAcquiredFrameTypes(context.Background())
	if err != nil {
		return nil, err
	}

	var fids = make([]int32, 0)
	var fTypeMap = make(map[int32]string)
	for _, frameType := range frameTypes {
		ids, err := global.GetFrameIDForFrameType(context.Background(), frameType)
		if err != nil {
			return nil, err
		}
		fids = append(fids, ids...)
		for _, id := range ids {
			fTypeMap[id] = frameType
		}
	}
	var acq sqlc.GetSavedFilesForAcquisitionParams
	acq.Acquisitionid = acqs[0].Acquisitionid
	acq.Systemname = acqs[0].Systemname

	savedFiles, err := global.GetSavedFilesForAcquisition(context.Background(), acq)
	if err != nil {
		return nil, err
	}
	if len(savedFiles) == 0 {
		return nil, fmt.Errorf("no files for the provided date and time")
	}
	var tbr = make([]SavedFile, 0)
	for _, savedFile := range savedFiles {
		if slices.Contains(fids, savedFile.Frameid) {
			var s SavedFile
			s.FileName = savedFile.Filepath
			s.FrameID = int(savedFile.Frameid)
			s.FrameType = fTypeMap[savedFile.Frameid]
			tbr = append(tbr, s)
		}
	}
	if len(tbr) == 0 {
		return nil, fmt.Errorf("no files for the provided date and time")
	}
	return tbr, nil
}

type ResultProfile struct {
	Profile sqlc.Resultprofiledetail
	Params  []string
}

// Gets all results in a profile
func GetAllResultsInProfile(profileID int) ([]ResultProfile, error) {
	profiles, err := global.GetResultProfileDetails(context.Background(), int32(profileID))
	if err != nil {
		return nil, err
	}
	var tbr = make([]ResultProfile, 0)
	for _, profile := range profiles {
		var r ResultProfile
		r.Profile = profile
		params, err := getParametersList(int(profile.Profileid), int(profile.Resultid))
		if err != nil {
			return nil, err
		}
		r.Params = params
		tbr = append(tbr, r)
	}
	return tbr, nil
}

func getParametersList(profileID int, resultID int) ([]string, error) {
	var parameterNames = make([]string, 0)
	var arg sqlc.GetResultProfileParametersParams
	arg.Profileid = int32(profileID)
	arg.Resultid = int32(resultID)
	params, err := global.GetResultProfileParameters(context.Background(), arg)
	if err != nil {
		return nil, err
	}
	for _, param := range params {
		p, err := global.GetParameterDetails(context.Background(), param.Paramid)
		if err != nil {
			continue
		}
		parameterNames = append(parameterNames, p.Parametername)
	}
	return parameterNames, nil
}

type AuxConsolidatedSummary struct {
	ConsolidatedReportID int
	AcquisitionID        int
	FrameType            string
	AuxDataCorrect       string
	NoOfFrames           string
}

type VideoConsolidatedSummary struct {
	ConsolidatedReportID int
	AcquisitionID        int
	FrameType            string
	FrameID              string
	MaxMean              float64
	MeanOfMean           float64
	MaxSD                float64
	MeanSD               float64
	RmsSD                float64
	Path                 string
	NoOfFrames           string
}

// Get consolidated summary for the given acq id
func GetConsolidatedSummary(acqID int, systemName string) ([]AuxConsolidatedSummary, []VideoConsolidatedSummary, error) {
	auxSummary := make([]AuxConsolidatedSummary, 0)
	videoSummary := make([]VideoConsolidatedSummary, 0)
	var arg sqlc.GetConsolidatedReportsParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	summary, err := global.GetConsolidatedReports(context.Background(), arg)
	if err != nil {
		return nil, nil, err
	}

	for _, sum := range summary {
		if sum.Auxdatacorrect.Valid {
			var temp AuxConsolidatedSummary
			temp.ConsolidatedReportID = int(sum.Consolidatedreportid)
			temp.AcquisitionID = int(sum.Acquisitionid)
			temp.AuxDataCorrect = sum.Auxdatacorrect.String
			temp.FrameType = sum.Frametype
			temp.NoOfFrames = sum.Noofframes
			auxSummary = append(auxSummary, temp)
		} else {
			var temp VideoConsolidatedSummary
			temp.ConsolidatedReportID = int(sum.Consolidatedreportid)
			temp.AcquisitionID = int(sum.Acquisitionid)
			temp.FrameType = sum.Frametype
			temp.FrameID = sum.Frameid.String
			temp.NoOfFrames = sum.Noofframes
			temp.MaxMean = sum.Maxmean.Float64
			temp.MaxSD = sum.Maxsd.Float64
			temp.MeanOfMean = sum.Meanofmean.Float64
			temp.MeanSD = sum.Meansd.Float64
			temp.RmsSD = sum.Rmssd.Float64
			temp.Path = sum.Path.String
			videoSummary = append(videoSummary, temp)
		}
	}
	return auxSummary, videoSummary, nil
}

// Get saved results for given result number
func GetSavedResult(resultNumber int, systemName string) (sqlc.Savedresult, error) {
	var arg sqlc.GetSavedResultsParams
	arg.Resultnumber = int32(resultNumber)
	arg.Systemname = systemName
	return global.GetSavedResults(context.Background(), arg)
}

// Get saved results for given acquisition ID and system name
func GetSavedResults(acqID int, systemName string) ([]sqlc.Savedresult, error) {
	var arg sqlc.GetSavedResultsForAcquisitionParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	return global.GetSavedResultsForAcquisition(context.Background(), arg)
}

// Gets all the result profiles in the database
func GetAllResultProfiles() ([]string, error) {
	return global.GetAllResultProfiles(context.Background())
}

type TrendAnalysisSingle struct {
	MaxMean    float64
	MeanOfMean float64
	MaxSD      float64
	MeanSD     float64
	FilePath   string
	NoOfFrames string
}

// Get Trend Analysis values for the specified AcqID and frame identifier
func GetTrendAnalysisValues(acqID int, systemName string, frameType string, frameIdentifier string) (TrendAnalysisSingle, error) {
	var trend TrendAnalysisSingle
	var arg sqlc.GetConsolidatedReportFrameIdentiferParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	arg.Frameid.String = frameIdentifier
	arg.Frameid.Valid = true
	arg.Frametype = frameType
	cons, err := global.GetConsolidatedReportFrameIdentifer(context.Background(), arg)
	if err != nil {
		return TrendAnalysisSingle{}, err
	}

	trend.MaxMean = cons.Maxmean.Float64
	trend.MeanOfMean = cons.Meanofmean.Float64
	trend.MaxSD = cons.Maxsd.Float64
	trend.MeanSD = cons.Meansd.Float64
	trend.FilePath = cons.Path.String
	trend.NoOfFrames = cons.Noofframes

	return trend, nil
}

// Create a new Saved Files, don't insert and return if entry already exists
func InsertSavedFiles(acqID int, frameID int, noOfLines uint64, filePath string, systemName string) error {
	var arg1 sqlc.CheckIfSavedFileExistsParams
	arg1.Acquisitionid = int32(acqID)
	arg1.Frameid = int32(frameID)
	arg1.Systemname = systemName
	count, err := global.CheckIfSavedFileExists(context.Background(), arg1)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var arg2 sqlc.CreateSavedFilesParams
	arg2.Acquisitionid = int32(acqID)
	arg2.Filepath = filePath
	arg2.Frameid = int32(frameID)
	arg2.Nooflines = int32(noOfLines)
	arg2.Systemname = systemName
	return global.CreateSavedFiles(context.Background(), arg2)
}

// Remove the result profile from the database
func RemoveResultProfile(profileName string) error {
	ctx := context.Background()
	tx, err := global.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	q := global.WithTx(tx)

	profileID, err := q.GetResultProfileID(ctx, profileName)
	if err != nil {
		return err
	}
	err = q.RemoveResultProfileParameters(ctx, profileID)
	if err != nil {
		return err
	}
	err = q.RemoveResultProfileDetails(ctx, profileID)
	if err != nil {
		return err
	}
	err = q.RemoveResultProfile(ctx, profileID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Add a new Result Profile and Add all the results
func AddResultProfile(profileName string, results []ResultProfile) error {
	ctx := context.Background()
	tx, err := global.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	q := global.WithTx(tx)

	var arg1 sqlc.CreateResultProfileParams
	arg1.Noofresults = int32(len(results))
	arg1.Profilename = profileName
	res, err := q.CreateResultProfile(ctx, arg1)
	if err != nil {
		return err
	}
	profileID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for i := 0; i < len(results); i++ {
		var arg2 sqlc.CreateResultProfileDetailsParams
		prof := results[i].Profile
		prof.Profileid = int32(profileID)
		arg2.Profileid = prof.Profileid
		arg2.Resultid = prof.Resultid
		arg2.Resultname = prof.Resultname
		arg2.Frametype = prof.Frametype
		arg2.Frameid = prof.Frameid
		arg2.Noofparameters = int32(len(results[i].Params))
		arg2.Startline = prof.Startline
		arg2.Endline = prof.Endline
		arg2.Startword = prof.Startword
		arg2.Endword = prof.Endword
		arg2.Startpixel = prof.Startpixel
		arg2.Endpixel = prof.Endpixel
		arg2.Decimalorhex = prof.Decimalorhex
		arg2.Raworprocessed = prof.Raworprocessed
		arg2.Ascordsc = prof.Ascordsc
		arg2.Ismean = prof.Ismean
		arg2.Issd = prof.Issd
		arg2.Stacklist = prof.Stacklist
		err = q.CreateResultProfileDetails(ctx, arg2)
		if err != nil {
			return err
		}

		for _, param := range results[i].Params {
			var arg3 sqlc.GetParamIDParams
			arg3.Frametype = prof.Frametype
			arg3.Parametername = param
			pid, err := global.GetParamID(ctx, arg3)
			if err != nil {
				return err
			}
			var arg4 sqlc.CreateResultProfileParametersParams
			arg4.Paramid = pid
			arg4.Profileid = prof.Profileid
			arg4.Resultid = prof.Resultid
			err = global.CreateResultProfileParameters(ctx, arg4)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// Add Aux consolidated summary for given dataset
func AddAuxConsolidatedSummary(acqID int, systemName string, frameType string, auxDataStatus string, noOfFrames string) error {
	var arg1 sqlc.CheckIfConsolidatedSummaryExistsParams
	arg1.Acquisitionid = int32(acqID)
	arg1.Frametype = frameType
	arg1.Systemname = systemName
	count, err := global.CheckIfConsolidatedSummaryExists(context.Background(), arg1)
	if err != nil {
		return err
	}
	if count > 0 {
		var arg2 sqlc.UpdateAuxConsolidatedReportParams
		arg2.Acquisitionid = int32(acqID)
		arg2.Auxdatacorrect.String = auxDataStatus
		arg2.Auxdatacorrect.Valid = true
		arg2.Frametype = frameType
		arg2.Noofframes = noOfFrames
		arg2.Systemname = systemName
		return global.UpdateAuxConsolidatedReport(context.Background(), arg2)
	}
	var arg2 sqlc.InsertAuxConsolidatedReportParams
	arg2.Acquisitionid = int32(acqID)
	arg2.Auxdatacorrect.String = auxDataStatus
	arg2.Auxdatacorrect.Valid = true
	arg2.Frametype = frameType
	arg2.Noofframes = noOfFrames
	arg2.Systemname = systemName
	return global.InsertAuxConsolidatedReport(context.Background(), arg2)
}

// Add Video consolidated summary for given dataset
func AddVideoConsolidatedSummary(acqID int, systemName string, frameType string, frameID string, maxMean float64, meanMean float64, maxSD float64, meanSD float64, RMSSD float64, path string, noOfFrames string) error {

	var arg1 sqlc.CheckIfVideoConsolidatedSummaryExistsParams
	arg1.Acquisitionid = int32(acqID)
	arg1.Frametype = frameType
	arg1.Frameid.String = frameID
	arg1.Frameid.Valid = true
	arg1.Systemname = systemName
	count, err := global.CheckIfVideoConsolidatedSummaryExists(context.Background(), arg1)
	if err != nil {
		return err
	}
	if count > 0 {
		var arg2 sqlc.UpdateVideoConsolidatedReportParams
		arg2.Acquisitionid = int32(acqID)
		arg2.Frametype = frameType
		arg2.Noofframes = noOfFrames
		arg2.Systemname = systemName
		arg2.Frameid = sql.NullString{String: frameID, Valid: true}
		arg2.Maxmean = sql.NullFloat64{Float64: maxMean, Valid: true}
		arg2.Maxsd = sql.NullFloat64{Float64: maxSD, Valid: true}
		arg2.Meansd = sql.NullFloat64{Float64: meanSD, Valid: true}
		arg2.Meanofmean = sql.NullFloat64{Float64: meanMean, Valid: true}
		arg2.Rmssd = sql.NullFloat64{Float64: RMSSD, Valid: true}
		arg2.Path = sql.NullString{String: path, Valid: true}
		return global.UpdateVideoConsolidatedReport(context.Background(), arg2)
	}
	var arg2 sqlc.InsertVideoConsolidatedReportParams
	arg2.Acquisitionid = int32(acqID)
	arg2.Frametype = frameType
	arg2.Noofframes = noOfFrames
	arg2.Systemname = systemName
	arg2.Frameid = sql.NullString{String: frameID, Valid: true}
	arg2.Maxmean = sql.NullFloat64{Float64: maxMean, Valid: true}
	arg2.Maxsd = sql.NullFloat64{Float64: maxSD, Valid: true}
	arg2.Meansd = sql.NullFloat64{Float64: meanSD, Valid: true}
	arg2.Meanofmean = sql.NullFloat64{Float64: meanMean, Valid: true}
	arg2.Rmssd = sql.NullFloat64{Float64: RMSSD, Valid: true}
	arg2.Path = sql.NullString{String: path, Valid: true}
	return global.InsertVideoConsolidatedReport(context.Background(), arg2)
}

// Add an offline result
func AddOfflineResult(acqID int, systemName string, blob []byte, frameType string, frameID string, resultName string, resultPath string) error {
	var arg sqlc.CreateOfflineResultParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	arg.Uniqueid = blob
	arg.Frametype = frameType
	arg.Frameidentifier = sql.NullString{String: frameID, Valid: true}
	arg.Resultname = resultName
	arg.Resultpath = sql.NullString{String: resultPath, Valid: true}
	return global.CreateOfflineResult(context.Background(), arg)
}

// Add an offline consolidated result
func AddConsolidatedOfflineResult(acqID int, systemName string, blob []byte, frameType string, resultName string, resultPath string) error {
	var arg sqlc.CreateOfflineResultParams
	arg.Acquisitionid = int32(acqID)
	arg.Systemname = systemName
	arg.Uniqueid = blob
	arg.Frametype = frameType
	arg.Resultname = resultName
	arg.Resultpath = sql.NullString{String: resultPath, Valid: true}
	return global.CreateOfflineResult(context.Background(), arg)
}

// Get total number of offline results in database
func GetOfflineResultsCount() (int, error) {
	res, err := global.GetMaxResultNumber(context.Background())
	if err != nil {
		return -1, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	return int(res[0]), nil
}
