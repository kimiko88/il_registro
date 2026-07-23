package pcto

import "fmt"

type CertificateGenerator struct{}

func (g *CertificateGenerator) GenerateCertificate(project *Project, studentName string) ([]byte, error) {
	html := fmt.Sprintf(`
		<html>
		<head><style>body{font-family: Arial, sans-serif; text-align: center; margin-top: 50px;}</style></head>
		<body>
			<h1>Certificato PCTO</h1>
			<h2>Attestato di Partecipazione</h2>
			<p>Si certifica che lo studente <strong>%s</strong></p>
			<p>ha partecipato con successo al progetto:</p>
			<h3>%s</h3>
			<p>Periodo: %s - %s</p>
		</body>
		</html>`, 
		studentName, project.Title, 
		project.StartDate.Format("02/01/2006"), project.EndDate.Format("02/01/2006"),
	)
	return []byte(html), nil
}

type AgreementGenerator struct{}

func (g *AgreementGenerator) GenerateAgreement(project *Project, company *Company) ([]byte, error) {
	html := fmt.Sprintf(`
		<html>
		<head><style>body{font-family: Arial, sans-serif; margin: 40px;}</style></head>
		<body>
			<h2>Convenzione PCTO</h2>
			<p>Tra l'istituto scolastico e l'azienda ospitante <strong>%s</strong> (P.IVA %s)</p>
			<p>Si stipula la convenzione per il progetto:</p>
			<h3>%s</h3>
			<p>Con le seguenti caratteristiche:</p>
			<ul>
				<li>Inizio: %s</li>
				<li>Fine: %s</li>
			</ul>
			<p>Firme:</p>
			<br><br>
			<p>___________________ (Scuola)</p>
			<br><br>
			<p>___________________ (Azienda)</p>
		</body>
		</html>`,
		company.Name, company.VatNumber, project.Title,
		project.StartDate.Format("02/01/2006"), project.EndDate.Format("02/01/2006"),
	)
	return []byte(html), nil
}
