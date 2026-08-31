import SwiftUI

struct SecretaryInventoryView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Cespiti Inventariati (D.I. 129/2018)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Smart Board 75\" 4K Touch")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Inventariato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("INV-2026-0491 • Aula 12")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("PNRR Piano Scuola 4.0 - Next Gen Labs")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Inventario & Cespiti")
        }
    }
}
