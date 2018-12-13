package currency

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type currency struct {
	UUID *uuid.UUID        `json:"id"`
	Cat  category.Category `json:"category"`
	Sym  string            `json:"symbol"`
	Nme  string            `json:"name"`
	Desc string            `json:"description"`
}

func createCurrency(id *uuid.UUID, cat category.Category, symbol string, name string, description string) (Currency, error) {

	if len(symbol) != amountOfCharactersForSymbol {
		str := fmt.Sprintf("the symbol (%s) needs %d characters", symbol, len(symbol))
		return nil, errors.New(str)
	}

	if len(name) > maxAountOfCharactersForName {
		str := fmt.Sprintf("the name (%s) contains %d characters, the limit is: %d", name, len(name), maxAountOfCharactersForName)
		return nil, errors.New(str)
	}

	if len(description) > maxAmountOfCharactersForDescription {
		str := fmt.Sprintf("the description (%s) contains %d characters, thelimit is: %d", description, len(description), maxAmountOfCharactersForDescription)
		return nil, errors.New(str)
	}

	out := currency{
		UUID: id,
		Cat:  cat,
		Sym:  symbol,
		Nme:  name,
		Desc: description,
	}

	return &out, nil
}

func fromNormalizedToCurrency(normalized *normalizedCurrency) (Currency, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	catMetaData := category.SDKFunc.CreateMetaData()
	catIns, catInsErr := catMetaData.Denormalize()(normalized.Category)
	if catInsErr != nil {
		return nil, catInsErr
	}

	if cat, ok := catIns.(category.Category); ok {
		return createCurrency(&id, cat, normalized.Symbol, normalized.Name, normalized.Description)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", catIns.ID().String())
	return nil, errors.New(str)
}

func fromStorableToCurrency(storable *storableCurrency, rep entity.Repository) (Currency, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	categoryID, categoryIDErr := uuid.FromString(storable.CategoryID)
	if categoryIDErr != nil {
		return nil, categoryIDErr
	}

	catMetaData := category.SDKFunc.CreateMetaData()
	catIns, catInsErr := rep.RetrieveByID(catMetaData, &categoryID)
	if catInsErr != nil {
		return nil, catInsErr
	}

	if cat, ok := catIns.(category.Category); ok {
		return createCurrency(&id, cat, storable.Symbol, storable.Name, storable.Description)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", catIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *currency) ID() *uuid.UUID {
	return obj.UUID
}

// Category returns the category
func (obj *currency) Category() category.Category {
	return obj.Cat
}

// Symbol returns the symbol
func (obj *currency) Symbol() string {
	return obj.Sym
}

// Name returns the name
func (obj *currency) Name() string {
	return obj.Nme
}

// Description returns the description
func (obj *currency) Description() string {
	return obj.Desc
}
