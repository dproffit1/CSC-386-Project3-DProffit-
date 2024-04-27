In this project, there are 9 supported functions:

whiami, prints author's name
help, list avalible functions
mkdir, will make a directory in the virtual file system
>>, can be used with cat to read in text from a file on the actual OS to the virtual OS
rm, can remove a file
more, read a file
exit, closes virtual disk enviorment
       /!\These commands currently dont work/!\
cp, to copy a file from one location to another
mv, moves a file from one location to another

when inputing location, use / to seporate files
/ is the home (bottom) directory so /test/dir would be root->test->dir
the cp and mv commands dont work to do a slicing error in the write function, I don't see
any reason why this should be happening
