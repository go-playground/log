# log
Highly configurable, structured logging that is a drop in replacement for the std library log

#### Log Level Definitions

**DebugLevel** - Info useful to developers for debugging the application, not useful during operations.

**TraceLevel** - Info useful to developers for debugging the application and reporting on possible bottlenecks.

**InfoLevel** - Normal operational messages - may be harvested for reporting, measuring throughput, etc. - no action required.

**NoticeLevel** - Normal but significant condition. Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required.

**WarnLevel** - Warning messages, not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time.

**ErrorLevel** - Non-urgent failures, these should be relayed to developers or admins; each item must be resolved within a given time.

**PanicLevel** - A "panic" condition usually affecting multiple apps/servers/sites. At this level it would usually notify all tech staff on call.

**AlertLevel** - Action must be taken immediately. Should be corrected immediately, therefore notify staff who can fix the problem. An example would be the loss of a primary ISP connection.

**FatalLevel** - Should be corrected immediately, but indicates failure in a primary system, an example is a loss of a backup ISP connection. ( same as SYSLOG CRITICAL )
