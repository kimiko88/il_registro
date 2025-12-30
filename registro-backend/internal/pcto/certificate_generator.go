package pcto

import "fmt"

type CertificateGenerator struct{}

func (g *CertificateGenerator) GenerateCertificate(project *Project, studentName string) ([]byte, error) {
	// Stub: Return dummy PDF bytes
	return []byte(fmt.Sprintf("Certificate for %s in %s", studentName, project.Title)), nil
}

type AgreementGenerator struct{}

func (g *AgreementGenerator) GenerateAgreement(project *Project, company *Company) ([]byte, error) {
	return []byte("Agreement Document Stub"), nil
}
