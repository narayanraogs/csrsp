-- name: GetParamID :one
Select `ParamID` from `ParameterDetails` 
where `FrameType` like ? and `ParameterName` like ?;

-- name: GetParameterDetails :one
Select * from `ParameterDetails` 
where `ParamID` like ?;

-- name: GetBitDetails :many
Select * From `BitDetails` 
where `ParamID` like ?;

-- name: GetFSCDetails :one
Select * from `FSCType` 
where `ParamID` like ?;

-- name: GetCRCDetails :one
Select * from `CRCType` 
where `ParamID` like ?;

-- name: GetIncrementDetails :one
Select * from `IncrementType`
where `ParamID` like ?;

-- name: GetRadixDetails :one
Select * from `RadixType` 
where `ParamID` like ?;

-- name: GetAnalogDetails :one
Select * from `AnalogType` 
where `ParamID` = ?;

-- name: GetStatusDetails :many
Select * from `StatusType` 
where `ParamID` like ?;

-- name: GetFrameTypeParamterDetails :many
Select * from `ParameterDetails`
where `FrameType` like ? order by `StartWord`, `StartBit` Asc;

-- name: GetFrameDetails :one
Select * from `Frames`
where `FrameID` like ?;

-- name: GetDASDetails :one
Select * from `DASDetails` 
where `DASID` = ?;

-- name: GetAllDASDetails :many
Select * from `DASDetails`;


-- name: GetDASID :one
Select `DASID` from `DASDetails` 
where `DASName` = ?;

-- name: GetChainDetails :many
Select * from `DASPath`
where `AcqMode` = ?;

-- name: GetFrameTypeFrameIdentifiers :many
Select * from `Frames` 
where `FrameType` like ? order by `FrameID` asc;

-- name: GetArchivableFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `ToBeArchived` = 1;

-- name: GetFrameTypeDetails :one
Select * from `FrameTypes` 
where `FrameType` = ?;

-- name: GetAcquiredFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `ProcessedOrAcquired` like '%acq%';

-- name: GetAuxFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `Type` like '%aux%';

-- name: GetOpticalFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `Type` like '%optical%';

-- name: GetMicrowaveFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `Type` like '%microwave%';

-- name: GetFrameTypes :many
Select `FrameType` from `FrameTypes`;


-- name: GetWritableFrameTypes :many
Select `FrameType` from `FrameTypes` 
where `ToBeWritten` = 1;

-- name: GetUniqueDASForBER :many
Select Distinct `DASID` from `BERLogging` 
where `AcqMode` = ?;

-- name: GetBERLoggingDetails :many
Select * from `BERLogging` 
where `AcqMode` like ?;

-- name: GetAcquisitionModes :many
Select `AcqMode` from `AcquisitionMode`
where `AcqType` like ?;

-- name: GetAcqID :one
Select `AcqID` from `AcquisitionMode` 
where `AcqMode` like ?;

-- name: GetAllConfigurations :many
Select `ConfigName` from `Configuration`;

-- name: GetConfigNameForAcqMode :many
Select `ConfigName` from `Configuration` 
where `AcqMode` like ?;

-- name: GetFrameID :one
Select `FrameID` from `Frames` 
where `FrameType` = ? and `FrameIdentifier` = ?;

-- name: GetPayloads :many
SELECT `PayloadName` FROM `Payloads`;

-- name: GetAcquiredSatelliteNames :many
SELECT Distinct `SatName` FROM `Acquisitions`;

-- name: GetTestPhases :many
SELECT `TestPhase` FROM `TestPhases`;

-- name: GetSelectedTestPhase :one
SELECT `TestPhase` FROM `TestPhases` 
where `Selected` = 1;

-- name: GetAcquiredTestPhases :many
SELECT `TestPhase` FROM `Acquisitions`
Where `SatName` like ?
Order by AcquisitionID desc;

-- name: CountOpticalSubTypes :one
Select count(*) from `VideoSubFrames` 
where `SubFrameType` like ? and `FrameID` = ?;

-- name: GetVideoDataDetails :one
Select * from `VideoDataProcessing`
where `FrameID` = ?;

