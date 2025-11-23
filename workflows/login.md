# Login

* The user has to login to CSRSP before using the software
* There are 2 kinds of login.
  * Whitelisted
  * Autenticated
* There are logically different sections of the client which the user can access
  * Acquisition - Acquire data or BER
  * Offline Processing - Process already acquired data
  * Data Transfer - Download the files for selected previous acquisitions for which the user has access
  * Database Admin - Restricted database tasks are avalable to the user 
    * changing the remark on the acquired datasets
    * change the current test phase of the project
    * Change the IP address from which the user can acquire data
  * Result Profile Editing - Edit the result profiles in the database
  * File Processing - Process the data from a file - Not previous acquired - Can be considered admin access
  * Trend Analysis - Analise multiple results together
  * Developer Options - Change how the software behaves, for example - change the encryption mode or change the last step of pipeline etc... - Can be considered admin access
* When the user is using autenticated login, the permissions list is sent to the user based on which the client will display only the allowed fields
* If it is a whitelisted IP all the features are enabled
* Session should persist on refresh of the webpage also.
