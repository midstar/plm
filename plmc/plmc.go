package main

import (
	"flag"
	"fmt"
	"os"
)

// PLMUrl to PLM server (daemon)
var PLMUrl string

// Matcher matcher -m flag
var Matcher string

// UIDs uid -u flag
var UIDs string

// FailLimit memory fail limit -f flag
var FailLimit int64

func printUsage() {
	fmt.Printf("usage: plmc [options] <command> [<args>]\n\n")
	fmt.Printf(" General options:\n")
	fmt.Printf("  -h   Help (overview)\n")
	fmt.Printf("  -v   Display version\n")
	fmt.Printf("  -a   PLM server (daemon) URL address. Default http://localhost:12124\n")
	fmt.Printf("\n Commands:\n")
	fmt.Printf("  help   Help for a command\n")
	fmt.Printf("  plot   Download plot for one or more processes\n")
	fmt.Printf("  info   List info about one or more processes\n")
	fmt.Printf("  maxmem Display max memory used by process\n")
	fmt.Printf("  minmem Display min memory used by process\n")
	fmt.Printf("  tagset Create a tag\n")
	fmt.Printf("  tagget Get a tag\n")
	fmt.Printf("  tags   List all tags\n")
}

func printUsageCommand(command string) {
	switch command {
	case "help":
		fmt.Printf("Get help for a command\n\n")
		fmt.Printf("Usage: plmc help <command>\n")
	case "plot":
		fmt.Printf("Plot memory usage of processes.\n")
		fmt.Printf("By default all processes are plotted. Can be resttricted\n")
		fmt.Printf("with options described below\n\n")
		fmt.Printf("Usage: plmc [options] plot <filename>\n\n")
		fmt.Printf(" Options:\n")
		printProcessFilterFlags()
	case "info":
		fmt.Printf("List process info.\n")
		fmt.Printf("By default all processes are listed. Can be resttricted\n")
		fmt.Printf("with options described below\n\n")
		fmt.Printf("Usage: plmc [options] info\n\n")
		fmt.Printf(" Options:\n")
		printProcessFilterFlags()
	case "maxmem":
		fmt.Printf("Display max memory used by process.\n")
		fmt.Printf("By default all processes are listed. Can be resttricted\n")
		fmt.Printf("with options described below\n\n")
		fmt.Printf("Usage: plmc [options] maxmem\n\n")
		fmt.Printf(" Options:\n")
		printProcessFilterFlags()
		fmt.Printf("  -f <int>      Fail (return code 1) if memory is above the\n")
		fmt.Printf("                specified value in KB. If more than one\n")
		fmt.Printf("                process match the command will fail if any\n")
		fmt.Printf("                of the processes max memory is over the limit\n")
	case "minmem":
		fmt.Printf("Display min memory used by process.\n")
		fmt.Printf("By default all processes are listed. Can be resttricted\n")
		fmt.Printf("with options described below\n\n")
		fmt.Printf("Usage: plmc [options] minmem\n\n")
		fmt.Printf(" Options:\n")
		printProcessFilterFlags()
		fmt.Printf("  -f <int>      Fail (return code 1) if memory is below the\n")
		fmt.Printf("                specified value in KB. If more than one\n")
		fmt.Printf("                process match the command will fail if any\n")
		fmt.Printf("                of the processes min memory is below the limit\n")
	case "tagset":
		fmt.Printf("Create a tag (i.e. a named timestamp).\n\n")
		fmt.Printf("Usage: plmc tagset <tagname>\n")
	case "tagget":
		fmt.Printf("Get the timestampe of a tag.\n\n")
		fmt.Printf("Usage: plmc tagget <tagname>\n")
	case "tags":
		fmt.Printf("List all tags.\n\n")
		fmt.Printf("Usage: plmc tags\n\n")
	default:
		fmt.Fprintf(os.Stderr, "No such command %s\n\n", command)
		printUsage()
	}
}

func printProcessFilterFlags() {
	fmt.Printf("  -m <string>   List all processes matching the string. To insert\n")
	fmt.Printf("                space surround your match with \". For example\n")
	fmt.Printf("                -m \"myprocess.exe -param 3\"\n")
	fmt.Printf("  -u <int>      List process with matching UID. More than\n")
	fmt.Printf("                one UID can be entered in a comma separated\n")
	fmt.Printf("                list. For example -u 3,6,7\n")
}

func invalidUsage(why string) {
	fmt.Fprintf(os.Stderr, why)
	fmt.Fprintf(os.Stderr, "\n\n")
	printUsage()
	os.Exit(1)
}

func invalidUsageCommand(why string, command string) {
	fmt.Fprintf(os.Stderr, why)
	fmt.Fprintf(os.Stderr, "\n\n")
	printUsageCommand(command)
	os.Exit(1)
}

func main() {
	var version = flag.Bool("v", false, "Display version")
	flag.StringVar(&PLMUrl, "a", "http://localhost:12124", "PLM server address")
	flag.StringVar(&Matcher, "m", "", "Matcher(s)")
	flag.StringVar(&UIDs, "u", "", "UID(s)")
	flag.Int64Var(&FailLimit, "f", -1, "Fail limit")
	flag.Usage = printUsage
	flag.Parse()

	if *version {
		fmt.Printf("Client Version:    %s\n", applicationVersion)
		fmt.Printf("Client Build Time: %s\n", applicationBuildTime)
		fmt.Printf("Client GIT Hash:   %s\n", applicationGitHash)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		invalidUsage("You need to provide a command!")
	}

	command := flag.Arg(0)
	var err error
	switch command {
	case "help":
		if flag.NArg() != 2 {
			invalidUsageCommand(fmt.Sprintf("help takes 1 argument but %d given!", flag.NArg()-1), command)
		}
		printUsageCommand(flag.Arg(1))
	case "plot":
		if flag.NArg() != 2 {
			invalidUsageCommand(fmt.Sprintf("plot takes 1 argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdPlot(flag.Arg(1))
	case "info":
		if flag.NArg() != 1 {
			invalidUsageCommand(fmt.Sprintf("info takes no argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdInfo()
	case "maxmem":
		if flag.NArg() != 1 {
			invalidUsageCommand(fmt.Sprintf("maxmem takes no argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdMax()
	case "minmem":
		if flag.NArg() != 1 {
			invalidUsageCommand(fmt.Sprintf("minmem takes no argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdMin()
	case "tagget":
		if flag.NArg() != 2 {
			invalidUsageCommand(fmt.Sprintf("tagget takes 1 argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdTagGet(flag.Arg(1))
	case "tagset":
		if flag.NArg() != 2 {
			invalidUsageCommand(fmt.Sprintf("tagset takes 1 argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdTagSet(flag.Arg(1))
	case "tags":
		if flag.NArg() != 1 {
			invalidUsageCommand(fmt.Sprintf("tags takes no argument but %d given!", flag.NArg()-1), command)
		}
		err = CmdTags()
	default:
		invalidUsage(fmt.Sprintf("Invalid command '%s'!", command))
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
