package log

// Segment represents a portion of a log chain
type Segment struct {
	logger Logger
	event  string
	fields Fields
}

// Fields adds fields to the segment
func (segment *Segment) Fields(fields ...Fields) *Segment {

	segment.fields = segment.fields.Merge(fields...)
	return segment
}

// Debug logs a debug message for this segment
func (segment *Segment) Debug(message string) string {

	segment.fields[Message] = message

	return segment.logger.Debug(
		segment.event,
		segment.fields,
		)
}

// Info logs an info message for this segment
func (segment *Segment) Info(message string) string {

	segment.fields[Message] = message

	return segment.logger.Info(
		segment.event,
		segment.fields,
	)
}

// Warn logs a warn message for this segment
func (segment *Segment) Warn(message string) string {

	segment.fields[Message] = message

	return segment.logger.Warn(
		segment.event,
		segment.fields,
	)
}

// Error logs a error message for this segment
func (segment *Segment) Error(err error) string {

	return segment.logger.Error(
		segment.event,
		err,
		segment.fields,
	)
}

// Fatal logs a fatal message for this segment
func (segment *Segment) Fatal(err error)  {

	segment.logger.Fatal(
		segment.event,
		err,
		segment.fields,
	)
}

// Audit logs an audit message for this segment
func (segment *Segment) Audit(message string) string {

	segment.fields[Message] = message

	return segment.logger.Audit(
		segment.event,
		segment.fields,
	)
}

