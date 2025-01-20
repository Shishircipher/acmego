package logger

import (
	"log"
	"os"
	"time"
	"fmt"
)

var (
	logFile *os.File
	fileLogger  *log.Logger // For logging all messages to the log file
)

func init() {
	var err error

	// Create the log directory if it doesn't exist
	err = os.MkdirAll("./.acmego/log", 0700)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
		os.Exit(1)
	}

	// Generate a timestamp for the log file name
	timestamp := time.Now().Format("2006-01-02_15-04-05") // Example: 2025-01-18_13-45-00
//	logFileName := fmt.Sprintf("./acmego/log/log_%s.log", timestamp)
	logFileName := "./.acmego/log/" + timestamp + ".log"
	// Open the log file with the generated name
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
		os.Exit(1)
	}

	// Set the log output to the file
	log.SetOutput(logFile)
	fileLogger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Fatalf(format string, args ...interface{}) {
        fmt.Printf(format, args...)
        fileLogger.Fatalf(format, args...)

}

// Fatal logs a critical error message and exits the program
func Fatal(format string, v ...interface{}) {
	log.Printf("FATAL: "+format, v...)
	os.Exit(1)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

// Info logs an informational message
func Info(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	log.Printf("DEBUG: "+format, v...)
}
func Spinner(seconds int) {
	fmt.Print("Processing ... ")
	timeout := time.After(time.Duration(seconds) * time.Second) // Timer for the given duration
	for {
		select {
		case <-timeout: // Exit the spinner after the timeout
			fmt.Println("\nDone")
			return
		default:
			// Spinner animation
			for _, r := range `-\|/` {
				fmt.Printf("\rProcessing ... %c", r)
				time.Sleep(100 * time.Millisecond) // Adjust the delay as needed
			}
		}
	}
}
func Spinner1(delay time.Duration) {
	fmt.Println("Processing ...")
	for {
		for _, r := range `-\|/`{
			fmt.Printf("\r %c", r)
                        fmt.Printf("%c", r)
			time.Sleep(delay)
		}
	}
}
