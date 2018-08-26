package router

import (
	"encoding/base64"
)

/*
 * Query Response
 */

type queryResponse struct {
	IsSucc bool   `json:"is_success"`
	IsAuth bool   `json:"is_authorized"`
	IsNFS  bool   `json:"is_nfs"`
	GzUsed int64  `json:"gaz_used"`
	Dat    []byte `json:"data"`
	Lg     string `json:"log"`
}

func createQueryResponse(isSuccess bool, isAuthorized bool, hasNFS bool, gazUsed int64, log string, data []byte) (QueryResponse, error) {
	out := queryResponse{
		IsSucc: isSuccess,
		IsAuth: isAuthorized,
		IsNFS:  hasNFS,
		GzUsed: gazUsed,
		Lg:     log,
		Dat:    data,
	}

	return &out, nil
}

// IsSuccess returns true if the query was successful, false otherwise
func (obj *queryResponse) IsSuccess() bool {
	return obj.IsSucc
}

// IsAuthorized returns true if the request was authorized, false otherwise
func (obj *queryResponse) IsAuthorized() bool {
	return obj.IsAuth
}

// HasInsufficientFunds returns true if the user had insufficient funds to conver the query costs, false otherwise
func (obj *queryResponse) HasInsufficientFunds() bool {
	return obj.IsNFS
}

// GazUsed returns the amount of gaz used to execute the query
func (obj *queryResponse) GazUsed() int64 {
	return obj.GzUsed
}

// Log returns log
func (obj *queryResponse) Log() string {
	return obj.Lg
}

// Data returns the data
func (obj *queryResponse) Data() []byte {
	return obj.Dat
}

// UnMarshal unmarshals the data to the pointer
func (obj *queryResponse) UnMarshal(v interface{}) error {
	jsErr := cdc.UnmarshalJSON(obj.Dat, v)
	if jsErr != nil {
		return jsErr
	}

	return nil
}

/*
 * Trx Chk Response
 */

type trxChkResponse struct {
	CnBeExecuted   bool   `json:"can_be_executed"`
	CnBeAuthorized bool   `json:"can_be_authorized"`
	GzWanted       int64  `json:"gaz_wanted"`
	Lg             string `json:"log"`
}

func createTrxChkResponse(canBeExecuted bool, canBeAuthorized bool, gazWanted int64, log string) (TrxChkResponse, error) {
	out := trxChkResponse{
		CnBeExecuted:   canBeExecuted,
		CnBeAuthorized: canBeAuthorized,
		GzWanted:       gazWanted,
		Lg:             log,
	}

	return &out, nil
}

// CanBeExecuted returns true if the transaction can be executed, false otherwise
func (obj *trxChkResponse) CanBeExecuted() bool {
	return obj.CnBeExecuted
}

// CanBeAuthorized returns true if the transaction can be executed, false otherwise
func (obj *trxChkResponse) CanBeAuthorized() bool {
	return obj.CnBeAuthorized
}

// GazWanted returns the amount of gaz wanted to execute the transaction
func (obj *trxChkResponse) GazWanted() int64 {
	return obj.GzWanted
}

// Log returns the log
func (obj *trxChkResponse) Log() string {
	return obj.Lg
}

/*
 * Trx Response
 */

type jsonTrxResponse struct {
	IsSuc  bool              `json:"is_success"`
	IsAuth bool              `json:"is_authorized"`
	IsNFS  bool              `json:"is_nfs"`
	Tgs    map[string]string `json:"tags"`
	GzUsed int64             `json:"gaz_used"`
	Lg     string            `json:"log"`
}

func createJSONTrxResponse(res TrxResponse) *jsonTrxResponse {

	input := res.Tags()
	tgs := map[string]string{}
	for keyname, oneTg := range input {
		encodedStr := base64.StdEncoding.EncodeToString(oneTg)
		tgs[keyname] = encodedStr
	}

	out := jsonTrxResponse{
		IsSuc:  res.IsSuccess(),
		IsAuth: res.IsAuthorized(),
		IsNFS:  res.HasInsufficientFunds(),
		Tgs:    tgs,
		GzUsed: res.GazUsed(),
		Lg:     res.Log(),
	}

	return &out
}

type trxResponse struct {
	isSuc  bool
	isAuth bool
	isNFS  bool
	tgs    map[string][]byte
	gzUsed int64
	lg     string
}

func createTrxResponse(isSuccess bool, isAuthorized bool, isNFS bool, tags map[string][]byte, gazUsed int64, log string) (TrxResponse, error) {
	out := trxResponse{
		isSuc:  isSuccess,
		isAuth: isAuthorized,
		isNFS:  isNFS,
		tgs:    tags,
		gzUsed: gazUsed,
		lg:     log,
	}

	return &out, nil
}

// IsSuccess returns true if the transaction is successful, false otherwise
func (obj *trxResponse) IsSuccess() bool {
	return obj.isSuc
}

// IsAuthorized returns true if the transaction is authorized, false otherwise
func (obj *trxResponse) IsAuthorized() bool {
	return obj.isAuth
}

// HasInsufficientFunds returns true if the user had insufficient funds, false otherwise
func (obj *trxResponse) HasInsufficientFunds() bool {
	return obj.isNFS
}

// Tags returns the tags
func (obj *trxResponse) Tags() map[string][]byte {
	return obj.tgs
}

// GazUsed returns the amount of gaz used
func (obj *trxResponse) GazUsed() int64 {
	return obj.gzUsed
}

// Log returns the logs
func (obj *trxResponse) Log() string {
	return obj.lg
}

// MarshalJSON converts the instance to JSON
func (obj *trxResponse) MarshalJSON() ([]byte, error) {
	jsState := createJSONTrxResponse(obj)
	js, jsErr := cdc.MarshalJSON(jsState)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *trxResponse) UnmarshalJSON(data []byte) error {
	jsTrxResponse := new(jsonTrxResponse)
	jsErr := cdc.UnmarshalJSON(data, jsTrxResponse)
	if jsErr != nil {
		return jsErr
	}

	input := jsTrxResponse.Tgs
	tgs := map[string][]byte{}
	for keyname, encodedStr := range input {
		decodedStr, decodedErr := base64.StdEncoding.DecodeString(encodedStr)
		if decodedErr != nil {
			return decodedErr
		}

		tgs[keyname] = decodedStr
	}

	obj.isSuc = jsTrxResponse.IsSuc
	obj.isAuth = jsTrxResponse.IsAuth
	obj.isNFS = jsTrxResponse.IsNFS
	obj.tgs = tgs
	obj.gzUsed = jsTrxResponse.GzUsed
	obj.lg = jsTrxResponse.Lg
	return nil
}
