package help

const (
	Default = `pgn-tools make chess database management easier.

Usage:

    	pgn-tools <command> [arguments] [--flags]

The commands are:

	bug 		start a bug report
	convert		convert a chessbase cbh to pgn
	merge		reconcile multiple databases into one database
	query		query a pgn database
	version		print pgn-tools version

Global flags available:

	experimental enables experimental features
	output		defines the output path in commands that use an output
	verbose		print debug messages

Use "pgn-tools help <command>" for more information about a command.`
	Bug = `Usage: pgn-tools bug

Bug opens the default browser and starts a new bug report.`
	Convert = `Usage: pgn-tools convert INPUT_PATH OUTPUT_PATH [--flags]

Convert is an experimental feature and will most likely panic mid conversion.
Please raise a bug report should such an instance occur.

Convert takes a chessbase database which must include a header file (.cbh),
player file (.cbp), tournament file (.cbt), and game file (.cbg), all in the
same directory as the input path, and converts it to a pgn database.`
	Merge = `Usage: pgn-tools merge PATH... '-o | --output PATH'  [--flags]

Merge takes multiple pgn database paths or directories containing pgn databases
and merges them in the output.`
	Query = `Usage: pgn-tools query PATH "key=value" [--flags]

Query takes a pgn database path and the query(ies) which is a string array of
key value pairs used to match the game.

There are two types of queries, string queries and integer queries.

String queries look for a partial string in the game metadata based on the key,
for example if you wanted to look for Fide events you would use "event=fide".

Integer queries on the other hand compare the value from the query to the key
to match results. For example, if you were looking for games above 2500 elo you
could use "elo>=2500".

Flags available:
	output		writes to output path

Example queries:
"elo>=2300"
"player=Carlsen"
"site!=chess.com"`
	Version = `Usage: pgn-tools version

Version prints the binaries version details.`
)
