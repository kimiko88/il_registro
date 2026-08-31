import SwiftUI

struct StudentCertificationsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Certificazioni Validate")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Cambridge English B2 First")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Validata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Livello QCER: B2 • Cambridge Assessment")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Rilascio: Marzo 2026 • Credito assegnato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Certificazioni")
        }
    }
}
