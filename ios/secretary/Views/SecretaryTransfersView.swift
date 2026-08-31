import SwiftUI

struct SecretaryTransfersView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Pratiche Nulla Osta & Trasferimenti")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Nulla Osta in Uscita: L. B.")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Rilasciato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Destinazione: Liceo 'Galilei' • Classe 2ª B")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Prot: NO-2026-0392 • SIDI Sincronizzato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Trasferimenti")
        }
    }
}
