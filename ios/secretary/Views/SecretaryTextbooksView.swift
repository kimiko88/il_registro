import SwiftUI

struct SecretaryTextbooksView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Tetti di Spesa Libri di Testo")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Classe 3ª A • Scientifico")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Entro il Tetto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Spesa: € 281,50 / € 310,00 Max")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("9 Testi adottati • Delibera N. 18")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Adozione Libri")
        }
    }
}
