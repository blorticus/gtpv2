package gtpv2

import (
	"fmt"
	"net"
	"testing"
)

type v2IEComparable struct {
	ieOctets   []byte
	matchingIE *IE
}

type v2IEFailCase struct {
	name        string
	inputStream []byte
}

type v2IENamesComparable struct {
	expectedName string
	ieType       IEType
}

func TestIENames(t *testing.T) {
	// This test set is mostly to make sure the list doesn't accidentally
	// get shifted if values are changed
	testCases := []v2IENamesComparable{
		{"Reserved", 0},
		{"Cause", 2},
		{"TMSI", 88},
		{"P-TMSI", 111},
		{"Throttling", 154},
		{"UP Function Selection Indication Flags", 202},
	}

	for _, testCase := range testCases {
		if got := NameOfIEForType(testCase.ieType); got != testCase.expectedName {
			t.Errorf("For IE type (%d) expected name string = (%s), got = (%s)", testCase.ieType, testCase.expectedName, got)
		}
	}
}

func TestV2IEDecodeInvalidCases(t *testing.T) {
	cases := []v2IEFailCase{
		{
			name:        "Empty stream",
			inputStream: []byte{},
		},
		{
			name:        "Stream too short for header",
			inputStream: []byte{0x01, 0x00, 0x06},
		},
		{
			name:        "Header only",
			inputStream: []byte{0x01, 0x00, 0x06, 0x00},
		},
		{
			name:        "Insufficient byte stream length",
			inputStream: []byte{0x01, 0x00, 0x06, 0x00, 0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, testCase := range cases {
		_, err := DecodeIE(testCase.inputStream)

		if err == nil {
			t.Errorf("(%s) Expected error on DecodeV2IE(), but received none", testCase.name)
		}
	}
}

func TestV2IEDecodeValidCases(t *testing.T) {
	cases := []v2IEComparable{
		{
			ieOctets: []byte{0x56, 0x00, 0x0d, 0x00, 0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			matchingIE: &IE{
				Type:           UserLocationInformation,
				TotalLength:    17,
				InstanceNumber: 0,
				Data:           []byte{0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			},
		},
		{
			ieOctets: []byte{0x52, 0x00, 0x01, 0x03, 0x06},
			matchingIE: &IE{
				Type:           RATType,
				TotalLength:    5,
				InstanceNumber: 3,
				Data:           []byte{0x06},
			},
		},
	}

	testCaseNumber := 0
	for _, testCase := range cases {
		testCaseNumber++
		ie, err := DecodeIE(testCase.ieOctets)

		if err != nil {
			t.Errorf("For test case number [%d], received error on decode: %s", testCaseNumber, err)
			continue
		}

		if err = compareTwoIEObjects(testCase.matchingIE, ie); err != nil {
			t.Errorf("For test case number [%d]: %s", testCaseNumber, err)
		}
	}
}

func TestIEEncodeValidCases(t *testing.T) {
	testCases := []v2IEComparable{
		{
			ieOctets: []byte{0x56, 0x00, 0x0d, 0x00, 0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			matchingIE: &IE{
				Type:           UserLocationInformation,
				TotalLength:    17,
				InstanceNumber: 0,
				Data:           []byte{0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			},
		},
		{
			ieOctets: []byte{0x52, 0x00, 0x01, 0x03, 0x06},
			matchingIE: &IE{
				Type:           RATType,
				TotalLength:    5,
				InstanceNumber: 3,
				Data:           []byte{0x06},
			},
		},
	}

	for _, testCase := range testCases {
		encodedBytes := testCase.matchingIE.Encode()
		if err := compareByteArrays(testCase.ieOctets, encodedBytes); err != nil {
			t.Errorf("Encoded() did not generate expected byte stream: %s", err)
		}
	}
}

type TypedFTEIDComparable struct {
	fteid             *TypedFTEID
	expectedDataBytes []byte
}

func TestTypedFTEID(t *testing.T) {
	testCases := []TypedFTEIDComparable{
		{
			fteid:             &TypedFTEID{IPv4Addr: nil, IPv6Addr: nil, InterfaceType: 1, Key: 0xaabbccdd},
			expectedDataBytes: []byte{0x01, 0xaa, 0xbb, 0xcc, 0xdd},
		},
		{
			fteid:             &TypedFTEID{IPv4Addr: net.ParseIP("10.11.12.13"), IPv6Addr: nil, InterfaceType: 1, Key: 0xaabbccdd},
			expectedDataBytes: []byte{0x81, 0xaa, 0xbb, 0xcc, 0xdd, 0x0a, 0x0b, 0x0c, 0x0d},
		},
		{
			fteid:             &TypedFTEID{IPv4Addr: nil, IPv6Addr: net.ParseIP("fd00:a:b:c:d::1"), InterfaceType: 1, Key: 0xaabbccdd},
			expectedDataBytes: []byte{0x41, 0xaa, 0xbb, 0xcc, 0xdd, 0xfd, 0x00, 0x00, 0x0a, 0x00, 0x0b, 0x00, 0x0c, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			fteid:             &TypedFTEID{IPv4Addr: net.ParseIP("10.11.12.13"), IPv6Addr: net.ParseIP("fd00:a:b:c:d::1"), InterfaceType: 1, Key: 0xaabbccdd},
			expectedDataBytes: []byte{0xc1, 0xaa, 0xbb, 0xcc, 0xdd, 0x0a, 0x0b, 0x0c, 0x0d, 0xfd, 0x00, 0x00, 0x0a, 0x00, 0x0b, 0x00, 0x0c, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
	}

	for testIndex, testCase := range testCases {
		ie, err := testCase.fteid.ToIEErrorable()

		if err != nil {
			t.Errorf("On TestIE_FTEID test number [%d], error on ToIE: %s", testIndex+1, err.Error())
		} else {
			if err := compareByteArrays(testCase.expectedDataBytes, ie.Data); err != nil {
				t.Errorf("On TestIE_FTEID test number [%d], data bytes from ToIE do not match: %s", testIndex+1, err.Error())
			}

			ie, err := NewIEWithRawDataErrorable(FTEID, testCase.expectedDataBytes)

			if err != nil {
				t.Errorf("On TestIE_FTEID test number [%d], unable to convert expectedDataBytes to IE: %s", testIndex+1, err.Error())
			} else {
				typedFteid, err := ie.TypedDataErrorable()

				if err != nil {
					t.Errorf("On TestIE_FTEID test number [%d], error on TypedData: %s", testIndex+1, err.Error())
				} else {
					fteid := typedFteid.(*TypedFTEID)

					if fteid.IPv4Addr == nil && testCase.fteid.IPv4Addr != nil {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, got nil ipv4Addr, but it should not be", testIndex+1)
					} else if fteid.IPv4Addr != nil && testCase.fteid.IPv4Addr == nil {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, got ipv4Addr, but it should be nil", testIndex+1)
					} else if !fteid.IPv4Addr.Equal(testCase.fteid.IPv4Addr) {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, resulting fteid.ipv4Addr is not expected value", testIndex+1)
					}

					if fteid.IPv6Addr == nil && testCase.fteid.IPv6Addr != nil {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, got nil ipv6Addr, but it should not be", testIndex+1)
					} else if fteid.IPv6Addr != nil && testCase.fteid.IPv6Addr == nil {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, got ipv6Addr, but it should be nil", testIndex+1)
					} else if !fteid.IPv6Addr.Equal(testCase.fteid.IPv6Addr) {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, resulting fteid.ipv6Addr is not expected value", testIndex+1)
					}

					if fteid.InterfaceType != testCase.fteid.InterfaceType {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, expected interfaceType [%d], got[%d]", testIndex+1, testCase.fteid.InterfaceType, fteid.InterfaceType)
					}

					if fteid.Key != testCase.fteid.Key {
						t.Errorf("On TestIE_FTEID test number [%d], on TypedData, expected key [%08x], got[%08x]", testIndex+1, testCase.fteid.Key, fteid.Key)
					}
				}
			}
		}
	}
}

type TypedIMSIComparable struct {
	imsi              *TypedIMSI
	expectedDataBytes []byte
}

func TestTypedIMSI(t *testing.T) {
	testCases := []TypedIMSIComparable{
		{
			imsi:              &TypedIMSI{AsString: "001002789012345"},
			expectedDataBytes: []byte{0x00, 0x01, 0x20, 0x87, 0x09, 0x21, 0x43, 0xf5},
		},
	}

	for testIndex, testCase := range testCases {
		testNumber := testIndex + 1

		ie, err := testCase.imsi.ToIEErrorable()

		if err != nil {
			t.Errorf("[TestTypedIMSI] on test number [%d] did not expect error, but got error = (%s)", testNumber, err.Error())
		} else {
			if err := compareByteArrays(testCase.expectedDataBytes, ie.Data); err != nil {
				t.Errorf("[TestTypeIMSI] on test number [%d] data in IE from ToIEErrorable does not match expected: %s", testNumber, err.Error())
			}
		}

		typedImsi, err := ie.TypedDataErrorable()

		if err != nil {
			t.Errorf("[TestTypeIMSI] on test number [%d] expected no error on TypedData but got error = (%s)", testNumber, err.Error())
		} else {
			imsi := typedImsi.(*TypedIMSI)

			if imsi.AsString != testCase.imsi.AsString {
				t.Errorf("[TestTypeIMSI] on test number [%d], expected AsString = (%s), got (%s)", testNumber, testCase.imsi.AsString, imsi.AsString)
			}
		}
	}
}

func TestGroupIECreation(t *testing.T) {
	ie, err := NewGroupedIEErrorable(BearerContext, []*IE{
		NewIEWithRawData(EBI, []byte{0x01}),
		(&TypedFTEID{IPv4Addr: net.IPv4(10, 11, 12, 13), InterfaceType: 1, Key: 0xaabbccdd}).ToIE(),
		(&TypedFTEID{IPv4Addr: net.IPv4(1, 2, 3, 4), InterfaceType: 3, Key: 0x01020344}).ToIE(),
	})

	if err != nil {
		t.Errorf("[TestGroupIECreation] expected no error, but got error = (%s)", err.Error())
	} else {
		compareByteArrays([]byte{
			73, 0x01, 0x00, 0x01,
			87, 0x81, 0xaa, 0xbb, 0xcc, 0xdd, 0x0a, 0x0b, 0x0c, 0x0d,
			87, 0x83, 0x01, 0x02, 0x03, 0x44, 0x01, 0x02, 0x03, 0x04,
		}, ie.Data)
	}
}

func compareByteArrays(expected []byte, got []byte) error {
	if len(expected) != len(got) {
		return fmt.Errorf("Byte array lengths differ; expected %d bytes, got = %d", len(expected), len(got))
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			return fmt.Errorf("At index %d, expected = %02x, got = %02x", i, expected[i], got[i])
		}
	}

	return nil
}

func compareTwoIEObjects(expected *IE, got *IE) error {
	if expected.Type != got.Type {
		return fmt.Errorf("Expected IE Type [%d] (%s), got [%d] (%s)", expected.Type, NameOfIEForType(expected.Type), got.Type, NameOfIEForType(got.Type))
	}

	if expected.TotalLength != got.TotalLength {
		return fmt.Errorf("Expected IE TotalLength = %d, got = %d", expected.TotalLength, got.TotalLength)
	}

	if expected.InstanceNumber != got.InstanceNumber {
		return fmt.Errorf("Expected IE InstanceNumber = %d, got = %d", expected.InstanceNumber, got.InstanceNumber)
	}

	if len(expected.Data) != len(got.Data) {
		return fmt.Errorf("Expected IE Data has (%d) bytes, got = (%d) bytes", len(expected.Data), len(got.Data))
	}

	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != got.Data[i] {
			return fmt.Errorf("Expected IE Data byte index (%d) is (0x%02x), got (0x%02x)", i, expected.Data[i], got.Data[i])
		}
	}

	return nil
}
