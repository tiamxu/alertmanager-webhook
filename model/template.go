package model

import (
	"bytes"
	"fmt"
	tmplhtml "html/template"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/tiamxu/alertmanager-webhook/log"
)

type Template struct {
	Name string
	tmpl *template.Template
}

func NewTemplate(filename string) (*Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	funcMap := template.FuncMap{
		"GetTimeDuration": GetTimeDuration,
		"GetCSTtime":      GetCSTtime,
		"TimeFormat":      TimeFormat,
		"GetTime":         GetTime,
		"toUpper":         strings.ToUpper,
		"toLower":         strings.ToLower,
		"title":           strings.Title,
		// join is equal to strings.Join but inverts the argument order
		// for easier pipelining in templates.
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"safeHtml": func(text string) tmplhtml.HTML {
			return tmplhtml.HTML(text)
		},
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"stringSlice": func(s ...string) []string {
			return s
		},
		"SplitString": func(pstring string, start int, stop int) string {
			log.Infof("SplitString", pstring)
			if stop < 0 {
				return pstring[start : len(pstring)+stop]
			}
			return pstring[start:stop]
		},
	}
	tmpl, err := template.New("").Funcs(funcMap).Parse(string(data))
	// tmpl, err := template.New("alert").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return &Template{tmpl: tmpl}, nil
}

func (t *Template) Execute(data interface{}) (string, error) {
	// log.Infof("Template data: %+v", data)
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 转换时间戳到时间字符串
func GetTime(timeStr interface{}, timeFormat ...string) string {
	var R_Time string
	//判断传入的timeStr是否为float64类型，如gerrit消息中时间戳就是float64
	switch timeStr.(type) {
	case string:
		S_Time, _ := strconv.ParseInt(timeStr.(string), 10, 64)
		if len(timeFormat) == 0 {
			timeFormat = append(timeFormat, "2006-01-02T15:04:05")
		}
		if len(timeStr.(string)) == 13 {
			R_Time = time.Unix(S_Time/1000, 0).Format(timeFormat[0])
		} else {
			R_Time = time.Unix(S_Time, 0).Format(timeFormat[0])
		}
	case float64:
		if len(timeFormat) == 0 {
			timeFormat = append(timeFormat, "2006-01-02T15:04:05")
		}
		R_Time = time.Unix(int64(timeStr.(float64)), 0).Format(timeFormat[0])
	}
	return R_Time
}

// 转换时间为持续时长
func GetTimeDuration(startTime string, endTime string) string {
	var tm = "N/A"
	if startTime != "" && endTime != "" {
		starT1 := startTime[0:10]
		starT2 := startTime[11:19]
		starT3 := starT1 + " " + starT2
		startm2, err := time.Parse("2006-01-02 15:04:05", starT3)
		if err != nil {
			return tm // 如果解析失败，则返回N/A
		}

		endT1 := endTime[0:10]
		endT2 := endTime[11:19]
		endT3 := endT1 + " " + endT2
		endm2, err := time.Parse("2006-01-02 15:04:05", endT3)
		if err != nil {
			return tm // 如果解析失败，则返回N/A
		}

		sub := endm2.UTC().Sub(startm2.UTC())

		t := int64(sub.Seconds())
		if t >= 86400 {
			days := t / 86400
			hours := (t % 86400) / 3600
			tm = fmt.Sprintf("%dd%dh", days, hours)
		} else {
			hours := t / 3600
			minutes := (t % 3600) / 60
			if hours > 0 {
				tm = fmt.Sprintf("%dh%dm", hours, minutes)
			} else {
				// 如果小时为0，则只显示分钟和秒
				seconds := t % 60
				tm = fmt.Sprintf("%dm%ds", minutes, seconds)
				if minutes == 0 {
					// 如果分钟也为0，则只显示秒
					tm = fmt.Sprintf("%ds", seconds)
				}
			}
		}
	}
	return tm
}

// 转换UTC时区到CST
func GetCSTtime(date string) string {
	var tm string
	tm = time.Now().Format("2006-01-02 15:04:05")
	if date != "" {
		T1 := date[0:10]
		T2 := date[11:19]
		T3 := T1 + " " + T2
		tm2, _ := time.Parse("2006-01-02 15:04:05", T3)
		h, _ := time.ParseDuration("-1h")
		tm3 := tm2.Add(-8 * h)
		tm = tm3.Format("2006-01-02 15:04:05")
	}
	return tm
}

func TimeFormat(timestr, format string) string {
	returnTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", timestr)
	if err != nil {
		returnTime, err = time.Parse("2006-01-02T15:04:05.999999999+08:00", timestr)
	}
	if err != nil {
		return err.Error()
	} else {
		return returnTime.Format(format)
	}
}
