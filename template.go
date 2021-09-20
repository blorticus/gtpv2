package gtpv2

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

var mapOfYamlPduTypeToMessageType = map[string]MessageType{
	"EchoRequest":                                EchoRequest,
	"EchoResponse":                               EchoResponse,
	"CreateSessionRequest":                       CreateSessionRequest,
	"CreateSessionResponse":                      CreateSessionResponse,
	"ModifyBearerRequest":                        ModifyBearerRequest,
	"ModifyBearerResponse":                       ModifyBearerResponse,
	"DeleteSessionRequest":                       DeleteSessionRequest,
	"DeleteSessionResponse":                      DeleteSessionResponse,
	"RemoteUEReportNotification":                 RemoteUEReportNotification,
	"RemoteUEReportAcknowlegement":               RemoteUEReportAcknowlegement,
	"ChangeNotificationRequest":                  ChangeNotificationRequest,
	"ChangeNotificationResponse":                 ChangeNotificationResponse,
	"ModifyBearerCommand":                        ModifyBearerCommand,
	"ModifyBearerFailureIndication":              ModifyBearerFailureIndication,
	"DeleteBearerCommand":                        DeleteBearerCommand,
	"DeleteBearerFailureIndication":              DeleteBearerFailureIndication,
	"BearerResourceCommand":                      BearerResourceCommand,
	"BearerResourceFailureIndication":            BearerResourceFailureIndication,
	"DownlinkDataNotificationFailureIndication":  DownlinkDataNotificationFailureIndication,
	"TraceSessionActivation":                     TraceSessionActivation,
	"TraceSessionDeactivation":                   TraceSessionDeactivation,
	"StopPagingIndication":                       StopPagingIndication,
	"CreateBearerRequest":                        CreateBearerRequest,
	"CreateBearerResponse":                       CreateBearerResponse,
	"UpdateBearerRequest":                        UpdateBearerRequest,
	"UpdateBearerResponse":                       UpdateBearerResponse,
	"DeleteBearerRequest":                        DeleteBearerRequest,
	"DeleteBearerResponse":                       DeleteBearerResponse,
	"DeletePDNConnectionSetRequest":              DeletePDNConnectionSetRequest,
	"DeletePDNConnectionSetResponse":             DeletePDNConnectionSetResponse,
	"PGWDownlinkTriggeringNotification":          PGWDownlinkTriggeringNotification,
	"PGWDownlinkTriggeringAcknowledge":           PGWDownlinkTriggeringAcknowledge,
	"IdentificationRequest":                      IdentificationRequest,
	"IdentificationResponse":                     IdentificationResponse,
	"ContextRequest":                             ContextRequest,
	"ContextResponse":                            ContextResponse,
	"ContextAcknowledge":                         ContextAcknowledge,
	"ForwardRelocationRequest":                   ForwardRelocationRequest,
	"ForwardRelocationResponse":                  ForwardRelocationResponse,
	"ForwardRelocationCompleteNotification":      ForwardRelocationCompleteNotification,
	"ForwardRelocationCompleteAcknowledge":       ForwardRelocationCompleteAcknowledge,
	"ForwardAccessContextNotification":           ForwardAccessContextNotification,
	"ForwardAccessContextAcknowledge":            ForwardAccessContextAcknowledge,
	"RelocationCancelRequest":                    RelocationCancelRequest,
	"RelocationCancelResponse":                   RelocationCancelResponse,
	"ConfigurationTransferTunnel":                ConfigurationTransferTunnel,
	"DetachNotification":                         DetachNotification,
	"DetachAcknowledge":                          DetachAcknowledge,
	"CSPagingIndication":                         CSPagingIndication,
	"RANInformationRelay":                        RANInformationRelay,
	"AlertMMENotification":                       AlertMMENotification,
	"AlertMMEAcknowledge":                        AlertMMEAcknowledge,
	"UEActivityNotification":                     UEActivityNotification,
	"UEActivityAcknowledge":                      UEActivityAcknowledge,
	"ISRStatusIndication":                        ISRStatusIndication,
	"UERegistrationQueryRequest":                 UERegistrationQueryRequest,
	"UERegistrationQueryResponse":                UERegistrationQueryResponse,
	"CreateForwardingTunnelRequest":              CreateForwardingTunnelRequest,
	"CreateForwardingTunnelResponse":             CreateForwardingTunnelResponse,
	"SuspendNotification":                        SuspendNotification,
	"SuspendAcknowledge":                         SuspendAcknowledge,
	"ResumeNotification":                         ResumeNotification,
	"ResumeAcknowledge":                          ResumeAcknowledge,
	"CreateIndirectDataForwardingTunnelRequest":  CreateIndirectDataForwardingTunnelRequest,
	"CreateIndirectDataForwardingTunnelResponse": CreateIndirectDataForwardingTunnelResponse,
	"DeleteIndirectDataForwardingTunnelRequest":  DeleteIndirectDataForwardingTunnelRequest,
	"DeleteIndirectDataForwardingTunnelResponse": DeleteIndirectDataForwardingTunnelResponse,
	"ReleaseAccessBearersRequest":                ReleaseAccessBearersRequest,
	"ReleaseAccessBearersResponse":               ReleaseAccessBearersResponse,
	"DownlinkDataNotification":                   DownlinkDataNotification,
	"DownlinkDataNotificationAcknowledge":        DownlinkDataNotificationAcknowledge,
	"PGWRestartNotification":                     PGWRestartNotification,
	"PGWRestartNotificationAcknowledge":          PGWRestartNotificationAcknowledge,
	"UpdatePDNConnectionSetRequest":              UpdatePDNConnectionSetRequest,
	"UpdatePDNConnectionSetResponse":             UpdatePDNConnectionSetResponse,
}

