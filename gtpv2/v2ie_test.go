package v2

import (
	"fmt"
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
		v2IENamesComparable{"Reserved", 0},
		v2IENamesComparable{"Cause", 2},
		v2IENamesComparable{"TMSI", 88},
		v2IENamesComparable{"P-TMSI", 111},
		v2IENamesComparable{"Throttling", 154},
		v2IENamesComparable{"UP Function Selection Indication Flags", 202},
	}

	for _, testCase := range testCases {
		if got := NameOfIEForType(testCase.ieType); got != testCase.expectedName {
			t.Errorf("For IE type (%d) expected name string = (%s), got = (%s)", testCase.ieType, testCase.expectedName, got)
		}
	}
}

func TestV2IEDecodeInvalidCases(t *testing.T) {
	cases := []v2IEFailCase{
		v2IEFailCase{
			name:        "Empty stream",
			inputStream: []byte{},
		},
		v2IEFailCase{
			name:        "Stream too short for header",
			inputStream: []byte{0x01, 0x00, 0x06},
		},
		v2IEFailCase{
			name:        "Header only",
			inputStream: []byte{0x01, 0x00, 0x06, 0x00},
		},
		v2IEFailCase{
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
		v2IEComparable{
			ieOctets: []byte{0x56, 0x00, 0x0d, 0x00, 0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			matchingIE: &IE{
				Type:           UserLocationInformation,
				TotalLength:    17,
				InstanceNumber: 0,
				Data:           []byte{0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			},
		},
		v2IEComparable{
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
		v2IEComparable{
			ieOctets: []byte{0x56, 0x00, 0x0d, 0x00, 0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			matchingIE: &IE{
				Type:           UserLocationInformation,
				TotalLength:    17,
				InstanceNumber: 0,
				Data:           []byte{0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			},
		},
		v2IEComparable{
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
