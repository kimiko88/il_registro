import SwiftUI

struct SecretaryCertificatesProtocolView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Certificati Emessi con Protocollo")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Certificato Iscrizione e Frequenza")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Sigillo QR")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Studente: Mario Rossi (3ª A) • Uso INPS")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Protocollo: PROT-2026-004819 • Rilasciato oggi")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Certificati & Protocollo")
        }
    }
}
