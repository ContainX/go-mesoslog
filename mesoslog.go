package main

import (
	"fmt"
	ml "github.com/ContainX/go-mesoslog/mesoslog"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	// StdErrFlag is a flag to output stderr logs vs stdout if true
	StdErrFlag string = "stderr"
	// CompletedFlag is a flag to output completed tasks if true
	CompletedFlag string = "completed"
	// LatestFlag is a flag to only capture the latest instance.
	// This applies to completed and non-completed tasks
	LatestFlag string = "latest"
	// MasterFlag is the mesos master host:port flag
	MasterFlag string = "master"
	// DurationFlag is how often to poll in seconds
	DurationFlag string = "duration"
	// EnvMesosMaster is the MESOS_MASTER env variable
	EnvMesosMaster string = "MESOS_MASTER"
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

var tailCmd = &cobra.Command{
	Use:   "tail [appID]",
	Short: "Tails logs [appId] to StdOut.  Each running instance/task log will be outputed",
	Run:   tailLog,
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

func main() {
	rootCmd.PersistentFlags().Bool(StdErrFlag, false, "Output stderr log instead of default stdout")
	printCmd.PersistentFlags().Bool(CompletedFlag, false, "Use completed tasks (default: running tasks)")
	printCmd.PersistentFlags().Bool(LatestFlag, false, "Use the latest instance only")
	rootCmd.PersistentFlags().StringP(MasterFlag, "m", "", "Mesos Master host:port (eg. 192.168.2.1:5050) or ENV [MESOS_MASTER]")
	tailCmd.Flags().IntP(DurationFlag, "d", 5, "Log poll time (duration) in seconds")
	rootCmd.AddCommand(appsCmd, printCmd, tailCmd, fileCmd)
	rootCmd.Execute()
}

func printLog(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fmt.Println("ERROR: An [appId] must be specified")
		return
	}

	completed, _ := cmd.Flags().GetBool(CompletedFlag)
	latest, _ := cmd.Flags().GetBool(LatestFlag)

	logs, err := client(&ml.MesosClientOptions{SearchCompletedTasks: completed, ShowLatestOnly: latest}).GetLog(args[0], getLogType(), "")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	for _, log := range logs {
		fmt.Printf("\n=========================[ %s - Log For Task: %s ]============================\n", log.AppID, log.TaskID)
		fmt.Printf("%s\n", log.Log)
	}
}

func tailLog(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fmt.Println("ERROR: An [appId] must be specified")
		return
	}

	duration, _ := cmd.Flags().GetInt(DurationFlag)
	err := client(nil).TailLog(args[0], getLogType(), duration)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
}

func fileLog(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fmt.Println("ERROR: An [appId] and [output_dir] must be specified")
		return
	}

	completed, _ := cmd.Flags().GetBool(CompletedFlag)
	logs, err := client(&ml.MesosClientOptions{SearchCompletedTasks: completed}).GetLog(args[0], getLogType(), args[1])
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	for _, log := range logs {
		fmt.Printf("Log written to %s\n", log.Log)
	}
}

func listApps(cmd *cobra.Command, args []string) {
	apps, err := client(nil).GetAppNames()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "\nAPP_ID\tINSTANCES\n")
	for k, v := range apps {
		fmt.Fprintf(w, "%s\t%v\n", k, v)
	}
	flushWriter(w)
}

func getLogType() ml.LogType {
	if rootCmd.PersistentFlags().Changed(StdErrFlag) {
		if b, err := rootCmd.PersistentFlags().GetBool(StdErrFlag); err == nil && b {
			return ml.STDERR
		}
	}
	return ml.STDOUT
}

func client(options *ml.MesosClientOptions) *ml.MesosClient {
	var host string
	var port = 5050
	master, err := rootCmd.PersistentFlags().GetString(MasterFlag)

	if master == "" {
		if os.Getenv(EnvMesosMaster) == "" {
			printErr(fmt.Errorf("Must define a Master host and optional port"))
			os.Exit(1)
		}
		master = os.Getenv(EnvMesosMaster)
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

	c, err := ml.NewMesosClientWithOptions(host, port, options)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	return c
}

func printErr(err error) {
	fmt.Printf("\nError: %s\n", err.Error())
}

func flushWriter(w *tabwriter.Writer) {
	fmt.Fprintln(w, "")
	w.Flush()
}

func newTabWriter(output io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(output, 0, 8, 2, '\t', 0)
	return w
}
