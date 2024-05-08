package tlshelper

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type TlsOptions struct {
	Tls          string
	Cert         string
	Key          string
	Cacert       string
	AllowedUsers string
}

func LoadKeyPairSingle(options TlsOptions) credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(options.Cert, options.Key)
	if err != nil {
		panic(err)
	}
	return creds
}

func LoadKeyPairMutual(options TlsOptions) credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair(options.Cert, options.Key)
	if err != nil {
		panic(err)
	}

	ca, err := ioutil.ReadFile(options.Cacert)
	if err != nil {
		panic(err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		panic("AppendCertsFromPEM failed")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
	}

	return credentials.NewTLS(tlsConfig)
}
