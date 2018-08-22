package tendermint

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	defaultValidatorPower = 10
)

func generateBlockchain(namespace string, name string, id *uuid.UUID) (Blockchain, error) {
	privKey := ed25519.GenPrivKey()
	return generateBlockchainWithPrivateKey(namespace, name, id, privKey)
}

func generateBlockchainWithPrivateKey(namespace string, name string, id *uuid.UUID, privKey crypto.PrivKey) (Blockchain, error) {

	//creation time:
	crOn := time.Now().UTC()

	//create the path:
	path := createPath(namespace, name, id)

	//create a validator:
	appHash := "" //tendermint needs an empty appHash on genesis
	validator := createValidator(appHash, defaultValidatorPower, privKey.PubKey())

	//create the genesis:
	gen, genErr := createGenesis([]byte(""), path, []Validator{
		validator,
	}, crOn)

	if genErr != nil {
		return nil, genErr
	}

	//create the private validator:
	pubKey := privKey.PubKey()
	addr := fmt.Sprintf("%X", pubKey.Address())
	pv := createPrivateValidator(addr, pubKey, privKey, 0, 0, 0)

	//create  the blockchain:
	blkChain := createBlockchain(gen, pv, privKey)
	return blkChain, nil

}

type jsonValidator struct {
	Name   string        `json:"name"`
	Power  int           `json:"power"`
	PubKey crypto.PubKey `json:"pub_key"`
}

func createJSONValidator(validator Validator) (*jsonValidator, error) {
	out := jsonValidator{
		Name:   validator.GetName(),
		Power:  validator.GetPower(),
		PubKey: validator.GetPubKey(),
	}

	return &out, nil
}

type validator struct {
	name   string
	power  int
	pubKey crypto.PubKey
}

func createValidator(name string, power int, pubKey crypto.PubKey) Validator {
	out := validator{
		name:   name,
		power:  power,
		pubKey: pubKey,
	}

	return &out
}

func createValidatorFromJSON(jsVal *jsonValidator) Validator {
	out := createValidator(jsVal.Name, jsVal.Power, jsVal.PubKey)
	return out
}

// GetName returns the name
func (obj *validator) GetName() string {
	return obj.name
}

// GetPower returns the power
func (obj *validator) GetPower() int {
	return obj.power
}

// GetPubKey returns the PublicKey
func (obj *validator) GetPubKey() crypto.PubKey {
	return obj.pubKey
}

// MarshalJSON converts the instance to JSON
func (obj *validator) MarshalJSON() ([]byte, error) {
	jsVal, jsValErr := createJSONValidator(obj)
	if jsValErr != nil {
		return nil, jsValErr
	}

	js, jsErr := cdc.MarshalJSON(jsVal)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *validator) UnmarshalJSON(data []byte) error {
	jsVal := new(jsonValidator)
	jsErr := cdc.UnmarshalJSON(data, jsVal)
	if jsErr != nil {
		return jsErr
	}

	obj.name = jsVal.Name
	obj.power = jsVal.Power
	obj.pubKey = jsVal.PubKey
	return nil
}

type path struct {
	namespace string
	name      string
	id        *uuid.UUID
}

func createPath(namespace string, name string, id *uuid.UUID) Path {
	out := path{
		namespace: namespace,
		name:      name,
		id:        id,
	}

	return &out
}

func createPathFromString(pathAsString string) (Path, error) {
	splits := strings.Split(pathAsString, string(filepath.Separator))
	if len(splits) != 3 {
		str := fmt.Sprintf("the path string (%s) is invalid", pathAsString)
		return nil, errors.New(str)
	}

	id, idErr := uuid.FromString(splits[2])
	if idErr != nil {
		return nil, idErr
	}

	out := createPath(splits[0], splits[1], &id)
	return out, nil
}

// GetNamespace returns the namespace
func (obj *path) GetNamespace() string {
	return obj.namespace
}

// GetName returns the name
func (obj *path) GetName() string {
	return obj.name
}

// GetID returns the ID
func (obj *path) GetID() *uuid.UUID {
	return obj.id
}

// String returns the string representation of the path
func (obj *path) String() string {
	return fmt.Sprintf("%s%s%s%s%s", obj.namespace, string(filepath.Separator), obj.name, string(filepath.Separator), obj.id.String())
}

type jsonGenesis struct {
	Head            string           `json:"app_hash"`
	ChainIdentifier string           `json:"chain_id"`
	Validators      []*jsonValidator `json:"validators"`
	CreatedOn       time.Time        `json:"genesis_time"`
}

func createJSONGenesis(gen Genesis) (*jsonGenesis, error) {
	validators := gen.GetValidators()
	jsValidators := []*jsonValidator{}
	for _, oneValidator := range validators {
		oneJSValidator, oneJSValidatorErr := createJSONValidator(oneValidator)
		if oneJSValidatorErr != nil {
			return nil, oneJSValidatorErr
		}

		jsValidators = append(jsValidators, oneJSValidator)
	}

	out := jsonGenesis{
		Head:            fmt.Sprintf("%X", gen.GetHead()),
		ChainIdentifier: strings.Replace(gen.GetPath().String(), string(filepath.Separator), "-", 2),
		Validators:      jsValidators,
		CreatedOn:       gen.CreatedOn(),
	}

	return &out, nil
}

type genesis struct {
	head       []byte
	validators []Validator
	crOn       time.Time
	path       Path
}

