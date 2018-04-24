package main

import (
	"flag"
	"fmt"
	"os"
)

// PLMUrl to PLM server (daemon)
var PLMUrl string

func printUsage() {
	fmt.Printf("usage: plmc [options] <command> [<args>]\n\n")
	fmt.Printf(" General options:\n")
	fmt.Printf("  -h   Help (overview)\n")
	fmt.Printf("  -v   Display version\n")
	fmt.Printf("  -u   PLM server (daemon) URL. Default http://localhost:12124\n")
	fmt.Printf("\n Commands:\n")
	fmt.Printf("  help   Help for a command\n")
	fmt.Printf("  plot   Download plot for one or more processes\n")
	fmt.Printf("  info   List info about one or more processes\n")
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
		fmt.Printf("  -m <string>   Plot all processes matching the string.\n")
	case "info":
		fmt.Printf("List process info.\n")
		fmt.Printf("By default all processes are listed. Can be resttricted\n")
		fmt.Printf("with options described below\n\n")
		fmt.Printf("Usage: plmc [options] info\n\n")
		fmt.Printf(" Options:\n")
		fmt.Printf("  -m <string>   List all processes matching the string.\n")
	default:
		fmt.Fprintf(os.Stderr, "No such command %s\n\n", command)
		printUsage()
	}
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
	flag.StringVar(&PLMUrl, "u", "http://localhost:12124", "PLM server URL")
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
	default:
		invalidUsage(fmt.Sprintf("Invalid command '%s'!", command))
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
