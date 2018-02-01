package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/qacafe/go-usp-protobuf/usp"
	"github.com/qacafe/go-usp-protobuf/usp_record"
)

//go:generate protoc --go_out=. usp/usp-msg.proto
//go:generate protoc --go_out=. usp_record/usp-record.proto
func main() {
	encodeMsg := flag.Bool("encode-msg", false, "Encode USP Msg JSON document from stdin")
	decodeMsg := flag.Bool("decode-msg", false, "Decode USP Msg hex string from stdin")
	encodeRecord := flag.Bool("encode-record", false, "Encode USP Record JSON document from stdin")
	decodeRecord := flag.Bool("decode-record", false, "Decode USP Record hex string from stdin")
	indent := flag.Bool("indent", false, "Indent JSON documents printed to stdout")
	allowUnknown := flag.Bool("allow-unknown", false, "Allow unknown fields when decoding")
	emitDefaults := flag.Bool("emit-defaults", false, "Emit default values when encoding")
	enumsAsInts := flag.Bool("enums-as-ints", false, "Emit enums as ints when encoding")

	flag.Parse()

	if !*encodeMsg && !*encodeRecord &&
		!*decodeMsg && !*decodeRecord {
		flag.Usage()
		os.Exit(1)
	}

	if *encodeMsg || *encodeRecord {
		dec := json.NewDecoder(os.Stdin)
		u := jsonpb.Unmarshaler{
			AllowUnknownFields: *allowUnknown,
		}

		var pb proto.Message

		if *encodeMsg {
			pb = &usp.Msg{}
		} else {
			pb = &usp_record.Record{}
		}

		err := u.UnmarshalNext(dec, pb)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		buf, err := proto.Marshal(pb)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}

		str := hex.EncodeToString(buf)
		fmt.Fprintln(os.Stdout, str)
		fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString([]byte(str)))
	} else if *decodeMsg || *decodeRecord {
		// read USP Record hex string from stdin
		r := bufio.NewReader(os.Stdin)
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatal("read error: ", err)
		}
		line = strings.TrimSuffix(line, "\n")
		buf, err := hex.DecodeString(line)
		if err != nil {
			log.Fatal("decode error: ", err)
		}

		var (
			pb     proto.Message
			msg    *usp.Msg
			record *usp_record.Record
		)

		if *decodeMsg {
			msg = &usp.Msg{}
			pb = msg
		} else {
			record = &usp_record.Record{}
			pb = record
		}

		// decode USP Record protobuf wire format
		err = proto.Unmarshal(buf, pb)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}
		// print USP Record JSON document to stdout
		m := jsonpb.Marshaler{
			EnumsAsInts:  *enumsAsInts,
			EmitDefaults: *emitDefaults,
		}
		if *indent {
			m.Indent = "  "
		}
		err = m.Marshal(os.Stdout, pb)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		fmt.Fprintf(os.Stdout, "\n")
		if *decodeMsg {
			fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString([]byte(line)))
		} else if *decodeRecord {
			if nsc, ok := record.RecordType.(*usp_record.Record_NoSessionContext); ok {
				fmt.Fprintln(os.Stdout, hex.EncodeToString(nsc.NoSessionContext.GetPayload()))
			}
			if sc, ok := record.RecordType.(*usp_record.Record_SessionContext); ok {
				for _, payload := range sc.SessionContext.GetPayload() {
					fmt.Fprintln(os.Stdout, hex.EncodeToString(payload))
				}
			}
		}
	}
}