func createGenesis(head []byte, path Path, validators []Validator, createdOn time.Time) (Genesis, error) {
	blocks := [][]byte{
		[]byte(path.String()),
	}

	for _, oneValidator := range validators {
		blocks = append(blocks, []byte(oneValidator.GetName()))
	}

	out := genesis{
		head:       head,
		validators: validators,
		crOn:       createdOn,
		path:       path,
	}

	return &out, nil
}

// GetHead returns the head
func (obj *genesis) GetHead() []byte {
	return obj.head
}

// GetPath returns the path
func (obj *genesis) GetPath() Path {
	return obj.path
}

// GetValidators returns the validators
func (obj *genesis) GetValidators() []Validator {
	return obj.validators
}

// CreatedOn returns the creation time
func (obj *genesis) CreatedOn() time.Time {
	return obj.crOn
}

// MarshalJSON converts the instance to JSON
func (obj *genesis) MarshalJSON() ([]byte, error) {
	jsGen, jsGenErr := createJSONGenesis(obj)
	if jsGenErr != nil {
		return nil, jsGenErr
	}

	js, jsErr := cdc.MarshalJSON(jsGen)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *genesis) UnmarshalJSON(data []byte) error {
	jsGenesis := new(jsonGenesis)
	jsErr := cdc.UnmarshalJSON(data, jsGenesis)
	if jsErr != nil {
		return jsErr
	}

	validators := []Validator{}
	for _, oneJSValidator := range jsGenesis.Validators {
		oneValidator := createValidatorFromJSON(oneJSValidator)
		validators = append(validators, oneValidator)
	}

	pathAsString := strings.Replace(jsGenesis.ChainIdentifier, "-", string(filepath.Separator), 2)
	path, pathErr := createPathFromString(pathAsString)
	if pathErr != nil {
		return pathErr
	}

	decodedStr, decodedStrErr := hex.DecodeString(jsGenesis.Head)
	if decodedStrErr != nil {
		return decodedStrErr
	}

	obj.head = decodedStr
	obj.path = path
	obj.validators = validators
	obj.crOn = jsGenesis.CreatedOn
	return nil
}

type jsonPrivateValidator struct {
	Address    string         `json:"address"`
	PubKey     crypto.PubKey  `json:"pub_key"`
	PrivKey    crypto.PrivKey `json:"priv_key"`
	LastHeight int64          `json:"last_height"`
	LastRound  int64          `json:"last_round"`
	LastStep   int8           `json:"last_step"`
}

func createJSONPrivateValidator(val PrivateValidator) *jsonPrivateValidator {
	out := jsonPrivateValidator{
		Address:    val.GetAddress(),
		PubKey:     val.GetPubKey(),
		PrivKey:    val.GetPrivKey(),
		LastHeight: val.GetLastHeight(),
		LastRound:  val.GetLastRound(),
		LastStep:   val.GetLastStep(),
	}

	return &out
}

type privateValidator struct {
	addr       string
	pubKey     crypto.PubKey
	privKey    crypto.PrivKey
	lastHeight int64
	lastRound  int64
	lastStep   int8
}

func createPrivateValidator(addr string, pubKey crypto.PubKey, privKey crypto.PrivKey, lastHeight int64, lastRound int64, lastStep int8) PrivateValidator {
	out := privateValidator{
		addr:       addr,
		pubKey:     pubKey,
		privKey:    privKey,
		lastHeight: lastHeight,
		lastRound:  lastRound,
		lastStep:   lastStep,
	}

	return &out
}

// GetAddress returns the address
func (obj *privateValidator) GetAddress() string {
	return obj.addr
}

// GetPubKey returns the public key
func (obj *privateValidator) GetPubKey() crypto.PubKey {
	return obj.pubKey
}

// GetPrivKey returns the private key
func (obj *privateValidator) GetPrivKey() crypto.PrivKey {
	return obj.privKey
}

// GetLastHeight returns the blockchain last height
func (obj *privateValidator) GetLastHeight() int64 {
	return obj.lastHeight
}

// GetLastRound returns the blockchain last round
func (obj *privateValidator) GetLastRound() int64 {
	return obj.lastRound
}

// GetLastStep returns the blockchain last step
func (obj *privateValidator) GetLastStep() int8 {
	return obj.lastStep
}

// MarshalJSON converts the instance to JSON
func (obj *privateValidator) MarshalJSON() ([]byte, error) {
	jsVal := createJSONPrivateValidator(obj)
	js, jsErr := cdc.MarshalJSON(jsVal)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *privateValidator) UnmarshalJSON(data []byte) error {
	jsVal := new(jsonPrivateValidator)
	jsErr := cdc.UnmarshalJSON(data, jsVal)
	if jsErr != nil {
		return jsErr
	}

	obj.addr = jsVal.Address
	obj.pubKey = jsVal.PubKey
	obj.privKey = jsVal.PrivKey
	obj.lastHeight = jsVal.LastHeight
	obj.lastRound = jsVal.LastRound
	obj.lastStep = jsVal.LastStep
	return nil
}

type blockchain struct {
	gen Genesis
	pk  crypto.PrivKey
	pv  PrivateValidator
}

func createBlockchain(gen Genesis, pv PrivateValidator, pk crypto.PrivKey) Blockchain {
	out := blockchain{
		gen: gen,
		pk:  pk,
		pv:  pv,
	}

	return &out
}

// GetGenesis returns the genesis
func (obj *blockchain) GetGenesis() Genesis {
	return obj.gen
}

// GetPK returns the private key
func (obj *blockchain) GetPK() crypto.PrivKey {
	return obj.pk
}

// GetPV returns the private validator
func (obj *blockchain) GetPV() PrivateValidator {
	return obj.pv
}
