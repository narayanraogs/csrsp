# BER Logging 

* This screen is visible to users having acquisition privilages
* Initially show the following things to the client
  * Acquisition modes 
    * multiple options are available from server 
    * client has to select one of them
  * Remarks - Any remark the user wants to give to the acquisition
* Work flow
  * The user starts the Ber Logging
  * The server starts the acquistion and processing. 
  * There can be a maximum of 6 streams of acquisition
  * For each stream there can be a max of 6 parameters to be processed
  * The requirement is to have a real-time plot of all the parameters till BER Logging is stopped by the user
  * The server will also send the present value for each parameter
  * The maximum Logging time to be accounted for is 2 hours and the sampling rate is minimum of 256 milliseconds - but commonly 1 second
* Once the BER Logging is stopped the acquired data is to be archived.
* There is a related BER Offline Viewer workflow
