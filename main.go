// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package main

import (
	"fmt"

	messageConverter "github.com/LoRaWanSoFa/Core/messageConverter"
)

func main() {
	test := messageConverter.New()

	bytes := []byte{0x3F, 0xB0, 0xFC, 0x00, 0xAA, 0xA0, 0x84, 0x41}
	fmt.Println(test.ConvertSingleValue(bytes, 1))
}
