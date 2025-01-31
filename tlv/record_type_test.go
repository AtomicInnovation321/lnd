package tlv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fakeCsvDelayType = 1

	fakeIsCoolType = 2
)

type fakeWireMsg struct {
	CsvDelay RecordT[TlvType1, uint16]

	IsCool RecordT[TlvType2, bool]
}

// TestRecordTFromPrimitive tests the RecordT type. We should be able to create
// types of both record types, and also primitive types, and encode/decode them
// as normal.
func TestRecordTFromPrimitive(t *testing.T) {
	t.Parallel()

	wireMsg := fakeWireMsg{
		CsvDelay: NewPrimitiveRecord[TlvType1](uint16(5)),
		IsCool:   NewPrimitiveRecord[TlvType2](true),
	}

	encodeStream, err := NewStream(
		wireMsg.CsvDelay.Record(), wireMsg.IsCool.Record(),
	)
	require.NoError(t, err)

	var b bytes.Buffer
	err = encodeStream.Encode(&b)
	require.NoError(t, err)

	var newWireMsg fakeWireMsg

	decodeStream, err := NewStream(
		newWireMsg.CsvDelay.Record(),
		newWireMsg.IsCool.Record(),
	)
	require.NoError(t, err)

	err = decodeStream.Decode(&b)
	require.NoError(t, err)

	require.Equal(t, wireMsg, newWireMsg)
}

type wireCsv uint16

func (w *wireCsv) Record() Record {
	return MakeStaticRecord(fakeCsvDelayType, (*uint16)(w), 2, EUint16, DUint16)
}

type coolWireMsg struct {
	CsvDelay RecordT[TlvType1, wireCsv]
}

// TestRecordTFromRecord tests that we can create a RecordT type from an
// existing record type and encode/decode as normal.
func TestRecordTFromRecord(t *testing.T) {
	t.Parallel()

	val := wireCsv(5)

	wireMsg := coolWireMsg{
		CsvDelay: NewRecordT[TlvType1](val),
	}

	encodeStream, err := NewStream(wireMsg.CsvDelay.Record())
	require.NoError(t, err)

	var b bytes.Buffer
	err = encodeStream.Encode(&b)
	require.NoError(t, err)

	var wireMsg2 coolWireMsg

	decodeStream, err := NewStream(wireMsg2.CsvDelay.Record())
	require.NoError(t, err)

	err = decodeStream.Decode(&b)
	require.NoError(t, err)

	require.Equal(t, wireMsg, wireMsg2)
}
