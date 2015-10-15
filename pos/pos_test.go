package pos

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var prover *Prover = nil
var verifier *Verifier = nil
var pk []byte
var size int = 4
var beta int = 2
var graphDir string = "graph"

//exp* gets setup in test.go

func TestPoS(t *testing.T) {
	challenges := verifier.SelectChallenges()
	hashes, proofs := prover.ProveSpace(challenges)
	if !verifier.VerifySpace(challenges, hashes, proofs) {
		log.Fatal("Verify space failed:", challenges)
	}
}

func TestOpenVerify(t *testing.T) {
	hash, proof := prover.Open(1)
	for i := range expProof {
		for j := range expProof[i] {
			if expProof[i][j] != proof[i][j] {
				log.Fatal("Open failed:", expProof[i], proof[i])
			}
		}
	}

	if !verifier.Verify(1,hash, proof) {
		log.Fatal("Verify failed:", hash, proof)
	}
}

//Sanity check using simple graph
//[0 0 0 0]
//[1 0 0 0]
//[1 0 0 0]
//[1 0 1 0]
func TestComputeHash(t *testing.T) {
	hashes := make([][]byte, size)
	for i := range hashes {
		f, _ := os.Open(fmt.Sprintf("%s/%d/hash", graphDir, i))
		hashes[i] = make([]byte, hashSize)
		f.Read(hashes[i])
	}

	var result [hashSize]byte

	for i := range expHashes {
		copy(result[:], hashes[i])
		if expHashes[i] != result {
			log.Fatal("Hash mismatch:", expHashes[i], result)
		}

	}
}

func TestMerkleTree(t *testing.T) {
	result := make([][hashSize]byte, 2*size)
	for i := 1; i < size; i++ {
		f, _ := os.Open(fmt.Sprintf("%s/merkle/%d", graphDir, i))
		buf := make([]byte, hashSize)
		f.Read(buf)
		copy(result[i][:], buf)
	}
	for i := 0; i < size; i++ {
		f, _ := os.Open(fmt.Sprintf("%s/%d/hash", graphDir, i))
		buf := make([]byte, hashSize)
		f.Read(buf)
		copy(result[i+size][:], buf)
	}

	for i := 2*size-1; i > 0; i-- {
		if expMerkle[i] != result[i] {
			log.Fatal("Merkle node mismatch:", i, expMerkle[i], result[i])
		}
	}

}

func TestMain(m *testing.M) {
	pk = []byte{1}
	prover = NewProver(pk, size, graphDir)
	setup(pk, size, graphDir)
	root := prover.InitGraph()

	verifier = NewVerifier(pk, size, beta, root)
	os.Exit(m.Run())
}
