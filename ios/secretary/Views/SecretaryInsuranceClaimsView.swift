import SwiftUI

struct SecretaryInsuranceClaimsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Denunce Infortuni & Sinistri")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Infortunio Palestra (Distorsione)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("INAIL Inviato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Studente 2ª B • 12 Maggio (10:15)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Protocollo INAIL: 9021 • Sinistro Assicurazione: 4412")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Infortuni & Assicurazione")
        }
    }
}
