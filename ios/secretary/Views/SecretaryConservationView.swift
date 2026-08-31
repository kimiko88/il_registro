import SwiftUI

struct SecretaryConservationView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Versamenti a Norma AgID / ParER")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Registri Docenti - 1° Quadrimestre")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Conservato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Hash SHA-256 Verificato • Conservatore Regionale")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Rapporto di Versamento N. 8192 archiviato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Conservazione a Norma")
        }
    }
}
