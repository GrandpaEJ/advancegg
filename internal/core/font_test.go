package core

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"os"
	"strings"
	"testing"
)

func TestLoadFontBytes_WOFFDetection(t *testing.T) {
	// Create a dummy file with WOFF header
	tmpFile, err := os.CreateTemp("", "testfont*.woff")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	header := "wOFF"
	if _, err := tmpFile.WriteString(header); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Load it
	_, err = LoadFontBytes(tmpFile.Name())
	// Should fail at ParseWOFF because it's incomplete, but it confirms detection path entered
	if err == nil {
		t.Error("Expected error for incomplete WOFF file")
	} else if err.Error() == "EOF" || err.Error() == "unexpected EOF" {
		// Expected error from binary.Read on truncated file
	} else {
		//t.Logf("Got expected error: %v", err)
	}
}

func TestParseWOFF_Structure(t *testing.T) {
	// Construct a minimal valid WOFF structure
	// Header + TableDirectory

	buf := new(bytes.Buffer)

	// WOFF Header
	// Signature: wOFF
	buf.WriteString("wOFF")
	// Flavor: 0x00010000 (TrueType 1.0)
	binary.Write(buf, binary.BigEndian, uint32(0x00010000))
	// Length: will patch later
	lengthOffset := buf.Len()
	binary.Write(buf, binary.BigEndian, uint32(0))
	// NumTables: 1
	binary.Write(buf, binary.BigEndian, uint16(1))
	// Reserved
	binary.Write(buf, binary.BigEndian, uint16(0))
	// TotalSfntSize: 0 (ignored for test?)
	binary.Write(buf, binary.BigEndian, uint32(100))
	// Major/Minor
	binary.Write(buf, binary.BigEndian, uint16(1))
	binary.Write(buf, binary.BigEndian, uint16(0))
	// Meta/Priv (0)
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))

	// Table Directory Entry
	// Tag: 'cmap'
	buf.WriteString("cmap")
	// Offset: will patch
	offsetOffset := buf.Len()
	binary.Write(buf, binary.BigEndian, uint32(0))
	// CompLength
	binary.Write(buf, binary.BigEndian, uint32(4))
	// OrigLength
	binary.Write(buf, binary.BigEndian, uint32(4))
	// Checksum
	binary.Write(buf, binary.BigEndian, uint32(0))

	// Data (uncompressed for simplicity, but WOFF expects compressed usually? code handles both if Comp < Orig)
	// Here Comp == Orig, so uncompressed.
	dataOffset := buf.Len()
	// Pad to 4 bytes? it is 4 bytes.
	buf.WriteString("DATA")

	// Patch Length
	totalLen := buf.Len()
	data := buf.Bytes()
	binary.BigEndian.PutUint32(data[lengthOffset:], uint32(totalLen))
	binary.BigEndian.PutUint32(data[offsetOffset:], uint32(dataOffset))

	// Test
	sfnt, err := ParseWOFF(data)
	if err != nil {
		t.Fatalf("ParseWOFF failed: %v", err)
	}

	// Verify SFNT structure
	// SFNT Header (12 bytes) + Table Dir (16 bytes) + Data (4 bytes)
	// Output should signature be 0x00010000
	if len(sfnt) < 12 {
		t.Fatal("SFNT too short")
	}
	flavor := binary.BigEndian.Uint32(sfnt[:4])
	if flavor != 0x00010000 {
		t.Errorf("Expected flavor 0x00010000, got 0x%x", flavor)
	}

	// Verify table tag found?
	// SFNT dir is after 12 bytes.
	// Tag is at 12.
	tag := string(sfnt[12:16])
	if tag != "cmap" {
		t.Errorf("Expected tag 'cmap', got '%s'", tag)
	}
}

func TestZlibCompression(t *testing.T) {
	// Test decompression path

	// Create compressible data
	originalData := []byte(strings.Repeat("A", 100))

	// Compress it
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(originalData)
	w.Close()
	compressed := b.Bytes()

	// Ensure it is actually compressed
	if len(compressed) >= len(originalData) {
		t.Skip("Compression didn't shrink data, skipping test")
	}

	buf := new(bytes.Buffer)
	// Header...
	buf.WriteString("wOFF")
	binary.Write(buf, binary.BigEndian, uint32(0x00010000))
	lengthOffset := buf.Len()
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint16(1)) // NumTables
	binary.Write(buf, binary.BigEndian, uint16(0))
	binary.Write(buf, binary.BigEndian, uint32(100))
	binary.Write(buf, binary.BigEndian, uint16(1))
	binary.Write(buf, binary.BigEndian, uint16(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(0))

	// Table Entry
	buf.WriteString("head")
	offsetOffset := buf.Len()
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, uint32(len(compressed)))   // CompLength
	binary.Write(buf, binary.BigEndian, uint32(len(originalData))) // OrigLength
	binary.Write(buf, binary.BigEndian, uint32(0))

	dataOffset := buf.Len()
	buf.Write(compressed)

	// Patch
	totalLen := buf.Len()
	data := buf.Bytes()
	binary.BigEndian.PutUint32(data[lengthOffset:], uint32(totalLen))
	binary.BigEndian.PutUint32(data[offsetOffset:], uint32(dataOffset))

	sfnt, err := ParseWOFF(data)
	if err != nil {
		t.Fatalf("ParseWOFF with compression failed: %v", err)
	}

	// Verify data decompressed
	if !bytes.Contains(sfnt, originalData) {
		t.Error("Decompressed data not found in output")
	}
}
