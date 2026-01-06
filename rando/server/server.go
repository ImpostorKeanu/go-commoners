package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	rando "github.com/impostorkeanu/go-commoners/rando"
)

// RandoServer is a Gin HTTP server used to offer various random value
// generation via REST API.
//
// It should be initialized to a string value formatted
// as a socket that the server will listen on.
type RandoServer string

func (rS RandoServer) Start() {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	gV0 := r.Group("/v0")
	gV0.GET("/uuid", uuidHandler)
	gV0.GET("/host", hostnameHandler)
	gV0.GET("/dns", dnsHandler)
	gV0.GET("/values", randValsHandler)
	gV0.GET("/passphrase", passphraseHandler)
	gV0.GET("/ascii", randASCIIValsHandler)
	gV0.GET("/ascii/string", randASCIIStringHandler)
	http.ListenAndServe(string(rS), r)
}

func randASCIIStringHandler(c *gin.Context) {
	q := AsciiQueryBind{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}

	v, _ := rando.UntilCleanString(q.CheckProfanity, func() (string, error) {
		return rando.AnyASCIIString(q.MinLen, q.AllowNumbers, q.Delimiter), nil
	})
	c.JSON(http.StatusOK, StandardResp{
		SuccessRespField: SuccessRespField{true},
		Value:            v,
	})
}

func randASCIIValsHandler(c *gin.Context) {
	q := AsciiQueryBind{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}
	v, _ := rando.UntilCleanSlice(q.CheckProfanity, func() ([]string, error) {
		return rando.AnyAsciiSlice(q.MinLen, q.AllowNumbers), nil
	})
	c.JSON(http.StatusOK, RandValsResp{
		SuccessRespField: SuccessRespField{true},
		Value:            v,
	})
}

// randValsHandler generates a list of random values up to a
// joined length.
func randValsHandler(c *gin.Context) {
	q := MinLenQueryBind{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}
	v, _ := rando.UntilCleanSlice(q.CheckProfanity, func() ([]string, error) {
		return rando.AnySlice(q.MinLen), nil
	})
	c.JSON(http.StatusOK, RandValsResp{
		SuccessRespField: SuccessRespField{true},
		Value:            v,
	})
}

// passphraseHandler generates a string of random values up to
// a joined length.
func passphraseHandler(c *gin.Context) {
	q := PassphraseQueryBind{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}
	v, err := rando.UntilCleanString(q.CheckProfanity, func() (string, error) {
		var err error
		v := rando.AnyString(q.MinLen, q.Delimiter)
		if q.RandFixLength > 0 {
			v, err = rando.AnyASCIIRandfix(v, q.Delimiter, q.RandFixLength)
		}
		return v, err
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{Success: false},
			Message:          fmt.Sprintf("error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, StandardResp{
		SuccessRespField: SuccessRespField{true},
		Value:            v,
	})
}

// uuidHandler returns a UUID value.
func uuidHandler(c *gin.Context) {
	c.JSON(http.StatusOK, StandardResp{
		SuccessRespField: SuccessRespField{true},
		Value:            rando.UUID(),
	})
}

// hostnameHandler returns a dynamically generated hostname.
//
// See dnsHandler to ensure that a hostname is unique to a given
// apex domain.
//
// See NounTypeQueryBind for expected query parameters.
func hostnameHandler(c *gin.Context) {
	q := struct {
		NounTypeQueryBind
		ProfanityCheckQueryBind
	}{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}

	hn, err := rando.UntilCleanString(q.CheckProfanity, func() (string, error) {
		return rando.Hostname(q.NounType)
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{Success: false},
			Message:          fmt.Sprintf("error: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, HostnameServerResp{
		StandardResp: StandardResp{
			SuccessRespField: SuccessRespField{true},
			Value:            hn,
		},
		NounTypeRespField: NounTypeRespField{q.NounType},
	})
}

// dnsHandler returns a dynamically generated hostname and
// performs DNS resolution to determine if an A record for
// it already exists.
func dnsHandler(c *gin.Context) {
	q := DnsQueryBind{}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}

	var resolved bool
	var ips []net.IP

	fqdn, err := rando.UntilCleanString(q.CheckProfanity, func() (string, error) {
		f, r, i, e := rando.DNSUntil(q.NounType, q.ApexDomain, q.UniqueRequired)
		resolved = r
		ips = i
		return f, e
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, ErrResp{
			SuccessRespField: SuccessRespField{false},
			Message:          err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DnsServerResp{
		StandardResp: StandardResp{
			SuccessRespField: SuccessRespField{Success: true},
			Value:            fqdn,
		},
		NounTypeRespField: NounTypeRespField{NounType: q.NounType},
		Resolution: DnsResolutionFields{
			Resolved: resolved,
			IPs:      ips,
		},
	})
}
