package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	//"strconv"
	//"encoding/hex"
	//"fmt"
	"math/big"
	"../lib"
	"./ripemd160"

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

/*
func PrvToString(prv ecdsa.PrivateKey) string {
	src := []byte(fmt.Sprintf("%d", prv.D))
	return hex.DecodeString(src)
}
*/

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
	return EncodeBig([]byte{}, bigJoin(30, r, s)), nil
}


func (w Wallet) SignatureVerify(sig, hash []byte) bool {
	pub := w.PrivateKey.PublicKey
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




/*



func (w Wallet) SignatureVerify2(sig, hash []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(sig)
	r.SetBytes(sig[:(sigLen / 2)])
	s.SetBytes(sig[(sigLen / 2):])

	fmt.Println("-------------")
	fmt.Println(*r)
	fmt.Println(*s)
		
	rawPubKey := w.PrivateKey.PublicKey
	if ecdsa.Verify(&rawPubKey, hash, &r, &s) == false {
		return false
	}
	return true

}



func (w Wallet) SignatureVerify(sig, hash []byte) bool {
	b, _ := DecodeToBig(w.publicKey)
	publ := splitBig(b, 2)
	x, y := publ[0], publ[1]
	b, _ = DecodeToBig(sig)
	sigg := splitBig(b, 2)
	r, s := sigg[0], sigg[1]
	pub := ecdsa.PublicKey{elliptic.P256(), x, y}
	return ecdsa.Verify(&pub, hash, r, s)
}

func (w Wallet) Sign(dataToSign string) {
	d, err := DecodeToBig(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	b, _ := DecodeToBig(w.PublicKey)
	pub := splitBig(b, 2)
	x, y := pub[0], pub[1]
	key := ecdsa.PrivateKey{ecdsa.PublicKey{elliptic.P256(), x, y}, d}
	r, s, _ := ecdsa.Sign(rand.Reader, &key, hash)
	return EncodeBig([]byte{}, bigJoin(KEY_SIZE, r, s)), nil
	
	
	curve := elliptic.P256()
	r, s, _ := ecdsa.Sign(rand.Reader, &key, hash)
	return EncodeBig([]byte{}, bigJoin(KEY_SIZE, r, s)), nil
	return ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
}
*/

