package log

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// LogFormatter represents a log formatter
type LogFormatter struct {
	Prefix       string
	DisableLevel bool
	ForceColors   bool // 添加颜色支持标志
}

// Format writes to log according to the log entry
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 添加颜色代码（如果启用）
	if f.ForceColors {
		b.WriteString(f.getLogLevelColor(entry.Level))
	}

	b.WriteString(utils.FormatUnixTimeToLongDateTimeInServerTimezone(time.Now().Unix()))
	b.WriteString(" ")

	if f.Prefix != "" {
		b.WriteString(f.Prefix)
		b.WriteString(" ")
	}

	if !f.DisableLevel {
		b.WriteString("[")
		b.WriteString(strings.ToUpper(entry.Level.String()))
		b.WriteString("] ")
	}

	if requestId, exists := entry.Data[logFieldRequestId]; exists && requestId != "" {
		b.WriteString(fmt.Sprintf("[%s] ", requestId))
	}

	b.WriteString(entry.Message)

	// 重置颜色
	if f.ForceColors {
		b.WriteString("\x1b[0m")
	}

	b.WriteString("\n")

	if extra, exists := entry.Data[logFieldExtra]; exists {
		b.WriteString(extra.(string))
	}

	return b.Bytes(), nil
}

// getLogLevelColor 返回不同日志级别的颜色代码
func (f *LogFormatter) getLogLevelColor(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "\x1b[90m" // 深灰色 (TRACE)
	case logrus.DebugLevel:
		return "\x1b[36m" // 青色 (DEBUG)
	case logrus.InfoLevel:
		return "\x1b[32m" // 绿色 (INFO)
	case logrus.WarnLevel:
		return "\x1b[33m" // 黄色 (WARN)
	case logrus.ErrorLevel:
		return "\x1b[31m" // 红色 (ERROR)
	case logrus.FatalLevel:
		return "\x1b[35m" // 紫色 (FATAL)
	case logrus.PanicLevel:
		return "\x1b[35;1m" // 加粗紫色 (PANIC)
	default:
		return ""
	}
}