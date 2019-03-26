# loggos

## Contents

  Fill me in

## Description

Loggos is a flexible logging package for sending logs in line format and json.
Born out of need for a logging package that does both standard line logging and also json logging.

## How To Use Loggos

Loggos is intended to be easy to use and flexible enough to fit most projects.

There are shortcut methods to make use of both line logger and the JSON Logger. They mimic the signature of the fmt package. These are used for simple usage of the loggers. More advanced logging is discussed further down.

Line messages are simple message that hold generally unstructured logs and are very useful in small project.

```go
// Sending line messages
Infoln("Test Message - Infoln")
Warnln("Test Message - Warnln")
Critln("Test Message - Critln")
Debugln("Test Message - Debugln")

Infof("Test Message - %s", "Infof")
Warnf("Test Message - %s", "Warnf")
Critf("Test Message - %s", "Critf")
Debugf("Test Message - %s", "Debugf")
```

JSON logging is used for structured logging and generally requires a bit more effort to use but is far
better in the long run as the data in the logs is easy to digest into centralized logging platforms and
aggregators.

```go
// Make a message using one of these lines.
m := JSONInfoln("Test Message - JSONInfoln")
m := JSONWarnln("Test Message - JSONWarnln")
m := JSONCritln("Test Message - JSONCritln")
m := JSONDebugln("Test Message - JSONDebugln")

m := JSONInfof("Test Message - %s", "JSONInfof")
m := JSONWarnf("Test Message - %s", "JSONWarnf")
m := JSONCritf("Test Message - %s", "JSONCritf")
m := JSONDebugf("Test Message - %s", "JSONDebugf")

//Decorate with the Add function
m.Add(key string, value interface{})

// Then Send
SendJSON(m)
```

Debugging needs to be turned on using the following functions.

```go
DefaultJSONLogger.EnableDebugLogging(true|false)
DefaultLineLoggerBuffer.EnableDebugLogging(true|false)
```

## Flexibility of Loggos

Loggos flexibility comes in 3 features.

* Allow users to override where logs get shipped to
* Allow users to override the time stamping
* Allow users to choose between best effort and audit shipping

### Overriding the Send function

Sending logs to a aggregator is something that all people that need to do it do differently.
It's dependent on your project, company and products that you make use of. It is therefore best left to the operator to make that choice and integration.

By default logs will get printed out to STDOUT. Which in most cases is fine. However when it is not you can override the where the logs go and maybe process them before shipping. See below an example of a override.
```go
// Overrides have a interface to satisfy. Very simple.
//type Overrider interface {
//	Send(string)
//}

// Therefore this is acceptable.
type overrider struct{}
func (o overrider) Send(){}

overRider := overrider{}

DefaultJSONLogger.OverridePrinter(overRide)
DefaultLineLogger.OverridePrinter(overRide)
```

### Overriding the time stamp

#### Line logger

When using the line logger you need to set and overrider for timestamps. The only requirement is that your passed in fuction need return a single string. It is not the business of loggos to decide what your time stamps should look like. In fact in testing we just set them to a static word. You are the master of your own time stamp.

```go
// Change the default loggers time stamping
func timeStampOverrider() string {
  return "overridden"
}
DefaultLineLogger.OverrideTimeStamping(timeStampOverrider)
```

#### JSON Messages

JSON messages offer greater flexibility by storing the time stamp as a message field that you can override after creating the message.

The time stamp uses the default time stamp gained from `time.Now().UnixNano()`.

If you to change the default timestamp you need to change the value stored in `loggos.jsonessage.JSONTimeStampFunc` with your own struct that satisfies the interface of `loggos.jsonessage.JSONTimeStamper`. Which is basically anything with a `Stamp() string` function.

The time stamp message field can also be changed by changing the string stored in `loggos.jsonessage.JSONTimeStampKey`.

```go
// Make something that satisfies loggos.jsonessage.JSONTimeStamper
type stamper struct {}
func (s stamper) Stamp() string {return "moreFunStuff"}

// Set the default vars
loggos.jsonessage.JSONTimeStampKey = "somethingFun"
loggos.jsonessage.JSONTimeStamper = &stamper{}
```

### Changing the loggers behavior

Logging generally follows two paths with regards to shipping.

1. Logging is less important than service and should never impact service delivery. This is called best effort logging.
1. Logging is more important than service delivery and must be done before moving on. This is called auditing.

With this in mind there are 2 ways that loggos will work. Unsurprisingly:

#### Best Effort

Logs are put into a buffer and printed as fast as the logger can push them out. The buffer is set by default in `loggos.DefaultLineLoggerBuffer` and `loggos.DefaultJSONLoggerBuffer`. If the buffer is filled up then new messages are dropped. If your logger is slow then you should increase the buffer limits. Just beware that larger buffer means more potential memory usage.

#### Audit Mode

Logs are put into a buffer, if the buffer is full then the logger waits till there is space available in the buffer. You can tweak this to allow for some risk by increasing the buffer size. If you are planning on using audit mode then you should make a really small buffer size, maybe 1.
