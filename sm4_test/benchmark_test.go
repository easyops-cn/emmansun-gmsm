package sm4_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/emmansun/gmsm/sm4"
)

func benchmarkCBCEncrypt1K(b *testing.B, block cipher.Block) {
	buf := make([]byte, 1024)
	b.SetBytes(int64(len(buf)))

	var iv [16]byte
	cbc := cipher.NewCBCEncrypter(block, iv[:])
	for i := 0; i < b.N; i++ {
		cbc.CryptBlocks(buf, buf)
	}
}

func BenchmarkAESCBCEncrypt1K(b *testing.B) {
	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	benchmarkCBCEncrypt1K(b, c)
}

func BenchmarkSM4CBCEncrypt1K(b *testing.B) {
	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	benchmarkCBCEncrypt1K(b, c)
}

func benchmarkSM4CBCDecrypt1K(b *testing.B, block cipher.Block) {
	buf := make([]byte, 1024)
	b.SetBytes(int64(len(buf)))

	var iv [16]byte
	cbc := cipher.NewCBCDecrypter(block, iv[:])
	for i := 0; i < b.N; i++ {
		cbc.CryptBlocks(buf, buf)
	}
}

func BenchmarkAESCBCDecrypt1K(b *testing.B) {
	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	benchmarkSM4CBCDecrypt1K(b, c)
}

func BenchmarkSM4CBCDecrypt1K(b *testing.B) {
	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	benchmarkSM4CBCDecrypt1K(b, c)
}

func benchmarkStream(b *testing.B, block cipher.Block, mode func(cipher.Block, []byte) cipher.Stream, buf []byte) {
	b.SetBytes(int64(len(buf)))

	//var key [16]byte
	var iv [16]byte
	//c, _ := sm4.NewCipher(key[:])
	stream := mode(block, iv[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream.XORKeyStream(buf, buf)
	}
}

func benchmarkSM4Stream(b *testing.B, mode func(cipher.Block, []byte) cipher.Stream, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	benchmarkStream(b, c, mode, buf)
}

func benchmarkAESStream(b *testing.B, mode func(cipher.Block, []byte) cipher.Stream, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	benchmarkStream(b, c, mode, buf)
}

// If we test exactly 1K blocks, we would generate exact multiples of
// the cipher's block size, and the cipher stream fragments would
// always be wordsize aligned, whereas non-aligned is a more typical
// use-case.
const almost1K = 1024 - 5
const almost8K = 8*1024 - 5

func BenchmarkAESCFBEncrypt1K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewCFBEncrypter, make([]byte, almost1K))
}

func BenchmarkSM4CFBEncrypt1K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewCFBEncrypter, make([]byte, almost1K))
}

func BenchmarkAESCFBDecrypt1K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewCFBDecrypter, make([]byte, almost1K))
}

func BenchmarkSM4CFBDecrypt1K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewCFBDecrypter, make([]byte, almost1K))
}

func BenchmarkAESCFBDecrypt8K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewCFBDecrypter, make([]byte, almost8K))
}

func BenchmarkSM4CFBDecrypt8K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewCFBDecrypter, make([]byte, almost8K))
}

func BenchmarkAESOFB1K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewOFB, make([]byte, almost1K))
}

func BenchmarkSM4OFB1K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewOFB, make([]byte, almost1K))
}

func BenchmarkAESCTR1K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewCTR, make([]byte, almost1K))
}

func BenchmarkSM4CTR1K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewCTR, make([]byte, almost1K))
}

func BenchmarkAESCTR8K(b *testing.B) {
	benchmarkAESStream(b, cipher.NewCTR, make([]byte, almost8K))
}

func BenchmarkSM4CTR8K(b *testing.B) {
	benchmarkSM4Stream(b, cipher.NewCTR, make([]byte, almost8K))
}

func benchmarkGCMSign(b *testing.B, aead cipher.AEAD, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var nonce [12]byte
	var out []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out = aead.Seal(out[:0], nonce[:], nil, buf)
	}
}

func benchmarkAESGCMSign(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	aesgcm, _ := cipher.NewGCM(c)
	benchmarkGCMSign(b, aesgcm, buf)
}

func benchmarkSM4GCMSign(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	sm4gcm, _ := cipher.NewGCM(c)
	benchmarkGCMSign(b, sm4gcm, buf)
}

func benchmarkGCMSeal(b *testing.B, aead cipher.AEAD, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var nonce [12]byte
	var ad [13]byte
	var out []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out = aead.Seal(out[:0], nonce[:], buf, ad[:])
	}
}

func benchmarkAESGCMSeal(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	sm4gcm, _ := cipher.NewGCM(c)
	benchmarkGCMSeal(b, sm4gcm, buf)
}

func benchmarkSM4GCMSeal(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	sm4gcm, _ := cipher.NewGCM(c)
	benchmarkGCMSeal(b, sm4gcm, buf)
}

func benchmarkGCMOpen(b *testing.B, aead cipher.AEAD, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var nonce [12]byte
	var ad [13]byte
	var out []byte
	out = aead.Seal(out[:0], nonce[:], buf, ad[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := aead.Open(buf[:0], nonce[:], out, ad[:])
		if err != nil {
			b.Errorf("Open: %v", err)
		}
	}
}

func benchmarkAESGCMOpen(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := aes.NewCipher(key[:])
	sm4gcm, _ := cipher.NewGCM(c)
	benchmarkGCMOpen(b, sm4gcm, buf)
}

func benchmarkSM4GCMOpen(b *testing.B, buf []byte) {
	var key [16]byte
	c, _ := sm4.NewCipher(key[:])
	sm4gcm, _ := cipher.NewGCM(c)
	benchmarkGCMOpen(b, sm4gcm, buf)
}

func BenchmarkAESGCMSeal1K(b *testing.B) {
	benchmarkAESGCMSeal(b, make([]byte, 1024))
}

func BenchmarkSM4GCMSeal1K(b *testing.B) {
	benchmarkSM4GCMSeal(b, make([]byte, 1024))
}

func BenchmarkAESGCMOpen1K(b *testing.B) {
	benchmarkAESGCMOpen(b, make([]byte, 1024))
}

func BenchmarkSM4GCMOpen1K(b *testing.B) {
	benchmarkSM4GCMOpen(b, make([]byte, 1024))
}

func BenchmarkAESGCMSign1K(b *testing.B) {
	benchmarkAESGCMSign(b, make([]byte, 1024))
}

func BenchmarkSM4GCMSign1K(b *testing.B) {
	benchmarkSM4GCMSign(b, make([]byte, 1024))
}

func BenchmarkAESGCMSign8K(b *testing.B) {
	benchmarkAESGCMSign(b, make([]byte, 8*1024))
}

func BenchmarkSM4GCMSign8K(b *testing.B) {
	benchmarkSM4GCMSign(b, make([]byte, 8*1024))
}

func BenchmarkAESGCMSeal8K(b *testing.B) {
	benchmarkAESGCMSeal(b, make([]byte, 8*1024))
}

func BenchmarkSM4GCMSeal8K(b *testing.B) {
	benchmarkSM4GCMSeal(b, make([]byte, 8*1024))
}

func BenchmarkAESGCMOpen8K(b *testing.B) {
	benchmarkAESGCMOpen(b, make([]byte, 8*1024))
}

func BenchmarkSM4GCMOpen8K(b *testing.B) {
	benchmarkSM4GCMOpen(b, make([]byte, 8*1024))
}
