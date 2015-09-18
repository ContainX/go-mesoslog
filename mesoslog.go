package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	StdErrFlag string = "stderr"
	MasterFlag string = "master"
)

var rootCmd = &cobra.Command{
	Use:   "mesoslog",
	Short: "Download container logs from a mesos cluster",
	Long: `
Mesos Log is a quick utility to quickly pull stdout or stderr container logs from a Mesos cluster.
	`,
}

var printCmd = &cobra.Command{
	Use:   "print [appID]",
	Short: "Outputs the log for the given [appId] to StdOut.  Each running instances log will be outputed",
	Run:   printLog,
}

var fileCmd = &cobra.Command{
	Use:   "file [appID] [output_dir]",
	Short: "Outputs the log for the given [appId] to a file. Multiple files will be created (1 per running instance)",
	Run:   fileLog,
}

var appsCmd = &cobra.Command{
	Use:   "list",
	Short: "List current application id's and task count (instances running)",
	Run:   listApps,
}

func printLog(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fmt.Println("ERROR: An [appId] must be specified")
		return
	}

	logs, err := client().GetLog(args[0], getLogType(), "")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	for _, log := range logs {
		fmt.Printf("\n=========================[ %s - Log For Task: %s ]============================\n", log.AppID, log.TaskID)
		fmt.Printf("%s\n", log.Log)
	}
}

func main() {
	rootCmd.PersistentFlags().Bool(StdErrFlag, false, "Output stderr log instead of default stdout")
	rootCmd.PersistentFlags().StringP(MasterFlag, "m", "", "Mesos Master host:port (eg. 192.168.2.1:5050)")
	rootCmd.AddCommand(appsCmd, printCmd, fileCmd)
	rootCmd.Execute()
}

func fileLog(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fmt.Println("ERROR: An [appId] and [output_dir] must be specified")
		return
	}
	logs, err := client().GetLog(args[0], getLogType(), args[1])
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	for _, log := range logs {
		fmt.Printf("Log written to %s\n", log.Log)
	}
}

func listApps(cmd *cobra.Command, args []string) {
	apps, err := client().GetAppNames()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	w := NewTabWriter(os.Stdout)
	fmt.Fprintf(w, "\nAPP_ID\tINSTANCES\n")
	for k, v := range apps {
		fmt.Fprintf(w, "%s\t%v\n", k, v)
	}
	FlushWriter(w)
}

func getLogType() LogType {
	if rootCmd.PersistentFlags().Changed(StdErrFlag) {
		if b, err := rootCmd.PersistentFlags().GetBool(StdErrFlag); err == nil && b {
			return STDERR
		}
	}
	return STDOUT
}

func client() *MesosClient {
	var host string
	var port int = 5050
	if master, err := rootCmd.PersistentFlags().GetString(MasterFlag); err != nil {
		printErr(err)
		os.Exit(1)
	} else {
		if master == "" {
			printErr(fmt.Errorf("Must define a Master host and optional port"))
			os.Exit(1)
		}
		if strings.Contains(master, ":") {
			hp := strings.Split(master, ":")
			host = hp[0]
			port, err = strconv.Atoi(hp[1])
			if err != nil {
				printErr(err)
				os.Exit(1)
			}
		} else {
			host = master
		}
	}

	c, err := NewMesosClient(host, port)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	return c
}

func printErr(err error) {
	fmt.Printf("\nError: %s\n", err.Error())
}

func FlushWriter(w *tabwriter.Writer) {
	fmt.Fprintln(w, "")
	w.Flush()
}

func NewTabWriter(output io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(output, 0, 8, 2, '\t', 0)
	return w
}
