import SwiftUI

struct SecretaryElectionsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Consultazioni Elettorali Online")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Rinnovo Consiglio d'Istituto")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Seggio Aperto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Affluenza: 68.4% (842/1230 Votanti)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Chiusura ore 18:00 • Scrutinio crittografato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Elezioni Collegiali")
        }
    }
}
