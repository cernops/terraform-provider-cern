package cern

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Auth struct {
	Token string `xml:"urn:NetworkDataTypes token"`
}

type DeviceInput struct {
	DeviceName         string          `xml:"urn:NetworkDataTypes DeviceName"`
	Location           Location        `xml:"urn:NetworkDataTypes Location"`
	Zone               string          `xml:"urn:NetworkDataTypes Zone,omitempty"`
	Manufacturer       string          `xml:"urn:NetworkDataTypes Manufacturer"`
	Model              string          `xml:"urn:NetworkDataTypes Model"`
	Description        string          `xml:"urn:NetworkDataTypes Description,omitempty"`
	Tag                string          `xml:"urn:NetworkDataTypes Tag,omitempty"`
	SerialNumber       string          `xml:"urn:NetworkDataTypes SerialNumber,omitempty"`
	OperatingSystem    OperatingSystem `xml:"urn:NetworkDataTypes OperatingSystem"`
	InventoryNumber    string          `xml:"urn:NetworkDataTypes InventoryNumber,omitempty"`
	LandbManagerPerson PersonInput     `xml:"urn:NetworkDataTypes LandbManagerPerson,omitempty"`
	ResponsiblePerson  PersonInput     `xml:"urn:NetworkDataTypes ResponsiblePerson"`
	UserPerson         PersonInput     `xml:"urn:NetworkDataTypes UserPerson,omitempty"`
	HCPResponse        bool            `xml:"urn:NetworkDataTypes HCPResponse,omitempty"`
	IPv6Ready          bool            `xml:"urn:NetworkDataTypes IPv6Ready,omitempty"`
	ManagerLocked      bool            `xml:"urn:NetworkDataTypes ManagerLocked,omitempty"`
}

type InterfaceCard struct {
	HardwareAddress string `xml:"urn:NetworkDataTypes HardwareAddress"`
	CardType        string `xml:"urn:NetworkDataTypes CardType"`
}

type Location struct {
	Building string `xml:"urn:NetworkDataTypes Building"`
	Floor    string `xml:"urn:NetworkDataTypes Floor"`
	Room     string `xml:"urn:NetworkDataTypes Room"`
}

type OperatingSystem struct {
	Name    string `xml:"urn:NetworkDataTypes Name"`
	Version string `xml:"urn:NetworkDataTypes Version"`
}

type PersonInput struct {
	Name       string `xml:"urn:NetworkDataTypes Name,omitempty"`
	FirstName  string `xml:"urn:NetworkDataTypes FirstName,omitempty"`
	Department string `xml:"urn:NetworkDataTypes Department,omitempty"`
	Group      string `xml:"urn:NetworkDataTypes Group,omitempty"`
	PersonID   int64  `xml:"urn:NetworkDataTypes PersonID,omitempty"`
}

type VMCreateOptions struct {
	VMParent string `xml:"urn:NetworkDataTypes VMParent,omitempty"`
}

type VMInterfaceOptions struct {
	IP                   string `xml:"urn:NetworkDataTypes IP,omitempty"`
	IPv6                 string `xml:"urn:NetworkDataTypes IPv6,omitempty"`
	ServiceName          string `xml:"urn:NetworkDataTypes ServiceName,omitempty"`
	InternetConnectivity string `xml:"urn:NetworkDataTypes InternetConnectivity,omitempty"`
	AddressType          string `xml:"urn:NetworkDataTypes AddressType,omitempty"`
	BindHardwareAddress  string `xml:"urn:NetworkDataTypes BindHardwareAddress,omitempty"`
}

type LandbClient struct {
	HTTPClient   *http.Client
	ResponseHook func(*http.Response) *http.Response
	RequestHook  func(*http.Request) *http.Request
	Endpoint     string
	Auth         Auth
}