type IEYaml struct {
	Type  string      `yaml:"Type"`
	Value interface{} `yaml:"Value"`
}

type Gtpv2PduYaml struct {
	Name string   `yaml:"Name"`
	Type string   `yaml:"Type"`
	IEs  []IEYaml `yaml:"IEs"`
}

type GtpDefinitionRootYaml struct {
	Gtpv2Pdus []Gtpv2PduYaml
}

func validateGtpv2PduYaml(yaml Gtpv2PduYaml) error {
	if _, providedPduTypeIsValid := mapOfYamlPduTypeToMessageType[yaml.Name]; !providedPduTypeIsValid {
		return fmt.Errorf("provided PDU Type (%s) is not recognized", yaml.Name)
	}

	return nil
}

// t := gtpv2.ReadYamlTemplateFromFile( "/path/to/file.yaml" )
// csr := t.GeneratePDUByName( "CSR01" )
type Template struct {
	mapOfGtpv2PduYamlByName map[string]Gtpv2PduYaml
}

func ReadYamlTemplateFromString(yamlDefinition string) (*Template, error) {
	unmarhalledYaml := &GtpDefinitionRootYaml{}
	if err := yaml.Unmarshal([]byte(yamlDefinition), &unmarhalledYaml); err != nil {
		return nil, err
	}

	if unmarhalledYaml == nil {
		return nil, fmt.Errorf("undefined YAML parsing error")
	}

	mapOfGtpv2PduYamlByName := make(map[string]Gtpv2PduYaml)

	for _, pduDefinitionYaml := range unmarhalledYaml.Gtpv2Pdus {
		if err := validateGtpv2PduYaml(pduDefinitionYaml); err != nil {
			return nil, err
		}

		mapOfGtpv2PduYamlByName[pduDefinitionYaml.Name] = pduDefinitionYaml
	}

	return &Template{
		mapOfGtpv2PduYamlByName: mapOfGtpv2PduYamlByName,
	}, nil
}

func ReadYamlTemplateFromFile(filePath string) (*Template, error) {
	return nil, nil
}

// Gtpv2Pdus:
//   - Name: CreateSessionRequest
//     Type: CreateSessionRequest
//     IEs:
//       - Type: "IMSI"
//         Value: "{{ Requestor.IMSI }}"
//       - Type: "MEI"
//         Value: "01-124500-012345-01"
//       - Type: "ServingNetwork"
//         Value: "001001"
//       - Type: "F-TEID"
//         Value: ...
//       - Type: "F-TEID"
//         Value: ...
//       - Type: "APN"
//         Value: "{{ Requestor.APN }}"
//       - Type: "PDNType"
//         Value: "{{ Requestor.PDNType }}"
//       - Type: "BearerContext"
//         Value:
//           - Type: "EBI"
//             Value: ...
//           - Type: "F-TEID"
//             Value: ...
//           - Type: "F-TEID"
//             Value: ...
//           -