-- name: GetResultProfileID :one
Select `ProfileID` from `ResultProfiles` 
where `ProfileName` like ?;

-- name: GetMicrowaveProcessingID :many
SELECT `ProcessID` FROM `SARVideoDataProcessing`;

-- name: GetMicrowavePolarization :one
SELECT `PolarizationPID` FROM `SARVideoDataProcessing`;

-- name: GetMicrowaveTimingState :one
SELECT `TimingStatePID` FROM `SARVideoDataProcessing`;

-- name: GetTypeOfFrame :one
SELECT `Type` FROM `FrameTypes` 
where `FrameType` like ?;

-- name: GetIntermediateFrames :many
SELECT `FrameType` FROM `FrameTypes` 
where `isIntermediateReport` = 1 and `Type` like ?;

-- name: GetDistinctDASIDAndDPUNumber :many
Select Distinct `DASID`, `DPUNumber` from `DASPath`;

-- name: GetDistinctDASID :many
Select Distinct `DASID` from `DASPath` 
where `AcqMode` = ?;

-- name: GetDASComponentDetails :one
Select * from `DASComponentDetails` 
where `ComponentID` = ?;

-- name: GetDASComponentDetailsByCompCode :one
Select * from `DASComponentDetails` 
where `ComponentCode` = ? and `Enabled` = 1;

-- name: GetEnabledDASComponents :many
Select * from `DASComponentDetails`
where `Enabled` = 1;

-- name: GetDependentDASComponents :many
Select * from `DASComponentDetails` 
where  `ParentCompName` IS NOT NULL;

-- name: GetAllFrameIDs :many
Select `FrameID` from `Frames`  
Order by `FrameID` asc;

-- name: GetAllAcquisitions :many
Select * from `Acquisitions` 
Order by `AcquisitionID` desc;

-- name: GetAcquistionDetail :one
Select * from `Acquisitions` 
WHERE `AcquisitionID`= ? and `SystemName` = ?;

-- name: GetAcquistionDetailByTime :one
Select * from `Acquisitions` 
WHERE `Date`= ? and `Time` = ?;

-- name: FilterSavedFiles :many
SELECT * from `SavedFiles` WHERE 
`AcquisitionID` = ? AND `SystemName` = ? 
AND FrameID IN (sqlc.slice(frameIDs)) ORDER BY `AcquisitionID` DESC;

-- name: GetSavedFilesForFrameIDs :many
SELECT * from `SavedFiles` 
WHERE `FrameID` IN (sqlc.slice(frameIDs)) ORDER BY `AcquisitionID` DESC;

-- name: GetAcqusition :one
Select * from `Acquisitions` 
WHERE `AcquisitionID` = ? AND `SystemName` like ? 
AND `TestPhase` = ? Order by `AcquisitionID` desc;

-- name: GetFrameIDsForDownload :many
Select `FrameID` from `DataDownloadPermissions` 
where `UserID` = ?;

-- name: GetOfflineProcessingResults :many
Select * from `Acquisitions` 
where `SatName` like ? and `TestPhase` like ?
and `ConfigName` like ? and `Date` like ? 
and `Time` like ? and `AcqMode` like ?
and `Remark` like ? Order by `AcquisitionID` desc;

-- name: GetFrameIDForFrameType :many
Select `FrameID` from `Frames`
Where `FrameType` like ?;

-- name: GetSavedFilesForAcquisition :many
Select * from `SavedFiles` 
where `AcquisitionID` = ? and `SystemName` = ?;

-- name: GetRemoteParameterDetails :one
Select * from `RemoteParameters` 
WHERE `ParamID` = ?;

-- name: GetResultProfileDetails :many
Select * from `ResultProfileDetails`
where `ProfileID` = ?;

-- name: GetResultProfileParameters :many
Select * from `ResultProfileParameters` 
where `ProfileID` = ? and `ResultID` = ?;

-- name: GetHKDetails :many
Select * from HKDetails;

-- name: GetConsolidatedReports :many
Select * from `ConsolidatedReport` 
Where `AcquisitionID` = ? and `SystemName` = ?;

