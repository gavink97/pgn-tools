package help


const (
	Default = `pgn-tools are designed to make chess database management easier.

Usage:

        pgn-tools <command> [arguments] [--flags]

The commands are:

		bug			start a bug report
        convert     convert a chessbase cbh to pgn
        merge       reconcile multiple databases into one database
        query      	query a pgn database
		version		print pgn-tools version

Global flags available:

		verbose		print debug messages

Use "pgn-tools help <command>" for more information about a command.`

)
