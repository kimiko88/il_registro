import SwiftUI

struct SecretaryInvoicingView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Flussi di Incasso PagoPA")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Flusso PagoPA #2026-05-28")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Quadratura 100%")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Incassato: € 14.850,00 • 165 IUV")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Viaggi d'istruzione, Assicurazione, Contributo")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Riconciliazione PagoPA")
        }
    }
}
