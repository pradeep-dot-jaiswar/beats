// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tlscommon

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/snappyflow/beats/v7/libbeat/logp"
)

const logSelector = "tls"

// LoadCertificate will load a certificate from disk and return a tls.Certificate or error
func LoadCertificate(config *CertificateConfig) (*tls.Certificate, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	certificate := config.Certificate
	key := config.Key
	if certificate == "" {
		return nil, nil
	}

	log := logp.NewLogger(logSelector)

	certPEM, err := ReadPEMFile(log, certificate, config.Passphrase)
	if err != nil {
		log.Errorf("Failed reading certificate file %v: %+v", certificate, err)
		return nil, fmt.Errorf("%v %v", err, certificate)
	}

	keyPEM, err := ReadPEMFile(log, key, config.Passphrase)
	if err != nil {
		log.Errorf("Failed reading key file %v: %+v", key, err)
		return nil, fmt.Errorf("%v %v", err, key)
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Errorf("Failed loading client certificate %+v", err)
		return nil, err
	}

	log.Debugf("tls", "loading certificate: %v and key %v", certificate, key)
	return &cert, nil
}

// ReadPEMFile reads a PEM format file on disk and decrypt it with the privided password and
// return the raw content.
func ReadPEMFile(log *logp.Logger, path, passphrase string) ([]byte, error) {
	pass := []byte(passphrase)
	var blocks []*pem.Block

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for len(content) > 0 {
		var block *pem.Block

		block, content = pem.Decode(content)
		if block == nil {
			if len(blocks) == 0 {
				return nil, errors.New("no pem file")
			}
			break
		}

		if x509.IsEncryptedPEMBlock(block) {
			var buffer []byte
			var err error
			if len(pass) == 0 {
				err = errors.New("No passphrase available")
			} else {
				// Note, decrypting pem might succeed even with wrong password, but
				// only noise will be stored in buffer in this case.
				buffer, err = x509.DecryptPEMBlock(block, pass)
			}

			if err != nil {
				log.Errorf("Dropping encrypted pem '%v' block read from %v. %+v",
					block.Type, path, err)
				continue
			}

			// DEK-Info contains encryption info. Remove header to mark block as
			// unencrypted.
			delete(block.Headers, "DEK-Info")
			block.Bytes = buffer
		}
		blocks = append(blocks, block)
	}

	if len(blocks) == 0 {
		return nil, errors.New("no PEM blocks")
	}

	// re-encode available, decrypted blocks
	buffer := bytes.NewBuffer(nil)
	for _, block := range blocks {
		err := pem.Encode(buffer, block)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// LoadCertificateAuthorities read the slice of CAcert and return a Certpool.
func LoadCertificateAuthorities(CAs []string) (*x509.CertPool, []error) {
	errors := []error{}

	if len(CAs) == 0 {
		return nil, nil
	}

	log := logp.NewLogger(logSelector)
	roots := x509.NewCertPool()
	for _, path := range CAs {
		pemData, err := ioutil.ReadFile(path)
		if err != nil {
			log.Errorf("Failed reading CA certificate: %+v", err)
			errors = append(errors, fmt.Errorf("%v reading %v", err, path))
			continue
		}

		if ok := roots.AppendCertsFromPEM(pemData); !ok {
			log.Error("Failed to add CA to the cert pool, CA is not a valid PEM file")
			errors = append(errors, fmt.Errorf("%v adding %v to the list of known CAs", ErrNotACertificate, path))
			continue
		}
		log.Debugf("tls", "successfully loaded CA certificate: %v", path)
	}

	return roots, errors
}

func extractMinMaxVersion(versions []TLSVersion) (uint16, uint16) {
	if len(versions) == 0 {
		versions = TLSDefaultVersions
	}

	minVersion := uint16(0xffff)
	maxVersion := uint16(0)
	for _, version := range versions {
		v := uint16(version)
		if v < minVersion {
			minVersion = v
		}
		if v > maxVersion {
			maxVersion = v
		}
	}

	return minVersion, maxVersion
}

// ResolveTLSVersion takes the integer representation and return the name.
func ResolveTLSVersion(v uint16) string {
	return TLSVersion(v).String()
}

// ResolveCipherSuite takes the integer representation and return the cipher name.
func ResolveCipherSuite(cipher uint16) string {
	return tlsCipherSuite(cipher).String()
}
