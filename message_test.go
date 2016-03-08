package netlog

import (
	"bytes"
	"hash/crc32"
	"math/rand"
	"testing"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	data := randData(rand.Intn(990) + 10)
	msg := MessageFromPayload(data)
	testMessage(t, data, msg)
}

func testMessage(t *testing.T, data []byte, msg Message) {
	dlen := len(data)

	if len(msg.Bytes()) != dlen+headerSize {
		t.Fatalf("Bad Message. Invalid message length from payload: %d vs expected %d", len(msg.Bytes()), dlen+headerSize)
	}

	if int(msg.PLength()) != dlen {
		t.Fatalf("Bad Message. Invalid payload length: %d vs expected %d", msg.PLength(), dlen)
	}

	crc := crc32.ChecksumIEEE(data)
	if crc != msg.CRC32() {
		t.Fatalf("Bad Message. Invalid CRC32: %d vs expected %d", crc, msg.CRC32())
	}

	if !msg.ChecksumOK() {
		t.Error("Bad Message. Self checksum failed.")
	}

	if !bytes.Equal(data, msg.Payload()) {
		t.Errorf("Bad Message. Payload not equal to original data.\n Got: % x\n Exp: % x\n", msg.Payload(), data)
	}
}

func TestUnpackSequence(t *testing.T) {
	t.Parallel()

	// random number of payloads with random bytes
	data := make([][]byte, rand.Intn(90)+10)
	for k := range data {
		data[k] = randData(rand.Intn(90) + 10)
	}

	messages := make([]Message, len(data))
	for k := range data {
		messages[k] = MessageFromPayload(data[k])
	}

	var sequence []byte
	for _, m := range messages {
		sequence = append(sequence, m.Bytes()...)
	}

	unpacked, err := Unpack(sequence)
	if err != nil {
		t.Error(err)
	}

	if len(unpacked) != len(messages) {
		t.Errorf("Unpacked %d messages vs expected %d", len(unpacked), len(messages))
	}

	for k, d := range data {
		testMessage(t, d, unpacked[k])
	}
}