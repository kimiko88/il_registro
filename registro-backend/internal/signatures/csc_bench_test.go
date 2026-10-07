package signatures

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkGenerateDigitalStamp(b *testing.B) {
	key := []byte("benchmark-key-seal-cad")
	sha := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateDigitalStamp("https://scuola.edu.it", "doc-bench", "pagella", sha, "Prof. Rossi", "Dirigente", key)
	}
}

func BenchmarkVerifyDigitalStamp(b *testing.B) {
	key := []byte("benchmark-key-seal-cad")
	sha := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
	stamp, _ := GenerateDigitalStamp("https://scuola.edu.it", "doc-bench", "pagella", sha, "Prof. Rossi", "Dirigente", key)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyDigitalStamp(stamp.QRCodePayload, sha, key)
	}
}

func BenchmarkCSCBatchSign(b *testing.B) {
	connector := NewCSCConnector(&SoftwareQSCD{}, nil)
	ctx := context.Background()

	const count = 10
	docs := make([]string, count)
	for i := 0; i < count; i++ {
		docs[i] = fmt.Sprintf("pagella-%d", i)
	}

	req := CSCBatchSignRequest{
		DocumentIDs: docs,
		PIN:         "1234",
		OTP:         "654321",
		SignerRole:  "Dirigente Scolastico",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = connector.BatchSign(ctx, "user-ds", "school-1", req)
	}
}
