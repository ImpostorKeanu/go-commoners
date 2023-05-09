package server

import (
    rando "github.com/arch4ngel/go-commoners/rando"
    "github.com/gin-gonic/gin"
    "net/http"
)

// RandoServer represents a Gin HTTP server used to offer
// various random value generation across the network.
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
    gV0.GET("/ascii", randAsciiValsHandler)
    gV0.GET("/ascii/string", randAsciiStringHandler)
    http.ListenAndServe(string(rS), r)
}

func randAsciiStringHandler(c *gin.Context) {
    q := AsciiQueryBind{}
    if err := c.ShouldBindQuery(&q); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
            SuccessRespField: SuccessRespField{false},
            Message:          err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, StandardResp{
        SuccessRespField: SuccessRespField{true},
        Value:            rando.AnyAsciiString(q.MinLen, q.AllowNumbers, q.Delimiter),
    })
}

func randAsciiValsHandler(c *gin.Context) {
    q := AsciiQueryBind{}
    if err := c.ShouldBindQuery(&q); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
            SuccessRespField: SuccessRespField{false},
            Message:          err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, RandValsResp{
        SuccessRespField: SuccessRespField{true},
        Value:            rando.AnyAsciiSlice(q.MinLen, q.AllowNumbers),
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
    c.JSON(http.StatusOK, RandValsResp{
        SuccessRespField: SuccessRespField{true},
        Value:            rando.AnySlice(q.MinLen),
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
    v := rando.AnyString(q.MinLen, q.Delimiter)
    if q.RandFix {
        v, _ = rando.AnyAsciiRandfix(v, q.Delimiter, 15)
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
        Value:            rando.Uuid(),
    })
}

// hostnameHandler returns a dynamically generated hostname.
//
// See dnsHandler to ensure that a hostname is unique to a given
// apex domain.
//
// See NounTypeQueryBind for expected query parameters.
func hostnameHandler(c *gin.Context) {
    q := NounTypeQueryBind{}
    if err := c.ShouldBindQuery(&q); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, ErrResp{
            SuccessRespField: SuccessRespField{false},
            Message:          err.Error(),
        })
        return
    }
    hn, _ := rando.Hostname(q.NounType)
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

    fqdn, resolved, ips, err := rando.DnsUntil(q.NounType, q.ApexDomain, q.UniqueRequired)
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
