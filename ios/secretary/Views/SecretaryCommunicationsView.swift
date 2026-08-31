import SwiftUI

struct SecretaryCommunicationsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Circolari & Comunicazioni Inviate")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Circolare N. 215: Chiusura Elezioni")
                                .fontWeight(.bold)
                            Spacer()
                            Text("89% Letti")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Destinatari: Famiglie & Docenti • Invio: Oggi")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Richiesta Firma: Sì • Scadenza: 2 Giugno")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Circolari & Avvisi")
        }
    }
}