-- name: GetConsolidatedReportFrameIdentifer :one
Select * from `ConsolidatedReport` 
where `AcquisitionID` = ? and `FrameType` like ? and 
`FrameID` like ? and `SystemName`= ?;

-- name: GetSavedResults :one
Select * from `SavedResults` 
Where `ResultNumber` = ? and `SystemName`= ?;

-- name: GetSavedResultsForAcquisition :many
Select * from `SavedResults`
Where `AcquisitionID` = ? and `SystemName` = ?;

-- name: GetAllResultProfiles :many
Select `ProfileName` from `ResultProfiles`;

-- name: GetUserDetails :one
SELECT * FROM `Users` WHERE `UserName` = ?;

-- name: GetDASSystemsByAcquisitionMode :many
SELECT DISTINCT d.*
FROM `DASDetails` d
JOIN `DASPath` p ON d.`DASID` = p.`DASID`
WHERE p.`AcqMode` = ?;

-- name: GetDASConfigurations :many
SELECT *
FROM DASConfigurations
WHERE AcquisitionMode = ? AND DASID = ?
ORDER BY ExecutionOrder ASC;

-- name: GetFSCParamID :one
Select `ParamID` from `ParameterDetails` 
where `FrameType` like ? and `ParameterType` like 'FSC';

-- name: GetMicrowaveProcessingDetails :one
Select * from `SARVideoDataProcessing` 
where `ProcessID` = ?;

-- name: GetOpticalLimitDetails :one
Select * from `VideoLimits` 
where `FrameID` = ?;

-- name: GetOpticalFaultyPixels :one
Select * from `VideoFaultyPixels` 
where `FrameID` = ?;

-- name: GetOpticalSubFrames :many
Select * from `VideoSubFrames` 
where `SubFrameType` like ? and `FrameID` = ?  
Order By `SubFrameNumber`;

-- name: GetOpticalSubFrameNumber :one
Select * from `VideoSubFrames` 
where `SubFrameType` like ? and `FrameID` = ? 
and `SubFrameName` like ?;

-- name: GetOpticalSubFrameDetails :one
Select * from `VideoSubFrames` 
where `SubFrameType` like ? and `SubFrameNumber` = ? 
and `FrameID` = ?;

-- name: CreateAcquisition :execresult
INSERT INTO `Acquisitions`
(`SystemName`, `SatName`, `TestPhase`, `AcqMode`, `Date`, `Time`, `ConfigName`, `Remark`) 
VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );

-- name: ChangeRemark :exec
Update `Acquisitions` Set `Remark` = ? 
where `Date` like ? and `Time` like ?;

-- name: DeselectTestPhase :exec
UPDATE `TestPhases` SET `Selected` = 0;

-- name: CreateTestPhase :exec
INSERT INTO `TestPhases` 
(`TestPhase`, `Selected`) 
VALUES (?,1);

-- name: SelectTestPhase :exec
UPDATE `TestPhases` SET `Selected` = 1 
WHERE `TestPhase` like ?;

-- name: CheckIfSavedFileExists :one
SELECT COUNT(*) FROM `SavedFiles` WHERE 
`AcquisitionID` = ? AND `FrameID` = ? 
AND `SystemName` = ?;


-- name: CreateSavedFiles :exec
Insert into `SavedFiles` 
(`AcquisitionID`, `FrameID`, `NoOfLines`, `SystemName`, `FilePath`)  
values ( ?, ?, ?, ?, ?);

-- name: RemoveResultProfileParameters :exec
Delete from `ResultProfileParameters` 
where `ProfileID` = ?;

-- name: RemoveResultProfileDetails :exec
Delete from `ResultProfileDetails`
where `ProfileID` = ?;

-- name: RemoveResultProfile :exec
Delete from `ResultProfiles` 
where `ProfileID` = ?;

-- name: CreateResultProfile :execresult
Insert into `ResultProfiles` 
(`ProfileName`, `NoOfResults`) values (?,?);

