CREATE TABLE `AcquisitionMode` (
  `AcqID` int(11) NOT NULL AUTO_INCREMENT,
  `SatName` varchar(45) NOT NULL,
  `AcqMode` varchar(45) NOT NULL,
  `AcqType` enum('Acquisition','BER') NOT NULL DEFAULT 'Acquisition',
  PRIMARY KEY (`AcqID`),
  UNIQUE KEY `SatName_AcqMode` (`SatName`, `AcqMode`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `Acquisitions` (
  `AcquisitionID` int(11) NOT NULL AUTO_INCREMENT,
  `SystemName` varchar(45) NOT NULL,
  `SatName` varchar(45) NOT NULL,
  `TestPhase` varchar(45) NOT NULL,
  `AcqMode` varchar(45) NOT NULL,
  `Date` varchar(45) NOT NULL,
  `Time` varchar(45) NOT NULL,
  `ConfigName` varchar(45) NOT NULL,
  `Remark` TEXT NOT NULL,
  PRIMARY KEY (`AcquisitionID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `AnalogType` (
  `ParamID` varchar(45) NOT NULL,
  `Equation` varchar(200) NOT NULL,
  `DataType` enum('int_32','int_64','float_ieee_32','float_ieee_64','float_1750a_32','float_1750a_48','int_16','uint_16','uint_32','uint_64') NOT NULL,
  `VariableName` varchar(45) NOT NULL,
  `ParameterIDs` varchar(250) DEFAULT NULL,
  `ValueValidationRequired` BOOLEAN NOT NULL DEFAULT 0,
  `DifferenceValidationRequired` BOOLEAN NOT NULL DEFAULT 0,
  `LowerLimitValue` double DEFAULT NULL,
  `UpperLimitValue` double DEFAULT NULL,
  `ToleranceValue` double DEFAULT NULL,
  `DifferenceValue` double DEFAULT NULL,
  `DifferenceTolerance` double DEFAULT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `BERLogging` (
  `BerID` int(11) NOT NULL AUTO_INCREMENT,
  `AcqMode` varchar(45) NOT NULL,
  `DASID` int(11) NOT NULL,
  `ComponentCode` int(11) NOT NULL,
  `FrameID` int(11) NOT NULL,
  PRIMARY KEY (`BerID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `BitDetails` (
  `BitID` int(11) NOT NULL AUTO_INCREMENT,
  `ParamID` varchar(45) NOT NULL,
  `WordNo` int(11) NOT NULL,
  `ValidBits` varchar(45) NOT NULL,
  PRIMARY KEY (`BitID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Configuration` (
  `ConfigID` int(11) NOT NULL AUTO_INCREMENT,
  `SatName` varchar(45) NOT NULL,
  `AcqMode` varchar(45) NOT NULL,
  `ConfigName` varchar(45) NOT NULL,
  `Logic` varchar(45) NOT NULL,
  PRIMARY KEY (`ConfigID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

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

CREATE TABLE `DASDetails` (
  `DASID` int(11) NOT NULL AUTO_INCREMENT,
  `DASName` varchar(45) NOT NULL,
  `DASType` enum('HDR','PDPS','TTCP','CDP') NOT NULL,
  `IPAddress` varchar(45) NOT NULL,
  `TimeOut` int(11) NOT NULL,
  PRIMARY KEY (`DASID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `DASPath` (
  `DASPathID` int(11) NOT NULL AUTO_INCREMENT,
  `AcqMode` varchar(45) NOT NULL,
  `ChainNumber` int(11) NOT NULL,
  `DASID` int(11) NOT NULL,
  `DPUNumber` int(11) NOT NULL,
  `dataFrameid` int(11) NOT NULL,
  `headerFrameId` int(11) NOT NULL,
  `trailerFrameId` int(11) NOT NULL,
  PRIMARY KEY (`DASPathID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `CRCType` (
  `ParamID` varchar(45) NOT NULL,
  `StartWord` int(11) NOT NULL,
  `NoOfWords` int(11) NOT NULL,
  `CRCPolynomial` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `DASComponentDetails` (
  `ComponentID` int(11) NOT NULL AUTO_INCREMENT,
  `ComponentName` varchar(200) NOT NULL,
  `ParentCompName` varchar(250) DEFAULT NULL,
  `ComponentCode` varchar(45) NOT NULL,
  `TotalLengthinWords` int(11) NOT NULL,
  `ConfigLengthinWords` int(11) NOT NULL,
  `ExcludedWords` varchar(500) DEFAULT NULL,
  `StartWord` int(11) NOT NULL,
  `StartBit` int(11) NOT NULL,
  `NoOfBits` int(11) NOT NULL,
  `Value` int(11) NOT NULL,
  `Enabled` BOOLEAN NOT NULL,
  PRIMARY KEY (`ComponentID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `DASConfigurations` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AcquisitionMode` varchar(45) NOT NULL,
  `DASID` int(11) NOT NULL,
  `ComponentCode` varchar(45) NOT NULL,
  `SetParameters` json NOT NULL,
  `VerifyParameters` json NOT NULL,
  `ExecutionOrder` int(11) NOT NULL,
  `DelayAfterMS` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  FOREIGN KEY (`DASID`) REFERENCES `DASDetails` (`DASID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `DataDownloadPermissions` (
  `UserID` int(11) NOT NULL,
  `FrameID` int(11) NOT NULL,
  PRIMARY KEY (`UserID`, `FrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `Frames` (
  `FrameID` int(11) NOT NULL AUTO_INCREMENT,
  `FrameType` varchar(255) NOT NULL,
  `FrameIdentifier` varchar(255) NOT NULL,
  `FrameFileName` varchar(255) NOT NULL,
  `FrameLength` int(11) NOT NULL,
  `IdentifyingValue` int(11) DEFAULT -1,
  PRIMARY KEY (`FrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `FrameTypes` (
  `FrameTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `SatName` varchar(255) NOT NULL,
  `FrameType` varchar(255) NOT NULL,
  `Type` enum('Aux','Optical','None','LBT','Image','RS', 'Microwave') NOT NULL DEFAULT 'None',
  `ProcessedOrAcquired` enum('Processed','Acquired') NOT NULL DEFAULT 'Processed',
  `ToBeWritten` BOOLEAN NOT NULL,
  `ToBeArchived` BOOLEAN NOT NULL,
  `IdentifyingParam` varchar(45) NOT NULL DEFAULT '',
  `FrameLengthParam` varchar(45) NOT NULL DEFAULT '',
  `isIntermediateReport` BOOLEAN NOT NULL,
  PRIMARY KEY (`FrameTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `FSCType` (
  `ParamID` varchar(45) NOT NULL,
  `FSC` varchar(100) NOT NULL,
  `NoOfBitErrorAllowed` int(11) NOT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `HKDetails` (
  `HKDetailsId` int(11) NOT NULL AUTO_INCREMENT,
  `SCCName` varchar(50) NOT NULL,
  `UserName` varchar(50) DEFAULT NULL,
  `Password` varchar(255) NOT NULL,
  `IPAddress` varchar(50) NOT NULL,
  `SatName` varchar(50) NOT NULL,
  `StationID` varchar(50) NOT NULL,
  `ChainName` varchar(50) NOT NULL,
  `Path` varchar(150) NOT NULL,
  `FrameLength` int(11) DEFAULT NULL,
  `StartWord` int(11) DEFAULT NULL,
  `NoOfWords` int(11) NOT NULL,
  PRIMARY KEY (`HKDetailsId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `IncrementType` (
  `ParamID` varchar(45) NOT NULL,
  `Increment` int(11) NOT NULL,
  `StartValue` int(11) NOT NULL DEFAULT -1,
  `EndValue` int(11) NOT NULL DEFAULT -1,
  `Tolerance` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ParameterDetails` (
  `ParamID` varchar(45) NOT NULL,
  `SatName` varchar(45) NOT NULL,
  `ParameterName` varchar(45) NOT NULL,
  `ParameterType` enum('Radix','FSC','Status','Analog','Increment','CRC') NOT NULL,
  `FrameType` varchar(45) NOT NULL,
  `BitWise` tinyint(1) NOT NULL,
  `StartWord` int(11) NOT NULL,
  `StartBit` int(11) NOT NULL,
  `NoOfBits` varchar(45) NOT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Payloads` (
  `PayloadID` int(11) NOT NULL AUTO_INCREMENT,
  `PayloadName` varchar(45) NOT NULL,
  `NoOfDevices` int(11) NOT NULL DEFAULT 1,
  PRIMARY KEY (`PayloadID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `RadixType` (
  `ParamID` varchar(45) NOT NULL,
  `Radix` int(11) NOT NULL,
  `ValueValidationRequired` BOOLEAN NOT NULL DEFAULT 0,
  `DifferenceValidationRequired` BOOLEAN NOT NULL DEFAULT 0,
  `UpperLimitValue` int(11) DEFAULT NULL,
  `LowerLimitValue` int(11) DEFAULT NULL,
  `ValueTolerance` int(11) DEFAULT NULL,
  `DifferenceValue` int(11) DEFAULT NULL,
  `DifferenceTolerance` int(11) DEFAULT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE `RemoteParameters` (
  `ParamID` bigint(15) NOT NULL AUTO_INCREMENT,
  `ParamType` varchar(40) NOT NULL,
  `ParamName` varchar(40) NOT NULL,
  `MaxLength` int(11) NOT NULL,
  `Datatype` enum('string','int') NOT NULL,
  PRIMARY KEY (`ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `ResultProfileDetails` (
  `ProfileID` int(11) NOT NULL,
  `ResultID` int(11) NOT NULL,
  `ResultName` varchar(45) NOT NULL,
  `FrameType` varchar(45) NOT NULL,
  `FrameID` varchar(40) NOT NULL,
  `NoOfParameters` int(11) NOT NULL,
  `StartLine` int(11) NOT NULL,
  `EndLine` int(11) NOT NULL,
  `StartWord` int(11) NOT NULL,
  `EndWord` int(11) NOT NULL,
  `StartPixel` int(11) NOT NULL,
  `EndPixel` int(11) NOT NULL,
  `DecimalOrHex` tinyint(1) NOT NULL,
  `RawOrProcessed` tinyint(1) NOT NULL,
  `AscOrDsc` tinyint(1) NOT NULL,
  `IsMean` tinyint(1) NOT NULL,
  `IsSD` tinyint(1) NOT NULL,
  `StackList` varchar(100) NOT NULL,
  PRIMARY KEY (`ProfileID`, `ResultID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `ResultProfileParameters` (
  `ProfileID` int(11) NOT NULL,
  `ResultID` int(11) NOT NULL,
  `ParamID` varchar(45) NOT NULL,
   PRIMARY KEY (`ProfileID`, `ResultID`, `ParamID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `ResultProfiles` (
  `ProfileID` int(11) NOT NULL AUTO_INCREMENT,
  `ProfileName` varchar(30) NOT NULL,
  `NoOfResults` int(11) NOT NULL,
  PRIMARY KEY (`ProfileID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



CREATE TABLE `SARVideoDataProcessing` (
  `ProcessID` int(11) NOT NULL AUTO_INCREMENT,
  `StartByte` int(11) NOT NULL,
  `PayloadI1FrameID` int(11) DEFAULT NULL,
  `PayloadI2FrameID` int(11) DEFAULT NULL,
  `PayloadQ1FrameID` int(11) DEFAULT NULL,
  `PayloadQ2FrameID` int(11) DEFAULT NULL,
  `PolarizationPID` varchar(45) NOT NULL,
  `TimingStatePID` varchar(45) NOT NULL,
  `BAQPID` varchar(45) NOT NULL,
  `DataWindowPID` varchar(45) NOT NULL,
  `ChirpBWPID` varchar(45) NOT NULL,
  `SARModePID` varchar(45) NOT NULL,
  `BAQBlockSize` int(11) NOT NULL,
  `SARModeValue` varchar(20) NOT NULL,
  `polValue` varchar(20) NOT NULL,
  `TimingStateToExclude` varchar(30) NOT NULL,
  PRIMARY KEY (`ProcessID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `SavedFiles` (
  `AcquisitionID` int(11) NOT NULL,
  `FrameID` int(11) NOT NULL,
  `NoOfLines` int(11) NOT NULL,
  `SystemName` varchar(45) NOT NULL,
  `FilePath` varchar(255) NOT NULL,
  PRIMARY KEY (`AcquisitionID`, `FrameID`, `SystemName`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `SavedResults` (
  `ResultNumber` int(11) NOT NULL AUTO_INCREMENT,
  `UniqueID` blob NOT NULL,
  `AcquisitionID` int(11) NOT NULL,
  `SystemName` varchar(45) NOT NULL,
  `FrameType` varchar(45) NOT NULL,
  `FrameIdentifier` varchar(45) DEFAULT NULL,
  `ResultName` varchar(45) NOT NULL,
  `ResultPath` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ResultNumber`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `StatusType` (
  `ParamID` varchar(45) NOT NULL,
  `BitValue` varchar(45) NOT NULL,
  `Interpretation` varchar(45) NOT NULL,
  `MultipleValuesAllowed` BOOLEAN NOT NULL,
  PRIMARY KEY (`ParamID`,`BitValue`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `TestPhases` (
  `TestPhaseID` int(11) NOT NULL AUTO_INCREMENT,
  `TestPhase` varchar(45) NOT NULL,
  `Selected` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`TestPhaseID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `Users` (
  `UserID` int(11) NOT NULL AUTO_INCREMENT,
  `UserName` varchar(45) NOT NULL,
  `Password` varchar(255) NOT NULL,
  `Permissions` text NOT NULL,
  PRIMARY KEY (`UserID`),
  UNIQUE KEY `UserName_UNIQUE` (`UserName`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `VideoDataProcessing` (
  `FrameID` int(11) NOT NULL,
  `StartByte` int(11) NOT NULL,
  `NoOfPixelsPerLine` int(11) NOT NULL,
  `NoOfLinesPerFrame` int(11) NOT NULL,
  `NoOfBitsPerPixelPayload` int(11) NOT NULL,
  `NoOfBitsPerPixelCheckout` int(11) NOT NULL,
  PRIMARY KEY (`FrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `VideoSubFrames` (
  `SubFrameID` int(11) NOT NULL,
  `FrameID` int(11) NOT NULL,
  `SubFrameType` enum('Band','Port','','') NOT NULL,
  `SubFrameName` varchar(45) NOT NULL,
  `SubFrameNumber` int(11) NOT NULL,
  `StartPixel` int(11) NOT NULL,
  `EndPixel` int(11) NOT NULL,
  `PixelAccessType` enum('Odd','Even','All','') NOT NULL,
  PRIMARY KEY (`SubFrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


CREATE TABLE `VideoFaultyPixels` (
  `FrameID` int(11) NOT NULL,
  `NoOfFaultyPixels` int(11) NOT NULL,
  `PixelNumbers` text NOT NULL,
  PRIMARY KEY (`FrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


CREATE TABLE `VideoLimits` (
  `FrameID` int(11) NOT NULL,
  `MaxSD` double NOT NULL,
  `MaxMean` double NOT NULL, 
  `CountDifference` double NOT NULL, 
  PRIMARY KEY (`FrameID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

