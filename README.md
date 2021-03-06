# syncFiles
Syncronize files between several folders on your pc.

Help Page  
Execution: syncFile.exe [Path of Config] (arg) (val)  
-t [value] = time to sleep between checks  
-h / --help = get this page  
[] - necessary | () - optional  
  
Config Format  
"Path from file" -> "Path to copy to (1)" "Path to copy to (2)"...  
Copies the file if it is newer than the file at the destination.  
"Path from file" |-> "Path to copy to (1)" "Path to copy to (2)"...  
Always copy if something in the first file has changed. -")  
"Path to copy to (1)" <-> "Path to copy to (2)" <-> "Path to copy to (3)"  
Takes the last edited file and copies it to all other locations.  
\# - To Comment a line  
