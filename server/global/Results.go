package global

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// --- Data Structures ---

// ResultStructure is the top-level container for result definitions, read from a JSON file.
type ResultStructure struct {
	ResultNames   []FrameTypeResults `json:"resultNames"`
	ResultDetails []ResultProperties `json:"resultDetails"`
}

// FrameTypeResults maps a frame type (e.g., "Aux", "Video") to the kinds of results it can produce.
type FrameTypeResults struct {
	FrameType  string   `json:"frameType"`
	Reports    []string `json:"reports,omitempty"`
	Displays   []string `json:"displays,omitempty"`
	Plots      []string `json:"plots,omitempty"`
	Histograms []string `json:"histograms,omitempty"`
	Media      []string `json:"media,omitempty"`
}

// ResultProperties defines the characteristics and UI options for a specific, named result.
type ResultProperties struct {
	ResultName                 string `json:"resultName"`
	IsFrameType                bool   `json:"isFrameType,omitempty"`
	IsLBT                      bool   `json:"isLBT,omitempty"`
	IsParameterListRequired    bool   `json:"isParameterListRequired,omitempty"`
	IsProcessedOrRawApplicable bool   `json:"isProcessedOrRawApplicable,omitempty"`
	IsDifferenceTypeApplicable bool   `json:"isDifferenceTypeApplicable,omitempty"`
	IsSLELApplicable           bool   `json:"isSLELApplicable,omitempty"`
	IsSPEPApplicable           bool   `json:"isSPEPApplicable,omitempty"`
	IsSortingApplicable        bool   `json:"isSortingApplicable,omitempty"`
	IsDecimalHexApplicable     bool   `json:"isDecimalHexApplicable,omitempty"`
	IsMeanSDApplicable         bool   `json:"isMeanSDApplicable,omitempty"`
	IsBandNumberApplicable     bool   `json:"isBandNumberApplicable,omitempty"`
	IsStackListApplicable      bool   `json:"isStackListApplicable,omitempty"`
	StackType                  string `json:"stackType,omitempty"`
	FilterType                 string `json:"filterType,omitempty"`
}

// LoadResults loads the result definitions from a JSON file and returns a queryable Results object.
func LoadResults(path string) (*Results, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open results file %s: %w", path, err)
	}
	defer file.Close()

	var structure ResultStructure
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&structure); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", path, err)
	}

	return newResultsFromStructure(&structure), nil
}

// newResultsFromStructure processes the raw ResultStructure and builds the lookup maps
// for efficient querying.
func newResultsFromStructure(structure *ResultStructure) *Results {
	r := &Results{
		typesByFrameType:  make(map[string][]string),
		namesByResultType: make(map[string]map[string][]string),
		propsByResultName: make(map[string]ResultProperties),
	}

	// Populate result types and names maps
	for _, result := range structure.ResultNames {
		frameType := strings.ToUpper(result.FrameType)
		r.namesByResultType[frameType] = make(map[string][]string)

		addResultType := func(typeName string, names []string) {
			if len(names) > 0 {
				r.typesByFrameType[frameType] = append(r.typesByFrameType[frameType], typeName)
				r.namesByResultType[frameType][typeName] = names
			}
		}

		addResultType("Reports", result.Reports)
		addResultType("Displays", result.Displays)
		addResultType("Plots", result.Plots)
		addResultType("Histograms", result.Histograms)
		addResultType("Media", result.Media)
	}

	// Populate result properties map
	for _, props := range structure.ResultDetails {
		r.propsByResultName[strings.ToUpper(props.ResultName)] = props
	}

	return r
}

// --- Queryable Results Object ---

// Results provides a queryable interface for the loaded result definitions.
// It is immutable after creation.
type Results struct {
	typesByFrameType  map[string][]string
	namesByResultType map[string]map[string][]string
	propsByResultName map[string]ResultProperties
}

// GetResultTypes returns the available result types (e.g., "Reports", "Plots") for a given frame type.
func (r *Results) GetResultTypes(frameType string) []string {
	return r.typesByFrameType[strings.ToUpper(frameType)]
}

// GetResultNames returns the available result names for a given frame type and result type.
func (r *Results) GetResultNames(frameType, resultType string) []string {
	if names, ok := r.namesByResultType[strings.ToUpper(frameType)]; ok {
		return names[resultType]
	}
	return nil
}

