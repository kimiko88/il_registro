import SwiftUI

struct ParentPaymentsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Avvisi di Pagamento Attivi (PagoPA)")) {
                    HStack {
                        VStack(alignment: .leading, spacing: 4) {
                            Text("Contributo e Assicurazione")
                                .fontWeight(.bold)
                            Text("Scadenza: 30 Set 2026 • € 65,00")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button("Paga") {}
                            .buttonStyle(.borderedProminent)
                            .tint(Color.blue)
                    }
                }
            }
            .navigationTitle("Pagamenti Scolastici")
        }
    }
}
