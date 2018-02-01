go-usp-protobuf
===============

A testing tool for encoding/decoding USP protobuf messages

Install
=======

``` shell
go get -u github.com/qacafe/go-usp-protobuf
```

Usage
=====

``` shell
$ go-usp-protobuf -h
Usage of go-usp-protobuf:
  -allow-unknown
    	Allow unknown fields when decoding
  -decode-msg
    	Decode USP Msg hex string from stdin
  -decode-record
    	Decode USP Record hex string from stdin
  -emit-defaults
    	Emit default values when encoding
  -encode-msg
    	Encode USP Msg JSON document from stdin
  -encode-record
    	Encode USP Record JSON document from stdin
  -enums-as-ints
    	Emit enums as ints when encoding
  -indent
    	Indent JSON documents printed to stdout
```

Either `-decode-msg`, `-decode-record`, `-encode-msg` or
`-encode-record` must be provided.  The message/record to
encode/decode is read from stdin.

Examples
========

Encoding a USP message:

``` shell
$ echo '{"header": {"msg_id": "deadadsad", "msg_type": "GET"}}' | go-usp-protobuf -encode-msg
0a0d0a096465616461647361641001
MGEwZDBhMDk2NDY1NjE2NDYxNjQ3MzYxNjQxMDAx
```

For convenience, the encoded USP message is printed first as a hex
string and then as a base64 string.

Encoding a USP record:

``` shell
$ echo '{"to_id": "deadbeef", "from_id": "feebdaed", "no_session_context": {"payload": "Cg0KCWRlYWRhZHNhZBAB"}}' | go-usp-protobuf -encode-record
120864656164626565661a0866656562646165643a11120f0a0d0a096465616461647361641001
MTIwODY0NjU2MTY0NjI2NTY1NjYxYTA4NjY2NTY1NjI2NDYxNjU2NDNhMTExMjBmMGEwZDBhMDk2NDY1NjE2NDYxNjQ3MzYxNjQxMDAx
```

For convenience, the encoded USP record is printed first as a hex
string and then as a base64 string.

Decoding a USP message:

``` shell
$ echo 0a0d0a096465616461647361641001 | go-usp-protobuf -decode-msg
{"header":{"msgId":"deadadsad","msgType":"GET"}}
MGEwZDBhMDk2NDY1NjE2NDYxNjQ3MzYxNjQxMDAx
```

For convenience, the decoded USP message is printed as a base64 string
after the JSON representation.

Decoding a USP record:

``` shell
$ echo 120864656164626565661a0866656562646165643a11120f0a0d0a096465616461647361641001 | go-usp-protobuf -decode-record
{"toId":"deadbeef","fromId":"feebdaed","noSessionContext":{"payload":"Cg0KCWRlYWRhZHNhZBAB"}}
0a0d0a096465616461647361641001
```

For convenience, the decoded USP record's payload is printed as a hex
string after the JSON representation.
