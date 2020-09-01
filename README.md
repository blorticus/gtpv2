# Overview

This package provides a module for the GPRS Tunnel Protocol (GTP) version 2.  GTPv2 is documented in [3GPP TS 29.274](https://portal.3gpp.org/desktopmodules/Specifications/SpecificationDetails.aspx?specificationId=1692).

# Installation

```bash
go get github.com/blorticus/gtpv2
```

# Basic Usage

```golang
import (
    gtpv2 "github.com/blorticus/gtpv2"
    "net" 
    "fmt"
)

modifyBearerRequest := gtpv2.PDU{
    Type:                     gtpv2.ModifyBearerRequest,
    IsCarryingPiggybackedPDU: false,
    PriorityFieldIsPresent:   false,
    TEIDFieldIsPresent:       true,
    SequenceNumber:           0x00001acc,
    Priority:                 0,
    TEID:                     0x05403b2e,
    TotalLength:              0x0042,
    InformationElements: []*gtpv2.IE{
        {
            Type:           gtpv2.UserLocationInformation,
            InstanceNumber: 0,
            TotalLength:    17,
            Data: []byte{
                0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
                0x00, 0x0f, 0x42, 0x4d, 0x00,
            },
        },
        {
            Type:           gtpv2.RATType,
            InstanceNumber: 0,
            TotalLength:    5,
            Data:           []byte{0x06},
        },
        {
            Type:           gtpv2.DelayValue,
            InstanceNumber: 0,
            TotalLength:    5,
            Data:           []byte{0x00},
        },
        {
            Type:           gtpv2.BearerContext,
            InstanceNumber: 0,
            TotalLength:    22,
            Data: []byte{
                0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
                0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13,
                0x01, 0xb2,
            },
        },
        {
            Type:           gtpv2.RecoveryRestartCounter,
            InstanceNumber: 0,
            TotalLength:    5,
            Data:           []byte{0x95},
        },
    },
}

conn, err := net.DialUDP("udp", "10.1.10.10:3386")
if err != nil {
    panic(err)
}

conn.Write(modifyBearerRequest.Encode())

incomingByteStream := make([]byte, 65536)

dgLength, err := conn.Read(incomingByteStream)
if err != nil {
    panic(err)
}

incomingGtpPDU, err := gtpv2.DecodePDU(incomingByteStream[:dgLength])
if err != nil {
    panic(err)
}

for _, ie := range incomingGtpPDU.InformationElements {
    fmt.Printf("IE name = (%s), value = (%02x)\n", gtpv2.NameOfIEForType(ie.Type), ie.Data)
}

# Information Elements

There is no support for interpretation of Information Elements.  They are created, stored, and presented as
a byte stream.  So, for example, when one uses the IMSI IE, it is not decoded to its natural string representation.

```
