import SwiftUI

struct SecretaryExamCommitteeView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Commissioni Esame di Stato (Maturità)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Commissione RMIS08102")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Insediata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Presidente: Prof.ssa Maria Ricci")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("3 Interni + 3 Esterni • 48 Candidati • Commissione Web")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Commissioni Esami")
        }
    }
}