type soapEnvelope struct {
	XMLName struct{} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  struct {
		Auth interface{} `xml:",any"`
	} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body struct {
		Message interface{} `xml:",any"`
		Fault   *struct {
			String string `xml:"faultstring,omitempty"`
			Code   string `xml:"faultcode,omitempty"`
			Detail string `xml:"detail,omitempty"`
		} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
	} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

//NewLandbClient initialises a new connection to LanDB
func NewLandbClient(endpoint string, username string, password string) (*LandbClient, error) {
	client := LandbClient{
		Endpoint: endpoint,
	}
	token, err := client.GetAuthToken(context.TODO(), username, password, "CERN")
	if err != nil {
		return nil, fmt.Errorf("Error requesting Landb auth token: %s", err)
	}
	client.Auth = Auth{
		Token: token,
	}
	return &client, nil
}

func (c *LandbClient) do(ctx context.Context, method, action string, in, out interface{}) error {
	var body io.Reader
	var envelope soapEnvelope
	if method == "POST" || method == "PUT" {
		var buf bytes.Buffer
		envelope.Body.Message = in
		envelope.Header.Auth = &c.Auth
		enc := xml.NewEncoder(&buf)
		if err := enc.Encode(envelope); err != nil {
			return err
		}
		if err := enc.Flush(); err != nil {
			return err
		}
		body = &buf
	}
	req, err := http.NewRequest(method, c.Endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("SOAPAction", action)

	req = req.WithContext(ctx)

	if c.RequestHook != nil {
		req = c.RequestHook(req)
	}

	httpClient := c.HTTPClient

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	rsp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if c.ResponseHook != nil {
		rsp = c.ResponseHook(rsp)
	}
	dec := xml.NewDecoder(rsp.Body)
	envelope.Body.Message = out
	if err := dec.Decode(&envelope); err != nil {
		println("Error decoding")
		return err
	}
	if envelope.Body.Fault != nil {
		return fmt.Errorf("%s: %s", envelope.Body.Fault.Code, envelope.Body.Fault.String)
	}
	return nil
}

//GetAuthToken  gets authentication token from login and password.
func (c *LandbClient) GetAuthToken(ctx context.Context, Login string, Password string, Type string) (string, error) {
	var input struct {
		XMLName  struct{} `xml:"urn:NetworkService getAuthToken"`
		Login    string   `xml:"urn:NetworkService Login"`
		Password string   `xml:"urn:NetworkService Password"`
		Type     string   `xml:"urn:NetworkService Type"`
	}
	input.Login = string(Login)
	input.Password = string(Password)
	input.Type = string(Type)
	var output struct {
		XMLName struct{} `xml:"getAuthTokenResponse"`
		Token   string   `xml:"token"`
	}
	err := c.do(ctx, "POST", "", &input, &output)

	return string(output.Token), err
}

// VMCreate creates a new Virtual Machine
func (c *LandbClient) VMCreate(ctx context.Context, vmDevice DeviceInput, vmCreateOptions VMCreateOptions) (bool, error) {
	var input struct {
		XMLName         struct{}        `xml:"urn:NetworkService vmCreate"`
		VMDevice        DeviceInput     `xml:"urn:NetworkService VMDevice"`
		VMCreateOptions VMCreateOptions `xml:"urn:NetworkService VMCreateOptions"`
	}
	input.VMDevice = DeviceInput(vmDevice)
	input.VMCreateOptions = VMCreateOptions(vmCreateOptions)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmCreateResponse"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}

//VMUpdate updates basic information on virtual machine
func (c *LandbClient) VMUpdate(ctx context.Context, deviceName string, deviceInput DeviceInput) (bool, error) {
	var input struct {
		XMLName     struct{}    `xml:"urn:NetworkService vmUpdate"`
		DeviceName  string      `xml:"urn:NetworkService DeviceName"`
		DeviceInput DeviceInput `xml:"urn:NetworkService DeviceInput"`
	}
	input.DeviceName = string(deviceName)
	input.DeviceInput = DeviceInput(deviceInput)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmUpdate"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}

//VMDestroy removes completely a virtual machine from database
func (c *LandbClient) VMDestroy(ctx context.Context, vmName string) (bool, error) {
	var input struct {
		XMLName struct{} `xml:"urn:NetworkService vmDestroy"`
		VMName  string   `xml:"urn:NetworkService VMName"`
	}
	input.VMName = string(vmName)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmDestroyResponse"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}

//VMAddInterfaceRequest defines the interface fields
type VMAddInterfaceRequest struct {
	VMName             string
	InterfaceName      string
	VMClusterName      string
	VMInterfaceOptions VMInterfaceOptions
}

//VMAddInterface creates an IP interface for a virtual machine
func (c *LandbClient) VMAddInterface(ctx context.Context, v VMAddInterfaceRequest) (bool, error) {
	var input struct {
		XMLName            struct{}           `xml:"urn:NetworkService vmAddInterface"`
		VMName             string             `xml:"urn:NetworkService VMName"`
		InterfaceName      string             `xml:"urn:NetworkService InterfaceName"`
		VMClusterName      string             `xml:"urn:NetworkService VMClusterName"`
		VMInterfaceOptions VMInterfaceOptions `xml:"urn:NetworkService VMInterfaceOptions"`
	}

	input.VMName = string(v.VMName)
	input.InterfaceName = string(v.InterfaceName)
	input.VMClusterName = string(v.VMClusterName)
	input.VMInterfaceOptions = VMInterfaceOptions(v.VMInterfaceOptions)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmAddInterfaceResponse"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}

//VMRemoveInterface removes an IP interface from a virtual machine
func (c *LandbClient) VMRemoveInterface(ctx context.Context, vmName string, interfaceName string) (bool, error) {
	var input struct {
		XMLName       struct{} `xml:"urn:NetworkService vmRemoveInterface"`
		VMName        string   `xml:"urn:NetworkService VMName"`
		InterfaceName string   `xml:"urn:NetworkService InterfaceName"`
	}
	input.VMName = string(vmName)
	input.InterfaceName = string(interfaceName)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmRemoveInterfaceResponse"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}

//VMAddCard attaches a new hardware address to the VM given. A 3 octet prefix can be used and the remainder will be auto-generated.
func (c *LandbClient) VMAddCard(ctx context.Context, vmName string, interfaceCard InterfaceCard) (string, error) {
	var input struct {
		XMLName       struct{}      `xml:"urn:NetworkService vmAddCard"`
		VMName        string        `xml:"urn:NetworkService VMName"`
		InterfaceCard InterfaceCard `xml:"urn:NetworkService InterfaceCard"`
	}
	input.VMName = string(vmName)
	input.InterfaceCard = InterfaceCard(interfaceCard)
	var output struct {
		XMLName         struct{} `xml:"urn:NetworkService vmAddCardResponse"`
		HardwareAddress string   `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return string(output.HardwareAddress), err
}

//VMRemoveCard  detaches the specified interface card from the VM and removes it from the database.
func (c *LandbClient) VMRemoveCard(ctx context.Context, vmName string, hardwareAddress string) (bool, error) {
	var input struct {
		XMLName         struct{} `xml:"urn:NetworkService vmRemoveCard"`
		VMName          string   `xml:"urn:NetworkService VMName"`
		HardwareAddress string   `xml:"urn:NetworkService HardwareAddress"`
	}
	input.VMName = string(vmName)
	input.HardwareAddress = string(hardwareAddress)
	var output struct {
		XMLName struct{} `xml:"urn:NetworkService vmRemoveCardResponse"`
		Result  bool     `xml:",any"`
	}
	err := c.do(ctx, "POST", "", &input, &output)
	return bool(output.Result), err
}
