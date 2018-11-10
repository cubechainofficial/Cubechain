package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	"github.com/cubechainofficial/Cubechain/lib"
	"github.com/cubechainofficial/Cubechain/wallet/ripemd160"
)

const prefix = byte(0x1C)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func CreateWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	prefixPayload := append([]byte{prefix}, pubKeyHash...)
	checksum := checksum(prefixPayload)
	fullPayload := append(prefixPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	prefix := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{prefix}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func (w Wallet) Sign(hash []byte) ([]byte, error) {
	key := w.PrivateKey
	r, s, _ := ecdsa.Sign(rand.Reader, &key, hash)
	// TODO: where is EncodeBig?
	return EncodeBig([]byte{}, bigJoin(30, r, s)), nil
}

func (w Wallet) SignatureVerify(sig, hash []byte) bool {
	pub := w.PrivateKey.PublicKey
	// TODO: where is DecodeBig?
	b, _ := DecodeToBig(sig)
	sigg := splitBig(b, 2)
	r, s := sigg[0], sigg[1]
	return ecdsa.Verify(&pub, hash, r, s)
}

func bigJoin(expectedLen int, bigs ...*big.Int) *big.Int {
	bs := []byte{}
	for i, b := range bigs {
		by := b.Bytes()
		dif := expectedLen - len(by)
		if dif > 0 && i != 0 {
			by = append(lib.ArrayOfBytes(dif, 0), by...)
		}
		bs = append(bs, by...)
	}
	b := new(big.Int).SetBytes(bs)
	return b
}

func splitBig(b *big.Int, parts int) []*big.Int {
	bs := b.Bytes()
	if len(bs)%2 != 0 {
		bs = append([]byte{0}, bs...)
	}
	l := len(bs) / parts
	as := make([]*big.Int, parts)
	for i, _ := range as {
		as[i] = new(big.Int).SetBytes(bs[i*l : (i+1)*l])
	}
	return as
}
