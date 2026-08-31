import SwiftUI

struct SecretaryStaffContractsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Contratti di Supplenza Stipulati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Supplenza: Prof.ssa Neri")
                                .fontWeight(.bold)
                            Spacer()
                            Text("NoiPA Sincronizzato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("A026 (Matematica) • 18 Ore Settimanali")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("18/05/2026 - 06/06/2026 • SIDI: CT-2026-0812")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Contratti & Supplenze")
        }
    }
}
