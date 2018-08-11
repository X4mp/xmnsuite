package hashes

type hashes struct {
	data map[string]map[string][]byte
}

func createHashes() Hashes {
	out := hashes{
		data: map[string]map[string][]byte{},
	}

	return &out
}

/*
   Exists returns if field is an existing field in the hash stored at key.

   Returns:
       true if the hash contains field
       false if the hash does not contain field, or key does not exist.
*/
func (obj *hashes) Exists(key string, field string) bool {
	value := obj.Get(key, field)
	return value != nil
}

/*
   Get returns the value associated with field in the hash stored at key.

   Returns:
       the value associated with field, or nil when field is not present in the
       hash or key does not exist.
*/
func (obj *hashes) Get(key string, field string) []byte {
	ins := obj.GetAll(key)
	if value, ok := ins[field]; ok {
		return value
	}

	return nil
}

/*
   GetAll returns all fields and values of the hash stored at key. In the
   returned value, every field name is followed by its value, so the length of
   the reply is twice the size of the hash.

   Returns:
        map of fields and their values stored in the hash, or an empty map
        when key does not exist.
*/
func (obj *hashes) GetAll(key string) map[string][]byte {
	if ins, ok := obj.data[key]; ok {
		return ins
	}

	return map[string][]byte{}
}

/*
   MultiGet returns the values associated with the specified fields in the hash
   stored at key.

   For every field that does not exist in the hash, a nil value is returned.
   Because non-existing keys are treated as empty hashes,
   running MultiGet against a non-existing key will return a list of nil values.

   Returns:
        map of field -> values associated with the given fields
*/
func (obj *hashes) MultiGet(key string, fields ...string) map[string][]byte {
	out := map[string][]byte{}
	all := obj.GetAll(key)
	for _, oneField := range fields {
		if value, ok := all[oneField]; ok {
			out[oneField] = value
			continue
		}

		out[oneField] = nil
	}

	return out
}

/*
   Set sets field in the hash stored at key to value. If key does not exist, a
   new key holding a hash is created. If field already exists in the hash, it is
   overwritten.

   Returns:
        true if field is a new field in the hash and value was set.
        false if field already exists in the hash and the value was updated.
*/
func (obj *hashes) Set(key string, field string, value []byte) bool {
	return false
}

/*
   SetNX sets field in the hash stored at key to value, only if field does not
   yet exist. If key does not exist, a new key holding a hash is created. If
   field already exists, this operation has no effect.

   Returns:
        true if field is a new field in the hash and value was set.
        false if field already exists in the hash and no operation was performed.
*/
func (obj *hashes) SetNX(key string, field string, value []byte) bool {
	return false
}

/*
   MultiSet sets the specified fields to their respective values in the hash
   stored at key. This command overwrites any specified fields already existing
   in the hash. If key does not exist, a new key holding a hash is created.

   Returns:
        Nothing
*/
func (obj *hashes) MultiSet(key string, keyValues ...map[string][]byte) {

}

/*
   IncrBy increments the number stored at field in the hash stored at key by
   increment. If key does not exist, a new key holding a hash is created. If
   field does not exist the value is set to 0 before the operation is performed.

   Note: since the increment argument is signed, both increment and decrement
   operations can be performed.

   Returns:
        the value at field after the increment operation.
*/
func (obj *hashes) IncrBy(key string, field string, increment int64) int64 {
	return 0
}

/*
   IncrByFloat increments the specified field of a hash stored at key, and
   representing a floating point number, by the specified increment. If the
   increment value is negative, the result is to have the hash field value
   decremented instead of incremented. If the field does not exist, it is set
   to 0 before performing the operation. An error is returned if one of the
   following conditions occur:

        1. The field contains a value of the wrong type.
        2. The current field content or the specified increment are not
           parsable as a double precision floating point number.

   Returns:
        the value of field after the increment.
*/
func (obj *hashes) IncrByFloat(key string, field string, increment float64) (float64, error) {
	return float64(0), nil
}

/*
   Len returns the number of fields contained in the hash stored at key.

   Returns:
        the number of fields in the hash, or 0 when key does not exist.
*/
func (obj *hashes) Len(key string) int64 {
	return 0
}

/*
   StrLen returns the string length of the value associated with field in the
   hash stored at key. If the key or the field do not exist, 0 is returned.

   Returns:
         the string length of the value associated with field, or zero when
         field is not present in the hash or key does not exist at all.
*/
func (obj *hashes) StrLen(key string, field string) int {
	return 0
}

/*
   Del removes the specified fields from the hash stored at key. Specified
   fields that do not exist within this hash are ignored. If key does not
   exist, it is treated as an empty hash and this command returns 0.

   Returns:
         the number of fields that were removed from the hash, not including
         specified but non existing fields.
*/
func (obj *hashes) Del(key string, fields ...string) int {
	return 0
}

/*
   Keys returns all field names in the hash stored at key.

   Returns:
         the list of fields in the hash, or an empty list when key does not exist.
*/
func (obj *hashes) Keys(key string) []string {
	return nil
}

/*
   Vals returns all values in the hash stored at key.

   Returns:
         the list of values in the hash, or an empty list when key does not exist.
*/
func (obj *hashes) Vals(key string) []byte {
	return nil
}
