import SwiftUI

struct StudentExamSimulationsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Simulazioni Esame di Stato (MIM)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("2ª Prova: Matematica")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Traccia Ufficiale")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("2 Problemi + 8 Quesiti (Scienze Applicate)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Durata: 6h • Griglia max 20 Punti")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        Button(action: {}) {
                            Label("Scarica Traccia & Griglia PDF", systemImage: "arrow.down.doc.fill")
                        }
                        .buttonStyle(.bordered)
                        .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Simulazioni Esami")
        }
    }
}