-- name: CreateResultProfileDetails :exec
Insert into `ResultProfileDetails` 
(`ProfileID`,`ResultID`,`ResultName`,`FrameType`,`FrameID`,`NoOfParameters`,
 `StartLine`,`EndLine`, `StartWord`, `EndWord`, `StartPixel`, `EndPixel` ,
`DecimalOrHex`, `RawOrProcessed`, `AscOrDsc`, `IsMean`, `IsSD`, `StackList` )
values 
(?, ?, ?, ?, ?, ?,
?, ?, ?, ?, ?, ?,
?, ?, ?, ?, ?, ?);


-- name: CreateResultProfileParameters :exec
Insert into `ResultProfileParameters` 
(`ProfileID`, `ResultID`, `ParamID`)
values ( ? , ? , ?);

-- name: ChangeDASIP :exec
Update `DASDetails` 
Set `IPAddress` = ? 
where `DASName` like ?;

CREATE TABLE `ConsolidatedReport` (
  `ConsolidatedReportID` int(11) NOT NULL AUTO_INCREMENT,
  `AcquisitionID` int(11) NOT NULL,
  `SystemName` varchar(45) NOT NULL,
  `FrameType` varchar(45) NOT NULL,
  `FrameID` varchar(45) DEFAULT NULL,
  `AuxDataCorrect` varchar(45) DEFAULT NULL,
  `MaxMean` double DEFAULT NULL,
  `MeanOfMean` double DEFAULT NULL,
  `MaxSD` double DEFAULT NULL,
  `MeanSD` double DEFAULT NULL,
  `RMSSD` double DEFAULT NULL,
  `Path` varchar(100) DEFAULT NULL,
  `NoOfFrames` varchar(45) NOT NULL,
  PRIMARY KEY (`ConsolidatedReportID`),
  FOREIGN KEY (`AcquisitionID`) REFERENCES `Acquisitions` (`AcquisitionID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- name: InsertAuxConsolidatedReport :exec
Insert into `ConsolidatedReport`
(`AcquisitionID`, `SystemName`, `FrameType`, `AuxDataCorrect`, `NoOfFrames`)  
values (?, ?, ?, ?, ?);

-- name: CheckIfConsolidatedSummaryExists :one
Select Count(*) from `ConsolidatedReport`
Where `AcquisitionID` = ? and `SystemName` = ? and `FrameType` = ?;

-- name: UpdateAuxConsolidatedReport :exec
Update `ConsolidatedReport` Set
`AuxDataCorrect` = ?, `NoOfFrames` = ?
Where `AcquisitionID` = ? and `SystemName` = ? and `FrameType` = ?;

-- name: CheckIfVideoConsolidatedSummaryExists :one
Select count(*) from `ConsolidatedReport` 
where `AcquisitionID` = ? and `FrameType` = ? 
and `FrameID`=? and `SystemName` = ?;

-- name: InsertVideoConsolidatedReport :exec
Insert into `ConsolidatedReport`
(`AcquisitionID`, `SystemName`, `FrameType`,`FrameID`, `MaxMean`, 
`MeanOfMean`, `MaxSD`, `MeanSD`, `RMSSD`, `Path`, `NoOfFrames`)
values (?, ?, ?, ?, ?, 
?, ?, ?, ?, ?, ?);

-- name: UpdateVideoConsolidatedReport :exec
UPDATE ConsolidatedReport SET 
`MaxMean` = ?,`MeanOfMean`=?,`MaxSD`=?,`MeanSD`=?,`RMSSD`=?,
`Path`=?,`NoOfFrames`=? 
WHERE `AcquisitionID` = ? and `FrameType` = ? 
and `FrameID` = ? and `SystemName` = ?;

-- name: CreateOfflineResult :exec
Insert into `SavedResults` 
(`UniqueID`,`AcquisitionID`, `SystemName`, `FrameType`, `FrameIdentifier`, `ResultName`, `ResultPath`)
Values ( ? , ?, ?, ?, ?, ?, ?);

-- name: GetMaxResultNumber :many
Select `ResultNumber` from `SavedResults` 
Order by `ResultNumber` desc;