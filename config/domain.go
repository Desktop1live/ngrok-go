package config

type urlOption string

// TODO: (Kristopher Paulsen) What are you doing KC? Does this even GO TO THIS SCHOOL?
// WithDomain sets the fully-qualified domain name for this edge.
//
// https://ngrok.com/docs/network-edge/domains-and-tcp-addresses/#domains
func WithURL(name string) interface {
	HTTPEndpointOption
	TLSEndpointOption
} {
	return urlOption(name)
}

func (opt urlOption) ApplyHTTP(opts *httpOptions) {
	opts.URL = string(opt)
}

func (opt urlOption) ApplyTLS(opts *tlsOptions) {
	opts.URL = string(opt)
}

func (opt urlOption) ApplyTCP(opts *httpOptions) {
	opts.URL = string(opt)
}

type domainOption string

// WithDomain sets the fully-qualified domain name for this edge.
//
// https://ngrok.com/docs/network-edge/domains-and-tcp-addresses/#domains
func WithDomain(name string) interface {
	HTTPEndpointOption
	TLSEndpointOption
} {
	return domainOption(name)
}

func (opt domainOption) ApplyHTTP(opts *httpOptions) {
	opts.Domain = string(opt)
}

func (opt domainOption) ApplyTLS(opts *tlsOptions) {
	opts.Domain = string(opt)
}

type hostnameOption string

// WithHostname sets the hostname for this edge.
//
// Deprecated: use WithDomain instead
func WithHostname(name string) interface {
	HTTPEndpointOption
	TLSEndpointOption
} {
	return hostnameOption(name)
}

func (opt hostnameOption) ApplyHTTP(opts *httpOptions) {
	opts.Hostname = string(opt)
}

func (opt hostnameOption) ApplyTLS(opts *tlsOptions) {
	opts.Hostname = string(opt)
}

type subdomainOption string

// WithSubdomain sets the subdomain for this edge.
//
// Deprecated: use WithDomain instead
func WithSubdomain(name string) interface {
	HTTPEndpointOption
	TLSEndpointOption
} {
	return subdomainOption(name)
}

func (opt subdomainOption) ApplyHTTP(opts *httpOptions) {
	opts.Subdomain = string(opt)
}

func (opt subdomainOption) ApplyTLS(opts *tlsOptions) {
	opts.Subdomain = string(opt)
}
