package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pg_backup/controllers"
	"runtime"
	"time"
)

var (
	cfg     controllers.Config
	dirPath string
	verbose bool
)

func helpMenu() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Program with the facility to facilitate the backup of the Postgresql database."+
			"\n\nHelp:\n\t-h\tDisplay this help\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v\t%v\n", f.Name, f.Usage) // f.Name, f.Value
		})

		fmt.Println()
		os.Exit(0)
	}

	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()
}

func init() {
	system := runtime.GOOS

	if system == "linux" {
		if !controllers.FileExists("/usr/bin/pg_dump") {
			log.Fatalln("command not found /usr/bin/pg_dump, needed to run")
		}
		if !controllers.FileExists("/usr/bin/psql") {
			log.Fatalln("command not found /usr/bin/psql, needed to run")
		}
	}

	helpMenu()

	controllers.ReadFile(&cfg)

	year, month, day := time.Now().Date()
	dirPath = fmt.Sprintf("%s/%d/%s/%d", cfg.Backup.Directory, year, month.String(), day)
	controllers.CreateBackupDirectory(dirPath)

}

func main() {

	var tables []string
	backupFormats := map[string]string{
		"plain":     "sql",
		"custom":    "backup",
		"tar":       "tar",
		"directory": "directory",
	}

	if cfg.Backup.Databases[0] == "all" {
		tables = controllers.ListTables(cfg.Server.Host, cfg.Server.Port, cfg.Database.Username, cfg.Database.Password)
	} else {
		tables = cfg.Backup.Databases
	}

	if verbose {
		fmt.Printf("Total databases found:\t%v\n\n", len(tables))
	}

	for _, table := range tables {
		var path = fmt.Sprintf("%s/%s.%s", dirPath, table, backupFormats[cfg.Backup.BackupFormat])
		controllers.ExecBackup(cfg.Server.Host, cfg.Server.Port, cfg.Database.Username, cfg.Database.Password, table, cfg.Backup.Encode, cfg.Backup.BackupFormat, path)
		if verbose {
			fmt.Printf("Database:\t%s\nSaved in:\t%s\n\n", table, path)
		}
	}
}
