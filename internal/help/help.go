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

	verbose		print debug messages

Use "pgn-tools help <command>" for more information about a command.`
	Bug = `Usage: pgn-tools bug

Bug opens the default browser and starts a new bug report.`
	Convert = `Convert has not been implemented yet, sorry :'(`
	Merge   = `Merge has not been implemented yet, sorry :'(`
	Query   = `Usage: pgn-tools query PATH "key=value"

Query takes a pgn database path and the query(ies) which is a string array of
key value pairs used to match the game.

There are two types of queries, string queries and integer queries.

String queries look for a partial string in the game metadata based on the key,
for example if you wanted to look for Fide events you would use "event=fide".

Integer queries on the other hand compare the value from the query to the key
to match results. For example, if you were looking for games above 2500 elo you
could use "elo>=2500".

Example queries:
"elo>=2300"
"player=Carlsen"
"site!=chess.com"`
	Version = `Usage: pgn-tools version

Version prints the binaries version details.`
)
