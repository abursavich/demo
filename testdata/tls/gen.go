// SPDX-License-Identifier: MIT
//
// Copyright 2020 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"bursavich.dev/demo/internal/tlstest"
)

func main() {
	now := time.Now()
	end := now.Add(time.Hour * 24 * 365 * 2)

	rootCA, rootCACertPEMBlock, _, err := tlstest.GenerateCert(&tlstest.CertOptions{
		Template: &x509.Certificate{
			NotBefore: now,
			NotAfter:  end,
		},
	})
	check("Failed to create CA", err)
	write("ca.crt", rootCACertPEMBlock)

	_, certPEMBlock, keyPEMBlock, err := tlstest.GenerateCert(&tlstest.CertOptions{
		Parent: rootCA,
		Template: &x509.Certificate{
			DNSNames:  []string{"localhost"},
			NotBefore: now,
			NotAfter:  end,
		},
	})
	check("Failed to create cert", err)
	write("tls.crt", certPEMBlock)
	write("tls.key", keyPEMBlock)
}

func write(path string, buf []byte) {
	check(
		fmt.Sprintf("Failed to write file: %s", path),
		ioutil.WriteFile(path, buf, 0644),
	)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%v: %v\n", msg, err)
		os.Exit(1)
	}
}
