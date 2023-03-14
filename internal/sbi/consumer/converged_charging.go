package consumer

import (
	"context"
	"time"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nchf_ConvergedCharging"
	"github.com/free5gc/openapi/models"
)

func SendConvergedChargingRequest(ue *amf_context.AmfUe) (*models.ChargingDataResponse, *models.ProblemDetails, error) {
	configuration := Nchf_ConvergedCharging.NewConfiguration()
	configuration.SetBasePath(ue.ChfUri)

	client := Nchf_ConvergedCharging.NewAPIClient(configuration)

	amfSelf := amf_context.AMF_Self()

	date := time.Now()

	req := models.ChargingDataRequest{
		SubscriberIdentifier: ue.Supi,
		NfConsumerIdentification: &models.NfIdentification{
			NodeFunctionality: models.NodeFunctionality_AMF,
			NFName:            amfSelf.NfId,
		},
		ChargingId:               ue.ChargingId,
		InvocationTimeStamp:      &date,
		InvocationSequenceNumber: 0,
		OneTimeEvent:             true,
		OneTimeEventType:         models.OneTimeEventType_PEC,
		ServiceSpecificationInfo: "TS 32.256 v16.10.0",
		RegistrationChargingInformation: &models.RegistrationChargingInformation{
			RegistrationMessagetype: models.RegistrationMessageType_INITIAL,
		},
	}

	rsp, httpResponse, err := client.DefaultApi.ChargingdataPost(context.Background(), req)
	if err == nil {
		return &rsp, nil, nil
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			return nil, nil, err
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return nil, &problem, nil
	} else {
		return nil, nil, openapi.ReportError("server no response")
	}
}
