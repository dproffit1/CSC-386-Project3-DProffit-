package main

import (
	"Project2Demo/FileSystem"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin) // outside the loop, only need one reader
	FileSystem.InitializeFileSystem()

	fmt.Println("\nPlease type in a command: 'mkdir', 'cp', 'mv', 'rm', 'more', '>>', or 'exit'")

	for {
		input, err := reader.ReadString('\n') //get user command

		if err != nil { // check input error
			log.Fatal("Failed to read input:", err)
		}

		input = strings.TrimSpace(input) // remove new lines

		if input == "" { //check for no input (prevent string pars crash)
			fmt.Print("Command not recognized")
			continue
		}

		inputArray := strings.Fields(input) // split inputs by spaces for args
		command := inputArray[0]            // command is its own var
		commandArgs := inputArray[1:]       // args are its own array

		switch command {
		case "whoami": // prints name
			fmt.Println("Dakotah Proffit, dproffit1")
		case "exit": // exit program
			os.Exit(0)
		case "more": // view text in sections
			moreCommand(commandArgs[0])
		case "mkdir": // make directory
			mkdirCommand(commandArgs[0])
		case "cp": // copy
			copyCommand(commandArgs[0], commandArgs[1])
		case "mv": // moves
			moveCommand(commandArgs[0], commandArgs[1])
		case "rm": // remove (only for files)
			removeCommand(commandArgs[0])
		case "help":
			fmt.Println("command: 'mkdir', 'cp', 'mv', 'rm', 'more', '>>', or 'exit'")
		default: // for commands not in program
			osCommand(command, commandArgs)
		}
	}
}

func getParentandChildInodes(path string) (parentinode FileSystem.INode, childinode FileSystem.INode, parentinodenum int, childinodenum int) {
	stringSlice := strings.Split(path, "/")
	newDirectory := stringSlice[len(stringSlice)-1]
	stringSlice = stringSlice[:len(stringSlice)-1]
	var toPath string
	for _, dir := range stringSlice {
		if dir != "" {
			toPath = toPath + "/" + dir
		}
	}
	parentinode, parentinodenum = FileSystem.FindSubdirectories(toPath)
	childinode, childinodenum = FileSystem.Open(FileSystem.CREATE, newDirectory, parentinode)
	return parentinode, childinode, parentinodenum, childinodenum
}

func moveCommand(frompath string, topath string) { //lazy function
	copyCommand(frompath, topath)
	removeCommand(frompath)
}

func copyCommand(frompath string, topath string) {
	_, movinginode, _, _ := getParentandChildInodes(frompath)
	filecontent := FileSystem.Read(&movinginode)
	toinode, _, toinodenum, _ := getParentandChildInodes(topath)
	inputcontent := []byte(filecontent)
	FileSystem.Write(&toinode, toinodenum, inputcontent)
}

func moreCommand(path string) {
	_, childinode, _, _ := getParentandChildInodes(path)
	if childinode.DirectBlock1 == 0 {
		fmt.Println("nothing to read in")
	} else {
		filecontent := FileSystem.Read(&childinode)
		fmt.Println(filecontent)
	}
}

func mkdirCommand(path string) {
	_, childinode, parentinodenum, childinodenum := getParentandChildInodes(path)
	directoryBlock, childinode := FileSystem.CreateDirectoryFile(parentinodenum, childinodenum)
	bytesForDirectoryBlock := FileSystem.EncodeToBytes(directoryBlock)
	FileSystem.Write(&childinode, childinodenum, bytesForDirectoryBlock)
	fmt.Println("New directory has been made!")
}

func removeCommand(path string) {
	parentinode, childinode, _, childinodenum := getParentandChildInodes(path)
	if childinode.DirectBlock1 == 0 {
		fmt.Println("nothing to remove")
	} else {
		FileSystem.Unlink(childinodenum, parentinode)
		fmt.Println("File has been removed!")
	}
}

func osCommand(command string, args []string) {
	modifiedArgs := []string{}
	var nextOutputFile string
	pathInOutputFile := false

	// Iterate through arguments to handle redirection and check for paths
	for i, arg := range args {
		if arg == ">>" {
			if i+1 < len(args) {
				nextOutputFile = args[i+1] // Store the next argument after '>>'
				if strings.Contains(nextOutputFile, "/") {
					pathInOutputFile = true // Set flag if '/' is found
				}
			}
			break // Stop processing arguments after '>>'
		}
		modifiedArgs = append(modifiedArgs, arg) // Add argument if not part of redirection
	}

	inputFileContent, err := os.ReadFile(modifiedArgs[len(modifiedArgs)-1])
	if err != nil {
		fmt.Println("couldn't read in file")
		fmt.Println(err)
	}

	// Output relevant information or take actions based on flags here, if necessary
	if pathInOutputFile {
		stringSlice := strings.Split(nextOutputFile, "/")
		fileName := stringSlice[len(stringSlice)-1]
		parentinode, _, _, _ := getParentandChildInodes(nextOutputFile)
		newFileInode, firstInodeNun := FileSystem.Open(FileSystem.CREATE, fileName, parentinode)
		contentToWrite := []byte(inputFileContent)
		FileSystem.Write(&newFileInode, firstInodeNun, contentToWrite)
		fmt.Println("file read in")
	} else {
		newFileInode, firstInodeNun := FileSystem.Open(FileSystem.CREATE, nextOutputFile, FileSystem.RootFolder)
		contentToWrite := []byte(inputFileContent)
		FileSystem.Write(&newFileInode, firstInodeNun, contentToWrite)
		fmt.Println("file read in")
	}

}
