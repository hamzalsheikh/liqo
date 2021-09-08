// Copyright 2019-2021 The Liqo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iptables

import (
	"testing"

	. "github.com/coreos/go-iptables/iptables"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIptables(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iptables Suite")
}

var _ = BeforeSuite(func() {
	var err error
	h, err = NewIPTHandler()
	Expect(err).To(BeNil())
	ipt, err = New()
	Expect(err).To(BeNil())
	err = h.Init()
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	err := h.Terminate()
	Expect(err).To(BeNil())
})
