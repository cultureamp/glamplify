package newrelic

import "github.com/cultureamp/glamplify/log"

type Segment struct {
	logger Logger
	event  string
	fields log.Fields
}

func (segment *Segment) Fields(fields ...log.Fields) *Segment {

	segment.fields = segment.fields.Merge(fields...)
	return segment
}

func (segment *Segment) Debug(message string) {

	segment.fields[log.Message] = message

	segment.logger.Debug(
		segment.event,
		segment.fields,
	)
}

func (segment *Segment) Info(message string) {

	segment.fields[log.Message] = message

	segment.logger.Info(
		segment.event,
		segment.fields,
	)
}

func (segment *Segment) Warn(message string) {

	segment.fields[log.Message] = message

	segment.logger.Warn(
		segment.event,
		segment.fields,
	)
}

func (segment *Segment) Error(err error) {

	segment.logger.Error(
		segment.event,
		err,
		segment.fields,
	)
}

func (segment *Segment) Fatal(err error) {

	segment.logger.Fatal(
		segment.event,
		err,
		segment.fields,
	)
}



