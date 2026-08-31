import SwiftUI

struct StudentCurriculumView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Dossier Esame di Stato")) {
                    HStack {
                        VStack(alignment: .leading, spacing: 4) {
                            Text("Curriculum dello Studente Ufficiale")
                                .fontWeight(.bold)
                            Text("Modello conforme alle Linee Guida MIM")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button(action: {}) {
                            Label("PDF", systemImage: "arrow.down.doc.fill")
                        }
                        .buttonStyle(.borderedProminent)
                    }
                }

                Section(header: Text("Riepilogo delle Sezioni")) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("1. Percorso degli Studi")
                            .font(.headline)
                        Text("Liceo Scientifico Statale (5 Anni) • Crediti: 35/40")
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        Divider()

                        Text("2. Certificazioni Riconosciute")
                            .font(.headline)
                        Text("Cambridge English B2 First • ICDL Full Standard")
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        Divider()

                        Text("3. Percorsi PCTO & Orientamento")
                            .font(.headline)
                        Text("120 ore PCTO validate • 30 ore Orientamento svolte")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Curriculum Studente")
        }
    }
}
