package scanner

// EncryptedFile is the transfer-object shape from the scanner's bridge module.
type EncryptedFile struct {
	OriginalName string
	PayloadSize  int
	LibraryUsed  string
	Algorithm    string
	CipherText   string
}

// ProcessPostQuantumEncryption returns the static sample records shipped with that module.
// It does not transform or protect a file.
func ProcessPostQuantumEncryption(targetDir string) []EncryptedFile {
	_ = targetDir
	return []EncryptedFile{
		{
			OriginalName: "mbank_swifts_transfers.dat (Межбанковские переводы)",
			PayloadSize:  142050,
			LibraryUsed:  "NIST Post-Quantum Compliance Core",
			Algorithm:    "ML-KEM-768 (Квантово-стойкая решетчатая криптография)",
			CipherText:   "pq-nist:1e30S2iu+37IDV5NHEA1887h3be9PTV2/PRCCCzLpAM=...",
		},
		{
			OriginalName: "china_clearance_payments.xml (Трансграничные расчеты)",
			PayloadSize:  98440,
			LibraryUsed:  "GmSSL Суверенная китайская библиотека",
			Algorithm:    "SM4-GCM (Суверенный стандарт КНР, стойкий к квантовым атакам)",
			CipherText:   "gmssl-sm4:pkxF4TZX8DHL6P140POddQ==...",
		},
	}
}
