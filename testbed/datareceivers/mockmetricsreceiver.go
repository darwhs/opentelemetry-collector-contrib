package datareceivers

import (
	"context"
	"crypto/x509"
	"fmt"
	mockawsemfreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/testbed/mockdatareceivers/mockawsemfreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/testbed/testbed"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
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
func (ar *MockAwsEmfDataReceiver) Start(_ consumer.TracesConsumer,tm consumer.MetricsConsumer, _ consumer.LogsConsumer) error {

	var err error
	os.Setenv("AWS_ACCESS_KEY_ID", "placeholder")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "placeholder")

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile("../mockdatareceivers/mockawsxrayreceiver/server.crt")

	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", "../mockdatareceivers/mockawsxrayreceiver/server.crt", err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}
	mockDatareceiverCFG := mockawsemfreceiver.Config{
		Endpoint: fmt.Sprintf("localhost:%d", ar.Port),
		TLSCredentials: &configtls.TLSSetting{
			CertFile: "../mockdatareceivers/mockawsxrayreceiver/server.crt",
			KeyFile:  "../mockdatareceivers/mockawsxrayreceiver/server.key",
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
    no_verify_ssl: true
    region: "us-west-2"`, ar.Port)
}

func (ar *MockAwsEmfDataReceiver) ProtocolName() string {
	return "awsemf"
}