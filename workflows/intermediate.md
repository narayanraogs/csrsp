# Intermediate Reports

Intermediate Reports are the reports which are available for the client during the execution of a pipeline. 

There are 3 Categories
* If the frame is an Aux type, diplay the below status for each parameter
  * ParameterName: Value - Color this line Red if error, amber if warning
  * Typically in a project there will be 10-12 aux frames with around 20 parameters per frame.
  * However, in extreme cases we have seen 200 parameters per frame - only 2 frames - other frames had 20 parameters only
* If the frame is Optical type, display the below status
  * Mean Of Mean: value
  * Mean of SD: Value
  * Max SD: Value
  * Red if any thing is out of limits
* If the frame is microwave, display the below status
  * No of pulses acquired for this till now
* Other than this if the data is being acquired, then dispaly a progress bar of each chain. For acquisition there is a minimum of 1 chain and max of 6 chains.
* If the data is being offline processed, then display an indeterminate progress bar. - It would be good to have a determinate progress bar for processing also. But as the number of steps can reach to 1000 or more, its difficult to determine as of now.
* Once the processing is complete the user will be transferred to the result display screen whcih will be a seperate workflow.
  