package secure

import (
	"bytes"
	"os"
	"time"
	"tp-system/config"
	"tp-system/model"

	log "github.com/sirupsen/logrus"
)

// InitLogRotate - Log Rotation initiated
func InitLogRotate() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	d := n.Sub(t)
	checkCurrentFile()
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour
		currentTime := time.Now()

		var LP bytes.Buffer
		LP.WriteString(config.Env.LogPath)
		LP.WriteString(currentTime.Format("2006-01-02"))
		LP.WriteString("_debug.log")
		logFile := LP.String()

		var file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err.Error())
		}
		log.SetOutput(file)
		log.SetFormatter(&log.TextFormatter{})
		// file.Close()
	}
}

func checkCurrentFile() {
	var LP bytes.Buffer
	currentTime := time.Now()
	LP.WriteString(config.Env.LogPath)
	LP.WriteString(currentTime.Format("2006-01-02"))
	LP.WriteString("_debug.log")
	logFile := LP.String()
	var file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println(err.Error())
	}
	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{})
}

// CronRoutine - It executes functions on routine
func CronRoutine() {
	Every15Day := 1
	Todaysdate := model.GetLocalDateSystem()
	model.FillInvRateRestData()
	//util.CacheError()
	for {
		time.Sleep(5 * time.Minute)
		Newdate := model.GetLocalDateSystem()
		if Newdate != Todaysdate {
			model.FillInvRateRestData()
			Todaysdate = Newdate
		}
		if Every15Day == 4320 {
			//util.CacheError()
			Every15Day = 0
		}
		// model.ErrorJob()
	}
}
