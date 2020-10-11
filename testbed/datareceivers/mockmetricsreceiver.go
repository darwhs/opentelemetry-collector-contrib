package datareceivers

import (
	"context"
	"fmt"
	mockawsemfreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/testbed/mockdatareceivers/mockawsemfreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/testbed/testbed"
	"go.uber.org/zap"
)

// MockAwsEmfDataReceiver implements AwsXray format receiver.
type MockAwsEmfDataReceiver struct {
	testbed.DataReceiverBase
	receiver component.MetricsReceiver
}


// NewMockMetricDataReceiver creates a new  mockmetricsDatareceiver
func NewMockMetricDataReceiver(port int) *MockAwsEmfDataReceiver {
	return &MockAwsEmfDataReceiver{DataReceiverBase: testbed.DataReceiverBase{Port: port}}
}

//Start listening on the specified port
func (ar *MockAwsEmfDataReceiver) Start(_ consumer.TraceConsumer,tm consumer.MetricsConsumer, _ consumer.LogsConsumer) error {
	var err error
	mockDatareceiverCFG := mockawsemfreceiver.Config{
		Endpoint: fmt.Sprintf("localhost:%d", ar.Port),
		TLSCredentials: &configtls.TLSSetting{
			CertFile: "../mockdatareceivers/mockawsemfreceiver/server.crt",
			KeyFile:  "../mockdatareceivers/mockawsemfreceiver/server.key",
		},
	}
	params := component.ReceiverCreateParams{Logger: zap.L()}
	ar.receiver, err = mockawsemfreceiver.New(tm, params, &mockDatareceiverCFG)

	if err != nil {
		return err
	}

	return ar.receiver.Start(context.Background(), ar)
}

func (ar *MockAwsEmfDataReceiver) Stop() error {
	ar.receiver.Shutdown(context.Background())
	return nil
}

func (ar *MockAwsEmfDataReceiver) GenConfigYAMLStr() string {
	// Note that this generates an exporter config for agent.
	return fmt.Sprintf(`
  awsemf:
    endpoint: localhost:%d
    region: "us-west-2"`, ar.Port)
}

func (ar *MockAwsEmfDataReceiver) ProtocolName() string {
	return "awsemf"
}