# loopit
Loop through combinations of parameters while executing a command on each permutation.

Flags:
  -c string
        Command to execute. Use MYFILE1, MYFILE2, ..., MYFILE9 as wildcards
  -l string
        Log progress to file with this name.
  -n integer
        number of simultaneous commands (default 1)
  -p  boolean
        Pretend. List commands that would have been run.

Special argument prefixes:
list: - a comma separated list of parameters
glob: - a string representing a file glob
stdin - read parameter list from stdin

Special variable:
MYJOB - an integer of the job number

Example:
loopit -p -c 'echo MYJOB MYFILE1 MYFILE2 MYFILE3' glob:*.txt list:1,2,3 file1.txt
echo 1 file1.txt 1 file1_line1
echo 2 file1.txt 1 file1_line2
echo 3 file1.txt 2 file1_line1
echo 4 file1.txt 2 file1_line2
echo 5 file1.txt 3 file1_line1
echo 6 file1.txt 3 file1_line2
echo 7 file2.txt 1 file1_line1
echo 8 file2.txt 1 file1_line2
echo 9 file2.txt 2 file1_line1
echo 10 file2.txt 2 file1_line2
echo 11 file2.txt 3 file1_line1
echo 12 file2.txt 3 file1_line2

