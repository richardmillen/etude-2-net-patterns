package logger

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub"
	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// Severity as per syslog
// https://en.wikipedia.org/wiki/Syslog
type Severity byte

const (
	// Alert is a condition that should be corrected immediately.
	Alert Severity = 1 + iota
	// Critical conditions, such as hard device errors.
	Critical
	// Error conditions.
	Error
	// Warning conditions.
	Warning
	// Notice conditions are normal but significant conditions.
	Notice
	// Info messages.
	Info
	// Debug information normally used when debugging a program.
	Debug
)

const (
	// AlertTopic ...
	AlertTopic = "1"
	// CritTopic ...
	CritTopic = "21"
	// ErrorTopic ...
	ErrorTopic = "321"
	// WarnTopic ...
	WarnTopic = "4321"
	// NoticeTopic ...
	NoticeTopic = "54321"
	// InfoTopic ...
	InfoTopic = "654321"
	// DebugTopic ...
	DebugTopic = "7654321"
)

var fieldSep = []byte{','}

// New constructs a new Logger.
//
// HACK: temporarily passing in Publisher. this should be wrapped in the api somehow.
// TODO: look into opentracing (http://opentracing.io/), zipkin etc.
func New(pub *pubsub.Publisher) *Logger {
	l := &Logger{
		event: uuid.New(),
		pub:   pub,
	}
	return l
}

// Start constructs a new Logger.
//
// HACK: temporarily passing in Publisher. this should be wrapped in the api somehow.
// TODO: look into opentracing (http://opentracing.io/), zipkin etc.
func Start(pub *pubsub.Publisher, name string) *Logger {
	l := &Logger{
		name:  name,
		event: uuid.New(),
		pub:   pub,
	}
	l.start()
	return l
}

// StartChild constructs a new Logger as a child of an existing Event.
func StartChild(parent *Logger, name string) *Logger {
	l := &Logger{
		name:  name,
		event: uuid.New(),
		pub:   parent.pub,
	}
	l.start()
	return l
}

// Logger is used to log system events.
//
// sample rate?
// collecting strategy e.g. store & forward, immediate send, opportunistic, or some refined combination.
type Logger struct {
	name    string
	event   uuid.Bytes
	pub     *pubsub.Publisher
	started bool
}

// start is called to start logging.
func (l *Logger) start() {
	l.pub.Publish(InfoTopic, l.newMessage(l.event, []byte(fmt.Sprintf("start,<hostname>,%s", l.name))))
	l.started = true
}

// Close is called to stop logging.
func (l *Logger) Close() (err error) {
	if !l.started {
		return
	}

	l.pub.Publish(InfoTopic, l.newMessage(l.event, []byte("stop")))
	l.started = false
	return
}

// Printf ...
func (l *Logger) Printf(severity Severity, format string, a ...interface{}) {
	l.pub.Publish(l.getSeverityTopic(severity), l.newMessage(l.event, []byte(fmt.Sprintf(format, a...))))
}

// Print ...
func (l *Logger) Print(severity Severity, a ...interface{}) {
	l.pub.Publish(l.getSeverityTopic(severity), l.newMessage(l.event, []byte(fmt.Sprint(a...))))
}

// newMessage is called to generate a new log message.
//
// message format (each field is separated):
// + message (correlation) id
// + event id
// + timestamp
// + data (specific to type of event; inc. node address/name)
func (l *Logger) newMessage(event uuid.Bytes, data []byte) []byte {
	return utils.JoinBytes(uuid.New(), fieldSep, event, []byte("<timestamp>"), fieldSep, data)
}

// getSeverityTopic turns a severity value into a valid topic string.
func (l *Logger) getSeverityTopic(severity Severity) string {
	v := int(severity)
	var buf bytes.Buffer

	for v > 0 {
		buf.WriteString(strconv.Itoa(v))
		v--
	}

	return buf.String()
}
