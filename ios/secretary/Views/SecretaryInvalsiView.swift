import SwiftUI

struct SecretaryInvalsiView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Prove Nazionali INVALSI CBT")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Matematica - Grado 13 (5ª)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Completata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("72/72 Studenti • Lab Informatica 1 & 2")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Somministratori: Bianchi, Neri • Trasmesso ad INVALSI")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Prove INVALSI")
        }
    }
}
