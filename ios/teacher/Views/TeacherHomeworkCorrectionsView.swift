import SwiftUI

struct TeacherHomeworkCorrectionsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Elaborati Consegnati da Valutare")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Relazione Laboratorio Fisica")
                                .fontWeight(.bold)
                            Spacer()
                            Text("22 / 24")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("Studente: Mario Rossi (Relazione_Rossi.pdf)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Consegnato puntuale • In attesa di voto")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        Button(action: {}) {
                            Label("Apri Elaborato e Valuta", systemImage: "pencil.and.outline")
                        }
                        .buttonStyle(.borderedProminent)
                        .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Correzione Compiti")
        }
    }
}
