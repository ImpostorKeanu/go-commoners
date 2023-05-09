package server

import (
    "github.com/arch4ngel/go-commoners/rando/nouns"
    "net"
)

// AllowNumbersQueryBind provides an AllowNumbers field.
type AllowNumbersQueryBind struct {
    AllowNumbers bool `form:"allow-numbers"`
}

// DelimiterQueryBind provides a Delimiter field.
type DelimiterQueryBind struct {
    Delimiter string `form:"delimiter"`
}

// MinLenQueryBind provides a MinLen field.
type MinLenQueryBind struct {
    MinLen uint32 `form:"min-len" binding:"gt=0"`
}

// NounTypeQueryBind provides a binding for noun type query values..
type NounTypeQueryBind struct {
    NounType nouns.NounType `form:"noun-type" binding:"required,oneof=common proper"`
}

// AsciiQueryBind provides a binding for ASCII value requests.
type AsciiQueryBind struct {
    MinLenQueryBind
    AllowNumbersQueryBind
    DelimiterQueryBind
}

// PassphraseQueryBind provides a binding for passphrase value requests.
type PassphraseQueryBind struct {
    MinLenQueryBind
    DelimiterQueryBind
}

// DnsQueryBind provides a binding for DNS value requests.
type DnsQueryBind struct {
    NounTypeQueryBind
    // ApexDomain for the value, e.g., google.com
    ApexDomain string `form:"apex-domain" binding:"hostname,required"`
    // UniqueRequired determines if DNS queries should be performed to
    // determine if the newly generated value is unique. Generation will
    // continue until a unique value is produced.
    UniqueRequired bool `form:"unique-required" binding:"boolean"`
}

// SuccessRespField is the success field for server
// responses.
type SuccessRespField struct {
    // Success determines if the call was successful.
    Success bool `json:"success" mapstructure:"success"`
}

// ErrResp is the standard server response structure.
type ErrResp struct {
    SuccessRespField `mapstructure:",squash" json:",squash"`
    // Message is any message accompanying the error.
    Message string `json:"message" mapstructure:"message"`
}

// StandardResp is the standard success structure for
// server responses.
type StandardResp struct {
    SuccessRespField `mapstructure:",squash" json:",squash"`
    // Value contains the string output for standard respones.
    Value string `json:"value" mapstructure:"value"`
}

// NounTypeRespField is the NounType field for server
// responses.
type NounTypeRespField struct {
    NounType nouns.NounType `json:"noun_type" mapstructure:"noun_type"`
}

// DnsResolutionFields is a structure for the DnsServerResp.
type DnsResolutionFields struct {
    Resolved bool     `json:"resolved" mapstructure:"resolved"`
    IPs      []net.IP `json:"ips" mapstructure:"ips"`
}

// DnsServerResp is the structure returned when generating
// randomized DNS hostnames.
type DnsServerResp struct {
    StandardResp      `mapstructure:",squash" json:",squash"`
    NounTypeRespField `mapstructure:",squash" json:",squash"`
    Resolution        DnsResolutionFields `json:"resolution" mapstructure:"resolution"`
}

// HostnameServerResp is the structure returned when generating
// randomized hostnames.
type HostnameServerResp struct {
    StandardResp      `mapstructure:",squash" json:",squash"`
    NounTypeRespField `mapstructure:",squash" json:",squash"`
}

// RandValsResp is the structure returned when generating a list
// of randomized values.
type RandValsResp struct {
    SuccessRespField `mapstructure:",squash" json:",squash"`
    Value            []string `json:"value" mapstructure:"value"`
}
