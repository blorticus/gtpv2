# Overview

This package provides a module for the GPRS Tunnel Protocol (GTP) version 2.  GTPv2 is documented in [3GPP TS 29.274](https://portal.3gpp.org/desktopmodules/Specifications/SpecificationDetails.aspx?specificationId=1692).

# Installation

```bash
go get github.com/blorticus/gtpv2
```

# Basic Usage

```golang
package main

import (
    gtpv2 "github.com/blorticus/gtpv2"
    "net"
    "fmt"
)

func main() {
    modifyBearerRequest := gtpv2.NewPDU(gtpv2.ModifyBearerRequest, 0x00001acc, []*gtpv2.IE{
        gtpv2.NewIEWithRawData(gtpv2.UserLocationInformation, []byte{
                0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
                0x00, 0x0f, 0x42, 0x4d, 0x00,
        }),
        gtpv2.NewIEWithRawData(gtpv2.RATType, []byte{0x06}),
        gtpv2.NewIEWithRawData(gtpv2.DelayValue, []byte{0x00}),
        gtpv2.NewIEWithRawData(gtpv2.BearerContext, []byte{
                0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
                0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13,
                0x01, 0xb2,
        }),
        gtpv2.NewIEWithRawData(gtpv2.RecoveryRestartCounter, []byte{0x95}),
    })
    
    conn, err := net.Dial("udp", "10.1.10.10:2123")
    if err != nil {
        panic(err)
    }
    
    conn.Write(modifyBearerRequest.Encode())
    
    incomingByteStream := make([]byte, 65536)
    
    dgLength, err := conn.Read(incomingByteStream)
    if err != nil {
        panic(err)
    }
    
    incomingGtpPDU, piggybackedGtpPDU, err := gtpv2.DecodePDU(incomingByteStream[:dgLength])
    if err != nil {
        panic(err)
    }
    
    for _, ie := range incomingGtpPDU.InformationElements {
        fmt.Printf("IE name = (%s), value = (%02x)\n", gtpv2.NameOfIEForType(ie.Type), ie.Data)
    }
}
```

# Information Elements

There is no support for interpretation of Information Elements.  They are created, stored, and presented as
a byte stream.  So, for example, when one uses the IMSI IE, it is not decoded to its natural string representation.

```
