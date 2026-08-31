import SwiftUI

struct ParentCyberbullyingReportView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Segnalazioni Anti-Bullismo (L. 71/17)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Pratica BULL-2026-081")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Gestione")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Episodio Chat di Gruppo • 20 Maggio 2026")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Referente: Prof.ssa Verdi • Supporto psicologico attivato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Anti-Bullismo")
        }
    }
}
