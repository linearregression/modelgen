// contivModel.go
// This file is auto generated by modelgen tool
// Do not edit this file manually

package contivModel

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objdb/modeldb"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error)

type Tenant struct {
	// every object has a key
	Key string `json:"key,omitempty"`

	TenantName string `json:"tenantName,omitempty"` //

	// add link-sets and links
	LinkSets TenantLinkSets `json:"link-sets,omitempty"`
}

type TenantLinkSets struct {
	Networks map[string]modeldb.Link `json:"Networks,omitempty"`
}

type TenantInspect struct {
	Config Tenant
}

type Network struct {
	// every object has a key
	Key string `json:"key,omitempty"`

	Encap       string   `json:"encap,omitempty"`     //
	IsPrivate   bool     `json:"isPrivate,omitempty"` //
	IsPublic    bool     `json:"isPublic,omitempty"`  //
	Labels      []string `json:"labels,omitempty"`
	NetworkName string   `json:"networkName,omitempty"` //
	PktTag      int      `json:"pktTag,omitempty"`      //
	Policies    []string `json:"policies,omitempty"`
	Subnet      string   `json:"subnet,omitempty"`     //
	TenantName  string   `json:"tenantName,omitempty"` //

	Links NetworkLinks `json:"links,omitempty"`
}

type NetworkLinks struct {
	Tenant modeldb.Link `json:"Tenant,omitempty"`
}

type NetworkInspect struct {
	Config Network
}

type NetTwo struct {
	// every object has a key
	Key string `json:"key,omitempty"`

	Encap       string   `json:"encap,omitempty"`     //
	IsPrivate   bool     `json:"isPrivate,omitempty"` //
	IsPublic    bool     `json:"isPublic,omitempty"`  //
	Labels      []string `json:"labels,omitempty"`
	NetworkName string   `json:"networkName,omitempty"` //
	PktTag      int      `json:"pktTag,omitempty"`      //
	Policies    []string `json:"policies,omitempty"`
	Subnet      string   `json:"subnet,omitempty"`     //
	TenantName  string   `json:"tenantName,omitempty"` //

	Links NetTwoLinks `json:"links,omitempty"`
}

type NetTwoLinks struct {
	Tenant modeldb.Link `json:"Tenant,omitempty"`
}

type NetTwoOper struct {
	CreateTime  string `json:"createTime,omitempty"`  //
	NetworkName string `json:"networkName,omitempty"` //
	TenantName  string `json:"tenantName,omitempty"`  //

}

type NetTwoInspect struct {
	Config NetTwo

	Oper NetTwoOper
}

type EndpointOper struct {

	// oper object key (present for oper only objects)
	Key string `json:"key,omitempty"`

	Labels string `json:"labels,omitempty"` //
	UUID   string `json:"uuid,omitempty"`   //

}

type EndpointInspect struct {
	Oper EndpointOper
}

type EpListOper struct {

	// oper object key (present for oper only objects)
	Key string `json:"key,omitempty"`

	Eps    []EndpointOper `json:"eps,omitempty"`
	Name   string         `json:"name,omitempty"`   //
	Subnet string         `json:"subnet,omitempty"` //

}

type EpListInspect struct {
	Oper EpListOper
}

type NetWithEpOper struct {

	// oper object key (present for oper only objects)
	Key string `json:"key,omitempty"`

	Ep     EndpointOper `json:"ep,omitempty"`     //
	Name   string       `json:"name,omitempty"`   //
	Subnet string       `json:"subnet,omitempty"` //

}

type NetWithEpInspect struct {
	Oper NetWithEpOper
}
type Collections struct {
	tenants  map[string]*Tenant
	networks map[string]*Network
	netTwos  map[string]*NetTwo
}

var collections Collections

type TenantCallbacks interface {
	TenantCreate(tenant *Tenant) error
	TenantUpdate(tenant, params *Tenant) error
	TenantDelete(tenant *Tenant) error
}

type NetworkCallbacks interface {
	NetworkCreate(network *Network) error
	NetworkUpdate(network, params *Network) error
	NetworkDelete(network *Network) error
}

