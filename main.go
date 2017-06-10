// Copyright 2017 The HopfieldPress Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pointlander/gopfield/hopfield"
)

const (
	Size = 32
)

func main() {
	net, err := hopfield.NewNetwork(Size, "hebbian")
	if err != nil {
		panic(err)
	}

	fmt.Println(net.Capacity())

	data, err := ioutil.ReadFile("alice30.txt")
	if err != nil {
		panic(err)
	}
	output := make([]byte, len(data))
	copy(output, data)

	var fbuffer [Size]float64
	var ffbuffer [Size]float64
	for i := range fbuffer {
		fbuffer[i] = -1
	}
	shift := func() {
		for i, f := 0, -1.0; i < Size; i++ {
			fbuffer[i], f = f, fbuffer[i]
		}
	}
	for x, b := range data {
		fmt.Printf("%c", b)
		for i := 0; i < 8; i++ {
			shift()
			copy(ffbuffer[:], fbuffer[:])
			pattern := hopfield.Encode(ffbuffer[:])
			res, err := net.Restore(pattern, "async", 2)
			if err != nil {
				panic(err)
			}
			out := 0
			if res.At(0) > 0 {
				out = 1
			}
			output[x] = output[x] ^ byte(out<<uint(i))

			if b&1 == 0 {
				fbuffer[0] = -1
			} else {
				fbuffer[0] = 1
			}
			b >>= 1
			pattern = hopfield.Encode(fbuffer[:])
			for i := 0; i < 1; i++ {
				err = net.Store([]*hopfield.Pattern{pattern})
				if err != nil {
					panic(err)
				}
			}
		}
	}

	write := func(name string, data []byte) {
		file, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := gzip.NewWriter(file)
		defer writer.Close()
		_, err = writer.Write(data)
		if err != nil {
			panic(err)
		}
	}
	write("alice30.hop", output)
	write("alice30.gz", data)
}