// GetResultProperties returns the properties for a given result name.
// The boolean return value indicates if the result name was found.
func (r *Results) GetResultProperties(resultName string) (ResultProperties, bool) {
	props, ok := r.propsByResultName[strings.ToUpper(resultName)]
	return props, ok
}

// --- Default Configuration Generator ---

// GenerateDefaultResults builds and returns a hard-coded default result structure.
// This function does not write to any files.
func GenerateDefaultResults() *ResultStructure {
	return &ResultStructure{
		ResultNames: []FrameTypeResults{
			{FrameType: "Aux", Reports: []string{"AUX_CONSOLIDATED_REPORT", "AUX_GLOBAL_REPORT", "AUX_ERROR_REPORT_FRAME_TYPE", "AUX_ERROR_REPORT_FRAME_IDENTIFIER"}, Displays: []string{"PARAMETER_DISPLAY", "AUX_RAW_DATA_DISPLAY"}, Plots: []string{"AUX_PARAMETER_PLOT"}},
			{FrameType: "LBT", Reports: []string{"LBT_COMPARSION_REPORT"}},
			{FrameType: "Video", Reports: []string{"DEVICE_SUMMARY_REPORT", "FAULTY_PIXEL_REPORT", "PORTWISE_DETAILED_REPORT", "PORTWISE_SUMMARY", "PAYLOAD_SUMMARY_REPORT", "BANDWISE_DETAILED_REPORT", "BANDWISE_SUMMARY", "PIXELS_WITH_SD_GT_LIMIT"}, Displays: []string{"PIXEL_PERFORMANCE_PIXELNO", "PIXEL_PERFORMANCE_MEAN", "PIXEL_PERFORMACE_SD", "PIXEL_PERFORMACE_DIFF", "ALL_PIXELS_OVER_SEL_LINES", "SEL_PIXELS_OVER_SEL_LINES", "SEL_ODD_PIXELS_OVER_SEL_LINES", "SEL_EVEN_PIXELS_OVER_SEL_LINES", "ALL_SAMPLES_OF_SEL_PIXEL"}, Plots: []string{"MEAN_AND_SD_PLOT", "ROW_WISE_PLOT", "COLUMN_WISE_PLOT", "PORTWISE_PLOT", "BANDWISE_PLOT", "STACKED_ROW_WISE_PLOT", "STACKED_COLUMN_WISE_PLOT", "STACKED_PORTWISE_PLOT", "STACKED_BANDWISE_PLOT", "MEAN_3D_PLOT", "SD_3D_PLOT"}, Histograms: []string{"MEAN_HISTOGRAM", "SD_HISTOGRAM", "SINGLE_PIXEL_HISTOGRAM"}, Media: []string{"MEAN_IMAGE", "IMAGES", "VIDEO"}},
			{FrameType: "Payload", Reports: []string{"SAR_SUMMARY_REPORT"}, Displays: []string{"SSPA_CAL", "LNA_CAL", "TX_RX_CAL", "NOISE_CAL", "PRE_REF_CAL_H", "IMAGING", "POST_REF_CAL_V", "PRE_NOISE_CAL", "POST_NOISE_CAL", "POST_REF_CAL_H", "PRE_REF_CAL_V", "IMAGING"}, Plots: []string{"SSPA_CAL", "LNA_CAL", "TX_RX_CAL", "NOISE_CAL", "PRE_REF_CAL_H", "IMAGING", "POST_REF_CAL_V", "PRE_NOISE_CAL", "POST_NOISE_CAL", "POST_REF_CAL_H", "PRE_REF_CAL_V", "IMAGING"}},
		},
		ResultDetails: []ResultProperties{
			{ResultName: "AUX_CONSOLIDATED_REPORT", IsFrameType: true},
			{ResultName: "AUX_GLOBAL_REPORT", IsFrameType: true},
			{ResultName: "AUX_ERROR_REPORT_FRAME_TYPE", IsFrameType: true},
			{ResultName: "AUX_ERROR_REPORT_FRAME_IDENTIFIER", IsFrameType: false},
			{ResultName: "PARAMETER_DISPLAY", IsParameterListRequired: true, IsSLELApplicable: true, IsProcessedOrRawApplicable: true, IsDecimalHexApplicable: true},
			{ResultName: "AUX_RAW_DATA_DISPLAY", IsSLELApplicable: true, IsDecimalHexApplicable: true, FilterType: "WORD"},
			{ResultName: "AUX_PARAMETER_PLOT", IsParameterListRequired: true, IsProcessedOrRawApplicable: true, IsDifferenceTypeApplicable: true, IsSLELApplicable: true},
			{ResultName: "LBT_COMPARSION_REPORT", IsFrameType: true, IsLBT: true},
			{ResultName: "DEVICE_SUMMARY_REPORT"},
			{ResultName: "FAULTY_PIXEL_REPORT", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "PORTWISE_DETAILED_REPORT"},
			{ResultName: "PORTWISE_SUMMARY"},
			{ResultName: "PAYLOAD_SUMMARY_REPORT"},
			{ResultName: "BANDWISE_DETAILED_REPORT"},
			{ResultName: "BANDWISE_SUMMARY"},
			{ResultName: "PIXELS_WITH_SD_GT_LIMIT", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "PIXEL_PERFORMANCE_PIXELNO", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "PIXEL_PERFORMANCE_MEAN", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "PIXEL_PERFORMACE_SD", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "PIXEL_PERFORMACE_DIFF", IsSPEPApplicable: true, IsSortingApplicable: true, FilterType: "PIXEL"},
			{ResultName: "ALL_PIXELS_OVER_SEL_LINES", IsSLELApplicable: true, IsDecimalHexApplicable: true},
			{ResultName: "SEL_PIXELS_OVER_SEL_LINES", IsSLELApplicable: true, IsSPEPApplicable: true, FilterType: "PIXEL", IsDecimalHexApplicable: true},
			{ResultName: "SEL_ODD_PIXELS_OVER_SEL_LINES", IsSLELApplicable: true, IsSPEPApplicable: true, FilterType: "PIXEL", IsDecimalHexApplicable: true},
			{ResultName: "SEL_EVEN_PIXELS_OVER_SEL_LINES", IsSLELApplicable: true, IsSPEPApplicable: true, FilterType: "PIXEL", IsDecimalHexApplicable: true},
			{ResultName: "ALL_SAMPLES_OF_SEL_PIXEL", IsSPEPApplicable: true, FilterType: "PIXEL", IsDecimalHexApplicable: true},
			{ResultName: "MEAN_AND_SD_PLOT", IsMeanSDApplicable: true},
			{ResultName: "ROW_WISE_PLOT", IsMeanSDApplicable: true},
			{ResultName: "COLUMN_WISE_PLOT", IsMeanSDApplicable: true},
			{ResultName: "PORTWISE_PLOT", IsMeanSDApplicable: true},
			{ResultName: "BANDWISE_PLOT", IsMeanSDApplicable: true},
			{ResultName: "STACKED_ROW_WISE_PLOT", IsMeanSDApplicable: true, IsStackListApplicable: true, StackType: "ROW"},
			{ResultName: "STACKED_COLUMN_WISE_PLOT", IsMeanSDApplicable: true, IsStackListApplicable: true, StackType: "COLUMN"},
			{ResultName: "STACKED_PORTWISE_PLOT", IsMeanSDApplicable: true, IsStackListApplicable: true, StackType: "PORT"},
			{ResultName: "STACKED_BANDWISE_PLOT", IsMeanSDApplicable: true, IsStackListApplicable: true, StackType: "BAND"},
			{ResultName: "MEAN_3D_PLOT"},
			{ResultName: "SD_3D_PLOT"},
			{ResultName: "MEAN_HISTOGRAM"},
			{ResultName: "SD_HISTOGRAM"},
			{ResultName: "SINGLE_PIXEL_HISTOGRAM", IsStackListApplicable: true, StackType: "PIXEL"},
			{ResultName: "MEAN_IMAGE"},
			{ResultName: "IMAGES"},
			{ResultName: "VIDEO"},
			{ResultName: "SAR_SUMMARY_REPORT"},
			{ResultName: "SSPA_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "LNA_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "TX_RX_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "NOISE_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "PRE_REF_CAL_H", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "POST_REF_CAL_H", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "PRE_REF_CAL_V", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "POST_REF_CAL_V", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "PRE_NOISE_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "POST_NOISE_CAL", IsParameterListRequired: true, IsSLELApplicable: true},
			{ResultName: "IMAGING", IsParameterListRequired: true, IsSLELApplicable: true},
		},
	}
}
