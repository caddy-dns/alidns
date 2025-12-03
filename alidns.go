package alidns

import (
	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/alidns"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *alidns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.alidns",
		New: func() caddy.Module { return &Provider{new(alidns.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AccessKeyID = repl.ReplaceAll(p.Provider.AccessKeyID, "")
	p.Provider.AccessKeySecret = repl.ReplaceAll(p.Provider.AccessKeySecret, "")
	p.Provider.SecurityToken = repl.ReplaceAll(p.SecurityToken, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	alidns {
//	    access_key_id "<access_key_id>"
//	    access_key_secret "<access_key_secret>"
//		region_id "<region_id,defaults to 'cn-hangzhou' if empty>"
//	    security_token "<security_token, if you use the STS authorization it's required by aliyun>"
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "access_key_id":
				if d.NextArg() {
					p.Provider.AccessKeyID = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "access_key_secret":
				if d.NextArg() {
					p.Provider.AccessKeySecret = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "region_id":
				if d.NextArg() {
					p.Provider.RegionID = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "security_token":
				if d.NextArg() {
					p.Provider.SecurityToken = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.AccessKeyID == "" || p.AccessKeySecret == "" {
		return d.Err("AccessKeyID or AccessKeySecret is empty")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
