// Code generated by github.com/alanctgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCES:
 *     block.avsc
 *     header.avsc
 */

package avro

import (
	"io"
)

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func readAvroContainerBlock(r io.Reader) (*AvroContainerBlock, error) {
	var str = &AvroContainerBlock{}
	var err error
	str.NumRecords, err = readLong(r)
	if err != nil {
		return nil, err
	}
	str.RecordBytes, err = readBytes(r)
	if err != nil {
		return nil, err
	}
	str.Sync, err = readSync(r)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func readAvroContainerHeader(r io.Reader) (*AvroContainerHeader, error) {
	var str = &AvroContainerHeader{}
	var err error
	str.Magic, err = readMagic(r)
	if err != nil {
		return nil, err
	}
	str.Meta, err = readMapBytes(r)
	if err != nil {
		return nil, err
	}
	str.Sync, err = readSync(r)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func readBytes(r io.Reader) ([]byte, error) {
	size, err := readLong(r)
	if err != nil {
		return nil, err
	}
	bb := make([]byte, size)
	_, err = io.ReadFull(r, bb)
	return bb, err
}

func readLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
}

func readMagic(r io.Reader) (Magic, error) {
	var bb Magic
	_, err := io.ReadFull(r, bb[:])
	return bb, err
}

func readMapBytes(r io.Reader) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		if blkSize == 0 {
			break
		}
		if blkSize < 0 {
			blkSize = -blkSize
			_, err := readLong(r)
			if err != nil {
				return nil, err
			}
		}
		for i := int64(0); i < blkSize; i++ {
			key, err := readString(r)
			if err != nil {
				return nil, err
			}
			val, err := readBytes(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readString(r io.Reader) (string, error) {
	len, err := readLong(r)
	if err != nil {
		return "", err
	}
	bb := make([]byte, len)
	_, err = io.ReadFull(r, bb)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func readSync(r io.Reader) (Sync, error) {
	var bb Sync
	_, err := io.ReadFull(r, bb[:])
	return bb, err
}

func writeAvroContainerBlock(r *AvroContainerBlock, w io.Writer) error {
	var err error
	err = writeLong(r.NumRecords, w)
	if err != nil {
		return err
	}
	err = writeBytes(r.RecordBytes, w)
	if err != nil {
		return err
	}
	err = writeSync(r.Sync, w)
	if err != nil {
		return err
	}

	return nil
}
func writeAvroContainerHeader(r *AvroContainerHeader, w io.Writer) error {
	var err error
	err = writeMagic(r.Magic, w)
	if err != nil {
		return err
	}
	err = writeMapBytes(r.Meta, w)
	if err != nil {
		return err
	}
	err = writeSync(r.Sync, w)
	if err != nil {
		return err
	}

	return nil
}

func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeMagic(r Magic, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}

func writeMapBytes(r map[string][]byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeBytes(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

func writeSync(r Sync, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}
