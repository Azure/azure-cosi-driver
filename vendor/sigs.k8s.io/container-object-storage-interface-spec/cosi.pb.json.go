// Code generated by protoc-gen-go-json. DO NOT EDIT.
// source: sigs.k8s.io/container-object-storage-interface-spec/cosi.proto

package cosi

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalJSON implements json.Marshaler
func (msg *S3) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *S3) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *AzureBlob) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *AzureBlob) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *GCS) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *GCS) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *Protocol) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Protocol) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *CredentialDetails) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *CredentialDetails) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverGetInfoRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverGetInfoRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverGetInfoResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverGetInfoResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverCreateBucketRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverCreateBucketRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverCreateBucketResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverCreateBucketResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverDeleteBucketRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverDeleteBucketRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverDeleteBucketResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverDeleteBucketResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverGrantBucketAccessRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverGrantBucketAccessRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverGrantBucketAccessResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverGrantBucketAccessResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverRevokeBucketAccessRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverRevokeBucketAccessRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *DriverRevokeBucketAccessResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *DriverRevokeBucketAccessResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}
