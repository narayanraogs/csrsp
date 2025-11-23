# Result View

The user can view results after processing, be it real-time or offline.

* As soon as the user lands in the Result view screen, he will have a consolidated report visible - this is a special report with complex structure
* All the other reports which are in result profile will come as a stream as and when the server finishes processing them. Once clicked the user can view the report.
* Crucially, the user can keep adding new results to the server as request and the server processes them as First come first serve - this can be changed to worker pool based to process multiple results in common
* The client displays the results based on all the frame types present during the acquisition and the list will be sent from the server based on the Results.go which we have parsed in global. Refer to it if you need more context.
* There are 4 kinds of reports
  * Results - Typically 20-30 rows, 10-15 columns - has error/warning/ok/header for each cell
  * Displays - Typically 10s of thousands of rows and 1000s of columns - no error/warning - useful feature to have - not implemented yet
  * Plots - Interative line charts - multiple parameters in a plot supported - 3d plots for mean and sd also available
  * Histograms - single parameter histograms