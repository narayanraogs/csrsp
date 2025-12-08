package das

import (
	"csrsp/server/utils/binary"
	"csrsp/server/utils/netutil"
	"log/slog"
	"strings"
)

const (
	monHeaderSize  = 12
	monTrailerSize = 4
	monSignature   = 1234567890
	monEndSig      = 0xB669FD2E
)

func getMonitorValue(conn *netutil.TimeoutConn, compCode int, equipmentType string) []byte {
	// Build Request
	requestArray := make([]byte, 0)
	requestArray = append(requestArray, binary.Uint32ToBytesBE(uint32(compCode))...)

	if strings.EqualFold(equipmentType, "hdr") {
		requestArray = append(requestArray, make([]byte, 12)...) // Append 3 zero uint32s
	}

	requestLength := len(requestArray) + monHeaderSize + monTrailerSize
	finalRequest := make([]byte, 0, requestLength)

	finalRequest = append(finalRequest, binary.Uint32ToBytesBE(monSignature)...)
	finalRequest = append(finalRequest, binary.Uint32ToBytesBE(uint32(requestLength))...)
	finalRequest = append(finalRequest, binary.Uint32ToBytesBE(0)...)
	finalRequest = append(finalRequest, requestArray...)
	finalRequest = append(finalRequest, binary.Uint32ToBytesBE(monEndSig)...)

	// Send Request
	if _, err := conn.Write(finalRequest); err != nil {
		slog.Error("Write timeout on monitor port", "ip", conn.RemoteAddr().String(), "error", err)
		return nil
	}

	// Read Header
	headerBytes := make([]byte, monHeaderSize)
	if _, err := conn.ReadFull(headerBytes); err != nil {
		slog.Error("Read timeout on monitor port (header)", "ip", conn.RemoteAddr().String(), "error", err)
		return nil
	}

	// Parse Length
	uLength, _ := binary.BytesToUint32BE(headerBytes[4:8])
	length := int(uLength)

	// Read Response
	if length <= monHeaderSize {
		return headerBytes // Should not happen for valid responses with data
	}

	responseBytes := make([]byte, length-monHeaderSize)
	if _, err := conn.ReadFull(responseBytes); err != nil {
		slog.Error("Read timeout on monitor port (body)", "ip", conn.RemoteAddr().String(), "error", err)
		return nil
	}

	return append(headerBytes, responseBytes...)
}