type NetTwoCallbacks interface {
	NetTwoGetOper(netTwo *NetTwoInspect) error

	NetTwoCreate(netTwo *NetTwo) error
	NetTwoUpdate(netTwo, params *NetTwo) error
	NetTwoDelete(netTwo *NetTwo) error
}

type EndpointCallbacks interface {
	EndpointGetOper(endpoint *EndpointInspect) error
}

type EpListCallbacks interface {
	EpListGetOper(epList *EpListInspect) error
}

type NetWithEpCallbacks interface {
	NetWithEpGetOper(netWithEp *NetWithEpInspect) error
}

type CallbackHandlers struct {
	TenantCb    TenantCallbacks
	NetworkCb   NetworkCallbacks
	NetTwoCb    NetTwoCallbacks
	EndpointCb  EndpointCallbacks
	EpListCb    EpListCallbacks
	NetWithEpCb NetWithEpCallbacks
}

var objCallbackHandler CallbackHandlers

func Init() {
	collections.tenants = make(map[string]*Tenant)
	collections.networks = make(map[string]*Network)
	collections.netTwos = make(map[string]*NetTwo)

	restoreTenant()
	restoreNetwork()
	restoreNetTwo()

}

func RegisterTenantCallbacks(handler TenantCallbacks) {
	objCallbackHandler.TenantCb = handler
}

func RegisterNetworkCallbacks(handler NetworkCallbacks) {
	objCallbackHandler.NetworkCb = handler
}

func RegisterNetTwoCallbacks(handler NetTwoCallbacks) {
	objCallbackHandler.NetTwoCb = handler
}

func RegisterEndpointCallbacks(handler EndpointCallbacks) {
	objCallbackHandler.EndpointCb = handler
}

func RegisterEpListCallbacks(handler EpListCallbacks) {
	objCallbackHandler.EpListCb = handler
}

func RegisterNetWithEpCallbacks(handler NetWithEpCallbacks) {
	objCallbackHandler.NetWithEpCb = handler
}

// Simple Wrapper for http handlers
func makeHttpHandler(handlerFunc HttpApiFunc) http.HandlerFunc {
	// Create a closure and return an anonymous function
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the handler
		resp, err := handlerFunc(w, r, mux.Vars(r))
		if err != nil {
			// Log error
			log.Errorf("Handler for %s %s returned error: %s", r.Method, r.URL, err)

			// Send HTTP response
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// Send HTTP response as Json
			err = writeJSON(w, http.StatusOK, resp)
			if err != nil {
				log.Errorf("Error generating json. Err: %v", err)
			}
		}
	}
}

// writeJSON: writes the value v to the http response stream as json with standard
// json encoding.
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	// Set content type as json
	w.Header().Set("Content-Type", "application/json")

	// write the HTTP status code
	w.WriteHeader(code)

	// Write the Json output
	return json.NewEncoder(w).Encode(v)
}

// Add all routes for REST handlers
func AddRoutes(router *mux.Router) {
	var route, listRoute, inspectRoute string

	// Register tenant
	route = "/api/v1/tenants/{key}/"
	listRoute = "/api/v1/tenants/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListTenants))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetTenant))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateTenant))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateTenant))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteTenant))

	inspectRoute = "/api/v1/inspect/tenants/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectTenant))

	// Register network
	route = "/api/v1/networks/{key}/"
	listRoute = "/api/v1/networks/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListNetworks))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetNetwork))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateNetwork))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateNetwork))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteNetwork))

	inspectRoute = "/api/v1/inspect/networks/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectNetwork))

	// Register netTwo
	route = "/api/v1/netTwos/{key}/"
	listRoute = "/api/v1/netTwos/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListNetTwos))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetNetTwo))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateNetTwo))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateNetTwo))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteNetTwo))

	inspectRoute = "/api/v1/inspect/netTwos/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectNetTwo))

	inspectRoute = "/api/v1/inspect/endpoints/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectEndpoint))

	inspectRoute = "/api/v1/inspect/epLists/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectEpList))

	inspectRoute = "/api/v1/inspect/netWithEps/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectNetWithEp))

}

