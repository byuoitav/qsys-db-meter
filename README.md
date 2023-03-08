# qsys-db-meter
Service to facilitate JSON file upload to a server which can then be used as a data source in scoreboard video software suites

Creates a MeterData.json file and continually updates this file with each frame of data downloaded from the Q-Sys designer file.  This is the file to be read by the video server software as the data stream.

Both LKFS and dBSPL are transmitted and stored in this JSON file.  Additionally, a splMax and lkfsMax is calculated and transmitted/stored.

Also creates a .ipInfo.dll file to store the last IP and port the program used. 
