package entity

type representation struct {
	met      MetaData
	keynames Keynames
	toData   ToStorable
	sync     Sync
}

func createRepresentation(met MetaData, toData ToStorable) Representation {
	return createRepresentationWithKeynamesAndSync(met, toData, nil, nil)
}

func createRepresentationWithKeynames(met MetaData, toData ToStorable, keynames Keynames) Representation {
	return createRepresentationWithKeynamesAndSync(met, toData, keynames, nil)
}

func createRepresentationWithSync(met MetaData, toData ToStorable, sync Sync) Representation {
	return createRepresentationWithKeynamesAndSync(met, toData, nil, sync)
}

func createRepresentationWithKeynamesAndSync(met MetaData, toData ToStorable, keynames Keynames, sync Sync) Representation {
	out := representation{
		met:      met,
		keynames: keynames,
		toData:   toData,
		sync:     sync,
	}

	return &out
}

// MetaData returns the MetaData instance
func (obj *representation) MetaData() MetaData {
	return obj.met
}

// ToStorable returns the instance to data
func (obj *representation) ToStorable() ToStorable {
	return obj.toData
}

// HasKeynames returns true if there is keynames, false otherwise
func (obj *representation) HasKeynames() bool {
	return obj.keynames != nil
}

// Keynames returns the keynames if any
func (obj *representation) Keynames() Keynames {
	return obj.keynames
}

// HasSync returns true if there is a Sync func, false otherwise
func (obj *representation) HasSync() bool {
	return obj.sync != nil
}

// Sync returns the sync func
func (obj *representation) Sync() Sync {
	return obj.sync
}
