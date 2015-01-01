connxray, `net.Listener` on steroids [![GoDoc](https://godoc.org/github.com/marcinwyszynski/connxray?status.svg)](https://godoc.org/github.com/marcinwyszynski/connxray)
========

This library provides a wrapper for net.Listener and net.Conn interfaces, allowing dynamic, callback-based introspection of those network primives. The rationale behind this library is that these network primitive are often used by higher level libraries (eg. net/http) that hide a lot of information that can be useful for debugging, monitoring or otherwise dynamically messing with a connection and its traffic. Please see the `examples` directory for some inspiration on how this library can be used.