// GET Oper REST call
func httpInspectTenant(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj TenantInspect
	log.Debugf("Received httpInspectTenant: %+v", vars)

	key := vars["key"]

	objConfig := collections.tenants[key]
	if objConfig == nil {
		log.Errorf("tenant %s not found", key)
		return nil, errors.New("tenant not found")
	}
	obj.Config = *objConfig

	// Return the obj
	return &obj, nil
}

// LIST REST call
func httpListTenants(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpListTenants: %+v", vars)

	list := make([]*Tenant, 0)
	for _, obj := range collections.tenants {
		list = append(list, obj)
	}

	// Return the list
	return list, nil
}

// GET REST call
func httpGetTenant(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetTenant: %+v", vars)

	key := vars["key"]

	obj := collections.tenants[key]
	if obj == nil {
		log.Errorf("tenant %s not found", key)
		return nil, errors.New("tenant not found")
	}

	// Return the obj
	return obj, nil
}

// CREATE REST call
func httpCreateTenant(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetTenant: %+v", vars)

	var obj Tenant
	key := vars["key"]

	// Get object from the request
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Errorf("Error decoding tenant create request. Err %v", err)
		return nil, err
	}

	// set the key
	obj.Key = key

	// Create the object
	err = CreateTenant(&obj)
	if err != nil {
		log.Errorf("CreateTenant error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return obj, nil
}

// DELETE rest call
func httpDeleteTenant(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpDeleteTenant: %+v", vars)

	key := vars["key"]

	// Delete the object
	err := DeleteTenant(key)
	if err != nil {
		log.Errorf("DeleteTenant error for: %s. Err: %v", key, err)
		return nil, err
	}

	// Return the obj
	return key, nil
}

// Create a tenant object
func CreateTenant(obj *Tenant) error {
	// Validate parameters
	err := ValidateTenant(obj)
	if err != nil {
		log.Errorf("ValidateTenant retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// Check if we handle this object
	if objCallbackHandler.TenantCb == nil {
		log.Errorf("No callback registered for tenant object")
		return errors.New("Invalid object type")
	}

	saveObj := obj

	// Check if object already exists
	if collections.tenants[obj.Key] != nil {
		// Perform Update callback
		err = objCallbackHandler.TenantCb.TenantUpdate(collections.tenants[obj.Key], obj)
		if err != nil {
			log.Errorf("TenantUpdate retruned error for: %+v. Err: %v", obj, err)
			return err
		}

		// save the original object after update
		saveObj = collections.tenants[obj.Key]
	} else {
		// save it in cache
		collections.tenants[obj.Key] = obj

		// Perform Create callback
		err = objCallbackHandler.TenantCb.TenantCreate(obj)
		if err != nil {
			log.Errorf("TenantCreate retruned error for: %+v. Err: %v", obj, err)
			delete(collections.tenants, obj.Key)
			return err
		}
	}

	// Write it to modeldb
	err = saveObj.Write()
	if err != nil {
		log.Errorf("Error saving tenant %s to db. Err: %v", saveObj.Key, err)
		return err
	}

	return nil
}

// Return a pointer to tenant from collection
func FindTenant(key string) *Tenant {
	obj := collections.tenants[key]
	if obj == nil {
		return nil
	}

	return obj
}

// Delete a tenant object
func DeleteTenant(key string) error {
	obj := collections.tenants[key]
	if obj == nil {
		log.Errorf("tenant %s not found", key)
		return errors.New("tenant not found")
	}

	// Check if we handle this object
	if objCallbackHandler.TenantCb == nil {
		log.Errorf("No callback registered for tenant object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.TenantCb.TenantDelete(obj)
	if err != nil {
		log.Errorf("TenantDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// delete it from modeldb
	err = obj.Delete()
	if err != nil {
		log.Errorf("Error deleting tenant %s. Err: %v", obj.Key, err)
	}

	// delete it from cache
	delete(collections.tenants, key)

	return nil
}

func (self *Tenant) GetType() string {
	return "tenant"
}

func (self *Tenant) GetKey() string {
	return self.Key
}

func (self *Tenant) Read() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to read tenant object")
		return errors.New("Empty key")
	}

	return modeldb.ReadObj("tenant", self.Key, self)
}

func (self *Tenant) Write() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Write tenant object")
		return errors.New("Empty key")
	}

	return modeldb.WriteObj("tenant", self.Key, self)
}

func (self *Tenant) Delete() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Delete tenant object")
		return errors.New("Empty key")
	}

	return modeldb.DeleteObj("tenant", self.Key)
}

func restoreTenant() error {
	strList, err := modeldb.ReadAllObj("tenant")
	if err != nil {
		log.Errorf("Error reading tenant list. Err: %v", err)
	}

	for _, objStr := range strList {
		// Parse the json model
		var tenant Tenant
		err = json.Unmarshal([]byte(objStr), &tenant)
		if err != nil {
			log.Errorf("Error parsing object %s, Err %v", objStr, err)
			return err
		}

		// add it to the collection
		collections.tenants[tenant.Key] = &tenant
	}

	return nil
}

// Validate a tenant object
func ValidateTenant(obj *Tenant) error {
	// Validate key is correct
	keyStr := obj.TenantName
	if obj.Key != keyStr {
		log.Errorf("Expecting Tenant Key: %s. Got: %s", keyStr, obj.Key)
		return errors.New("Invalid Key")
	}

	// Validate each field

	return nil
}

// GET Oper REST call
func httpInspectNetwork(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj NetworkInspect
	log.Debugf("Received httpInspectNetwork: %+v", vars)

	key := vars["key"]

	objConfig := collections.networks[key]
	if objConfig == nil {
		log.Errorf("network %s not found", key)
		return nil, errors.New("network not found")
	}
	obj.Config = *objConfig

	// Return the obj
	return &obj, nil
}

// LIST REST call
func httpListNetworks(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpListNetworks: %+v", vars)

	list := make([]*Network, 0)
	for _, obj := range collections.networks {
		list = append(list, obj)
	}

	// Return the list
	return list, nil
}

// GET REST call
func httpGetNetwork(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetNetwork: %+v", vars)

	key := vars["key"]

	obj := collections.networks[key]
	if obj == nil {
		log.Errorf("network %s not found", key)
		return nil, errors.New("network not found")
	}

	// Return the obj
	return obj, nil
}

// CREATE REST call
func httpCreateNetwork(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetNetwork: %+v", vars)

	var obj Network
	key := vars["key"]

	// Get object from the request
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Errorf("Error decoding network create request. Err %v", err)
		return nil, err
	}

	// set the key
	obj.Key = key

	// Create the object
	err = CreateNetwork(&obj)
	if err != nil {
		log.Errorf("CreateNetwork error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return obj, nil
}

// DELETE rest call
func httpDeleteNetwork(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpDeleteNetwork: %+v", vars)

	key := vars["key"]

	// Delete the object
	err := DeleteNetwork(key)
	if err != nil {
		log.Errorf("DeleteNetwork error for: %s. Err: %v", key, err)
		return nil, err
	}

	// Return the obj
	return key, nil
}

// Create a network object
func CreateNetwork(obj *Network) error {
	// Validate parameters
	err := ValidateNetwork(obj)
	if err != nil {
		log.Errorf("ValidateNetwork retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// Check if we handle this object
	if objCallbackHandler.NetworkCb == nil {
		log.Errorf("No callback registered for network object")
		return errors.New("Invalid object type")
	}

	saveObj := obj

	// Check if object already exists
	if collections.networks[obj.Key] != nil {
		// Perform Update callback
		err = objCallbackHandler.NetworkCb.NetworkUpdate(collections.networks[obj.Key], obj)
		if err != nil {
			log.Errorf("NetworkUpdate retruned error for: %+v. Err: %v", obj, err)
			return err
		}

		// save the original object after update
		saveObj = collections.networks[obj.Key]
	} else {
		// save it in cache
		collections.networks[obj.Key] = obj

		// Perform Create callback
		err = objCallbackHandler.NetworkCb.NetworkCreate(obj)
		if err != nil {
			log.Errorf("NetworkCreate retruned error for: %+v. Err: %v", obj, err)
			delete(collections.networks, obj.Key)
			return err
		}
	}

	// Write it to modeldb
	err = saveObj.Write()
	if err != nil {
		log.Errorf("Error saving network %s to db. Err: %v", saveObj.Key, err)
		return err
	}

	return nil
}

// Return a pointer to network from collection
func FindNetwork(key string) *Network {
	obj := collections.networks[key]
	if obj == nil {
		return nil
	}

	return obj
}

// Delete a network object
func DeleteNetwork(key string) error {
	obj := collections.networks[key]
	if obj == nil {
		log.Errorf("network %s not found", key)
		return errors.New("network not found")
	}

	// Check if we handle this object
	if objCallbackHandler.NetworkCb == nil {
		log.Errorf("No callback registered for network object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.NetworkCb.NetworkDelete(obj)
	if err != nil {
		log.Errorf("NetworkDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// delete it from modeldb
	err = obj.Delete()
	if err != nil {
		log.Errorf("Error deleting network %s. Err: %v", obj.Key, err)
	}

	// delete it from cache
	delete(collections.networks, key)

	return nil
}

func (self *Network) GetType() string {
	return "network"
}

func (self *Network) GetKey() string {
	return self.Key
}

func (self *Network) Read() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to read network object")
		return errors.New("Empty key")
	}

	return modeldb.ReadObj("network", self.Key, self)
}

func (self *Network) Write() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Write network object")
		return errors.New("Empty key")
	}

	return modeldb.WriteObj("network", self.Key, self)
}

func (self *Network) Delete() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Delete network object")
		return errors.New("Empty key")
	}

	return modeldb.DeleteObj("network", self.Key)
}

func restoreNetwork() error {
	strList, err := modeldb.ReadAllObj("network")
	if err != nil {
		log.Errorf("Error reading network list. Err: %v", err)
	}

	for _, objStr := range strList {
		// Parse the json model
		var network Network
		err = json.Unmarshal([]byte(objStr), &network)
		if err != nil {
			log.Errorf("Error parsing object %s, Err %v", objStr, err)
			return err
		}

		// add it to the collection
		collections.networks[network.Key] = &network
	}

	return nil
}

// Validate a network object
func ValidateNetwork(obj *Network) error {
	// Validate key is correct
	keyStr := obj.TenantName + ":" + obj.NetworkName
	if obj.Key != keyStr {
		log.Errorf("Expecting Network Key: %s. Got: %s", keyStr, obj.Key)
		return errors.New("Invalid Key")
	}

	// Validate each field

	if len(obj.Encap) > 32 {
		return errors.New("encap string too long")
	}

	if obj.IsPrivate == false {
		obj.IsPrivate = true
	}

	if obj.IsPublic == false {
		obj.IsPublic = false
	}

	if obj.PktTag == 0 {
		obj.PktTag = 1
	}

	if obj.PktTag < 1 {
		return errors.New("pktTag Value Out of bound")
	}

	if obj.PktTag > 4094 {
		return errors.New("pktTag Value Out of bound")
	}

	subnetMatch := regexp.MustCompile("^([0-9]{1,3}?.[0-9]{1,3}?.[0-9]{1,3}?.[0-9]{1,3}?/[0-9]{1,2}?)$")
	if subnetMatch.MatchString(obj.Subnet) == false {
		return errors.New("subnet string invalid format")
	}

	return nil
}

// GET Oper REST call
func httpInspectNetTwo(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj NetTwoInspect
	log.Debugf("Received httpInspectNetTwo: %+v", vars)

	key := vars["key"]

	objConfig := collections.netTwos[key]
	if objConfig == nil {
		log.Errorf("netTwo %s not found", key)
		return nil, errors.New("netTwo not found")
	}
	obj.Config = *objConfig

	if err := GetOperNetTwo(&obj); err != nil {
		log.Errorf("GetNetTwo error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return &obj, nil
}

// Get a netTwoOper object
func GetOperNetTwo(obj *NetTwoInspect) error {
	// Check if we handle this object
	if objCallbackHandler.NetTwoCb == nil {
		log.Errorf("No callback registered for netTwo object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.NetTwoCb.NetTwoGetOper(obj)
	if err != nil {
		log.Errorf("NetTwoDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	return nil
}

// LIST REST call
func httpListNetTwos(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpListNetTwos: %+v", vars)

	list := make([]*NetTwo, 0)
	for _, obj := range collections.netTwos {
		list = append(list, obj)
	}

	// Return the list
	return list, nil
}

// GET REST call
func httpGetNetTwo(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetNetTwo: %+v", vars)

	key := vars["key"]

	obj := collections.netTwos[key]
	if obj == nil {
		log.Errorf("netTwo %s not found", key)
		return nil, errors.New("netTwo not found")
	}

	// Return the obj
	return obj, nil
}

// CREATE REST call
func httpCreateNetTwo(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGetNetTwo: %+v", vars)

	var obj NetTwo
	key := vars["key"]

	// Get object from the request
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Errorf("Error decoding netTwo create request. Err %v", err)
		return nil, err
	}

	// set the key
	obj.Key = key

	// Create the object
	err = CreateNetTwo(&obj)
	if err != nil {
		log.Errorf("CreateNetTwo error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return obj, nil
}

// DELETE rest call
func httpDeleteNetTwo(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpDeleteNetTwo: %+v", vars)

	key := vars["key"]

	// Delete the object
	err := DeleteNetTwo(key)
	if err != nil {
		log.Errorf("DeleteNetTwo error for: %s. Err: %v", key, err)
		return nil, err
	}

	// Return the obj
	return key, nil
}

// Create a netTwo object
func CreateNetTwo(obj *NetTwo) error {
	// Validate parameters
	err := ValidateNetTwo(obj)
	if err != nil {
		log.Errorf("ValidateNetTwo retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// Check if we handle this object
	if objCallbackHandler.NetTwoCb == nil {
		log.Errorf("No callback registered for netTwo object")
		return errors.New("Invalid object type")
	}

	saveObj := obj

	// Check if object already exists
	if collections.netTwos[obj.Key] != nil {
		// Perform Update callback
		err = objCallbackHandler.NetTwoCb.NetTwoUpdate(collections.netTwos[obj.Key], obj)
		if err != nil {
			log.Errorf("NetTwoUpdate retruned error for: %+v. Err: %v", obj, err)
			return err
		}

		// save the original object after update
		saveObj = collections.netTwos[obj.Key]
	} else {
		// save it in cache
		collections.netTwos[obj.Key] = obj

		// Perform Create callback
		err = objCallbackHandler.NetTwoCb.NetTwoCreate(obj)
		if err != nil {
			log.Errorf("NetTwoCreate retruned error for: %+v. Err: %v", obj, err)
			delete(collections.netTwos, obj.Key)
			return err
		}
	}

	// Write it to modeldb
	err = saveObj.Write()
	if err != nil {
		log.Errorf("Error saving netTwo %s to db. Err: %v", saveObj.Key, err)
		return err
	}

	return nil
}

// Return a pointer to netTwo from collection
func FindNetTwo(key string) *NetTwo {
	obj := collections.netTwos[key]
	if obj == nil {
		return nil
	}

	return obj
}

// Delete a netTwo object
func DeleteNetTwo(key string) error {
	obj := collections.netTwos[key]
	if obj == nil {
		log.Errorf("netTwo %s not found", key)
		return errors.New("netTwo not found")
	}

	// Check if we handle this object
	if objCallbackHandler.NetTwoCb == nil {
		log.Errorf("No callback registered for netTwo object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.NetTwoCb.NetTwoDelete(obj)
	if err != nil {
		log.Errorf("NetTwoDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// delete it from modeldb
	err = obj.Delete()
	if err != nil {
		log.Errorf("Error deleting netTwo %s. Err: %v", obj.Key, err)
	}

	// delete it from cache
	delete(collections.netTwos, key)

	return nil
}

func (self *NetTwo) GetType() string {
	return "netTwo"
}

func (self *NetTwo) GetKey() string {
	return self.Key
}

func (self *NetTwo) Read() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to read netTwo object")
		return errors.New("Empty key")
	}

	return modeldb.ReadObj("netTwo", self.Key, self)
}

func (self *NetTwo) Write() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Write netTwo object")
		return errors.New("Empty key")
	}

	return modeldb.WriteObj("netTwo", self.Key, self)
}

func (self *NetTwo) Delete() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Delete netTwo object")
		return errors.New("Empty key")
	}

	return modeldb.DeleteObj("netTwo", self.Key)
}

func restoreNetTwo() error {
	strList, err := modeldb.ReadAllObj("netTwo")
	if err != nil {
		log.Errorf("Error reading netTwo list. Err: %v", err)
	}

	for _, objStr := range strList {
		// Parse the json model
		var netTwo NetTwo
		err = json.Unmarshal([]byte(objStr), &netTwo)
		if err != nil {
			log.Errorf("Error parsing object %s, Err %v", objStr, err)
			return err
		}

		// add it to the collection
		collections.netTwos[netTwo.Key] = &netTwo
	}

	return nil
}

// Validate a netTwo object
func ValidateNetTwo(obj *NetTwo) error {
	// Validate key is correct
	keyStr := obj.TenantName + ":" + obj.NetworkName
	if obj.Key != keyStr {
		log.Errorf("Expecting NetTwo Key: %s. Got: %s", keyStr, obj.Key)
		return errors.New("Invalid Key")
	}

	// Validate each field

	if len(obj.Encap) > 32 {
		return errors.New("encap string too long")
	}

	if obj.IsPrivate == false {
		obj.IsPrivate = true
	}

	if obj.IsPublic == false {
		obj.IsPublic = false
	}

	if obj.PktTag == 0 {
		obj.PktTag = 1
	}

	if obj.PktTag < 1 {
		return errors.New("pktTag Value Out of bound")
	}

	if obj.PktTag > 4094 {
		return errors.New("pktTag Value Out of bound")
	}

	subnetMatch := regexp.MustCompile("^([0-9]{1,3}?.[0-9]{1,3}?.[0-9]{1,3}?.[0-9]{1,3}?/[0-9]{1,2}?)$")
	if subnetMatch.MatchString(obj.Subnet) == false {
		return errors.New("subnet string invalid format")
	}

	return nil
}

// GET Oper REST call
func httpInspectEndpoint(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj EndpointInspect
	log.Debugf("Received httpInspectEndpoint: %+v", vars)

	obj.Oper.Key = vars["key"]

	if err := GetOperEndpoint(&obj); err != nil {
		log.Errorf("GetEndpoint error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return &obj, nil
}

// Get a endpointOper object
func GetOperEndpoint(obj *EndpointInspect) error {
	// Check if we handle this object
	if objCallbackHandler.EndpointCb == nil {
		log.Errorf("No callback registered for endpoint object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.EndpointCb.EndpointGetOper(obj)
	if err != nil {
		log.Errorf("EndpointDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	return nil
}

// GET Oper REST call
func httpInspectEpList(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj EpListInspect
	log.Debugf("Received httpInspectEpList: %+v", vars)

	obj.Oper.Key = vars["key"]

	if err := GetOperEpList(&obj); err != nil {
		log.Errorf("GetEpList error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return &obj, nil
}

// Get a epListOper object
func GetOperEpList(obj *EpListInspect) error {
	// Check if we handle this object
	if objCallbackHandler.EpListCb == nil {
		log.Errorf("No callback registered for epList object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.EpListCb.EpListGetOper(obj)
	if err != nil {
		log.Errorf("EpListDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	return nil
}

// GET Oper REST call
func httpInspectNetWithEp(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var obj NetWithEpInspect
	log.Debugf("Received httpInspectNetWithEp: %+v", vars)

	obj.Oper.Key = vars["key"]

	if err := GetOperNetWithEp(&obj); err != nil {
		log.Errorf("GetNetWithEp error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return &obj, nil
}

// Get a netWithEpOper object
func GetOperNetWithEp(obj *NetWithEpInspect) error {
	// Check if we handle this object
	if objCallbackHandler.NetWithEpCb == nil {
		log.Errorf("No callback registered for netWithEp object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.NetWithEpCb.NetWithEpGetOper(obj)
	if err != nil {
		log.Errorf("NetWithEpDelete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	return nil
}
