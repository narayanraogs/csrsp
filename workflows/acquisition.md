# Acquisition 

* This screen is visible to users having acquisition privilages
* Initially show the following things to the client
  * Acquisition modes 
    * multiple options are available from server 
    * client has to select one of them
  * Payloads 
    * multiple options are available from server 
    * client has to select one of them
  * Config Names
    * multiple options are available from server 
    * client has to select one of them
    * Options depend on the Acquisition mode selected
  * Acquisition type
    * Frame Based - If user selects this No of frames to be populated
    * Time Based - If user selects this time to be populated
    * Till Stopped - No specific user input required for this
  * Result Profiles
    * multiple options are available from server 
    * client has to select one of them
  * Remarks - Any remark the user wants to give to the acquisition
* Work flow
  * The user starts the acquistion
  * If any one hardware are ready, then go ahead - other wise negative ack
  * During acquisition the user should get updates of the current number of frame acquired as well as the data integrity in report format from the server store
  * On completion of acquisition
    * Show consolidated report
    * Also generate all the results in the result profile and show the name of generated results to the user - to be opened on click
    * If there is any error - show consolidated error report
    * If auto archival is selected - archive the data
    * If auto archival is not selected - when the user navigates to some other screen other than result view - prompt for manual archival
  * The results to be shown during acquisition are called intermedite reports.
