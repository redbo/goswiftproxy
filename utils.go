package main

import (
	"errors"
	"fmt"
	"log/syslog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/vaughan0/go-ini"
)

type httpRange struct {
	start, end int64
}

func Urlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

var GMT *time.Location

func ParseDate(date string) (time.Time, error) {
	if GMT == nil {
		GMT, _ = time.LoadLocation("GMT")
	}
	if ius, err := time.ParseInLocation(time.RFC1123, date, GMT); err == nil {
		return ius, nil
	}
	if ius, err := time.ParseInLocation(time.RFC1123Z, date, GMT); err == nil {
		return ius, nil
	}
	if ius, err := time.ParseInLocation(time.ANSIC, date, GMT); err == nil {
		return ius, nil
	}
	if ius, err := time.ParseInLocation(time.RFC850, date, GMT); err == nil {
		return ius, nil
	}
	if timestamp, err := strconv.ParseFloat(date, 64); err == nil {
		nans := int64((timestamp - float64(int64(timestamp))) * 1.0e9)
		return time.Unix(int64(timestamp), nans).In(GMT), nil
	}
	return time.Now(), errors.New("invalid time")
}

func LooksTrue(check string) bool {
	check = strings.TrimSpace(strings.ToLower(check))
	return check == "true" || check == "yes" || check == "1" || check == "on" || check == "t" || check == "y"
}

type IniFile struct{ ini.File }

func (f IniFile) GetDefault(section string, key string, dfl string) string {
	if value, ok := f.Get(section, key); ok {
		return value
	}
	return dfl
}

func LoadIniFile(filename string) (IniFile, error) {
	file := IniFile{make(ini.File)}
	return file, file.LoadFile(filename)
}

func SetupLogger(facility string, prefix string) *syslog.Writer {
	facility_mapping := map[string]syslog.Priority{"LOG_USER": syslog.LOG_USER,
		"LOG_MAIL": syslog.LOG_MAIL, "LOG_DAEMON": syslog.LOG_DAEMON,
		"LOG_AUTH": syslog.LOG_AUTH, "LOG_SYSLOG": syslog.LOG_SYSLOG,
		"LOG_LPR": syslog.LOG_LPR, "LOG_NEWS": syslog.LOG_NEWS,
		"LOG_UUCP": syslog.LOG_UUCP, "LOG_CRON": syslog.LOG_CRON,
		"LOG_AUTHPRIV": syslog.LOG_AUTHPRIV, "LOG_FTP": syslog.LOG_FTP,
		"LOG_LOCAL0": syslog.LOG_LOCAL0, "LOG_LOCAL1": syslog.LOG_LOCAL1,
		"LOG_LOCAL2": syslog.LOG_LOCAL2, "LOG_LOCAL3": syslog.LOG_LOCAL3,
		"LOG_LOCAL4": syslog.LOG_LOCAL4, "LOG_LOCAL5": syslog.LOG_LOCAL5,
		"LOG_LOCAL6": syslog.LOG_LOCAL6, "LOG_LOCAL7": syslog.LOG_LOCAL7}
	logger, err := syslog.Dial("udp", "127.0.0.1:514", facility_mapping[facility], prefix)
	if err != nil || logger == nil {
		panic(fmt.Sprintf("Unable to dial logger: %s", err))
	}
	return logger
}

func GetTimestamp() string {
	return fmt.Sprintf("%016.05f", float64(time.Now().UnixNano())/1000000000.0)
}

func GetTransactionId() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
