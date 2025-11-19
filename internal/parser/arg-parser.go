package parser

import (
	"fmt"
	"log/slog"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/help"
)

func ParseArgs(args []string) {
	argLength := len(args)
	var program string
	var argument string

	if argLength == 0 {
		program = ""
	} else {
		program = args[0]
	}

	if argLength >= 2 {
		argument = args[1]
	} else {
		argument = ""
	}

	switch program {
	case "bug":
		openBugReport()

	case "version":
		printVersion()

	case "help":
		printHelp(argument)

	case "query":
		ParseFlags(args)
		if !VerifyPGNInput(argument) {
			os.Exit(1)
		}

	case "merge":
		ParseFlags(args)

	case "convert":
		ParseFlags(args)

	default:
		if program != "" {
			printHelp(program)
		} else {
			printHelp("")
		}
	}
}

func ParseFlags(args []string) {
	for i, arg := range args {
		if strings.EqualFold(arg, "--verbose") {
			global.ProgramLevel.Set(slog.LevelDebug)
		}
		if strings.EqualFold(arg, "--output") {
			if VerifyPGNOutput(args[i+1]) {
				global.Output = args[i+1]
			}
		}
		if strings.EqualFold(arg, "--experimental") {
			global.AllowExperimental = true
		}
	}
}

func printVersion() {
	fmt.Printf("pgn-tools version %s %s/%s\n", global.VERSION, runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}

func printHelp(command string) {
	switch command {
	case "":
		fmt.Println(help.Default)
	case "bug":
		fmt.Println(help.Bug)
	case "convert":
		fmt.Println(help.Convert)
	case "merge":
		fmt.Println(help.Merge)
	case "query":
		fmt.Println(help.Query)
	case "version":
		fmt.Println(help.Version)
	default:
		global.Logger.Info(fmt.Sprintf("unknown command: %v", command))
		fmt.Println(help.Default)
	}

	os.Exit(0)
}

func VerifyPGNInput(file string) bool {
	if file == "" {
		global.Logger.Error("Enter input filepath")
		return false
	}

	global.Logger.Debug(fmt.Sprintf("input file: %s", file))

	ext := filepath.Ext(file)
	filetype := mime.TypeByExtension(ext)

	global.Logger.Debug(fmt.Sprintf("Mime input: %s", filetype))

	if !strings.Contains(filetype, "chess") {
		global.Logger.Error(fmt.Sprintf("Invalid Filetype: %s", file))
		return false
	}

	_, err := os.Stat(file)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Invalid Filepath: %s", file))
		return false
	}
	return true
}

func VerifyPGNOutput(file string) bool {
	if file == "" {
		global.Logger.Error("Enter output filepath")
		return false
	}

	global.Logger.Debug(fmt.Sprintf("output file: %s", file))

	ext := filepath.Ext(file)
	filetype := mime.TypeByExtension(ext)

	global.Logger.Debug(fmt.Sprintf("Mime output: %s", filetype))

	if !strings.Contains(filetype, "chess") {
		global.Logger.Error(fmt.Sprintf("Invalid Filetype: %s", file))
		return false
	}

	return true
}

// include version, etc
func openBugReport() {
	var cmd string
	var args []string

	url := "https://github.com/gavink97/pgn-tools/issues/new?labels=bug&title=Title&body=Please+check+if+an+issue+containing+this+error+exists+before+submitting.+Also+try+to+provide+any+steps+we+can+use+to+reproduce+the+error."

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}

	if len(args) > 1 {
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("An unexpected error occured: %v", err))
		os.Exit(1)
	}

	os.Exit(0)
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
