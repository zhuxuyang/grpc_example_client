package resource

import (
	"fmt"
	"github.com/phachon/go-logger"
	"github.com/spf13/viper"
	"log"
)

var Logger *go_logger.Logger

func InitLogger() {
	log.Println("init logger start")
	Logger = go_logger.NewLogger()
	logFormat := "[%timestamp_format%] [%level_string%] %function%:%line% %body%"
	logDir := viper.GetString("logDir")
	fileName := fmt.Sprintf("%s/grpc-example-client.log", logDir)
	// file adapter config
	fileConfig := &go_logger.FileConfig{
		//Filename: "./reviewer-test.log", // The file name of the logger output, does not exist automatically
		// If you want to separate separate logs into files, configure LevelFileName parameters.
		LevelFileName: map[int]string{
			Logger.LoggerLevel("error"): fileName, // The error level log is written to the error.log file.
			Logger.LoggerLevel("info"):  fileName, // The info level log is written to the info.log file.
			Logger.LoggerLevel("debug"): fileName, // The debug level log is written to the debug.log file.
		},
		MaxSize:    0,         // File maximum (KB), default 0 is not limited
		MaxLine:    0,         // The maximum number of lines in the file, the default 0 is not limited
		DateSlice:  "y",       // Cut the document by date, support "Y" (year), "m" (month), "d" (day), "H" (hour), default "no".
		JsonFormat: false,     // Whether the file data is written to JSON formatting
		Format:     logFormat, // JsonFormat is false, logger message written to file format string
	}
	// add output to the file
	_ = Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	log.Println("log to dir : " + logDir)
}

