import SwiftUI

struct ParentTaxReceiptsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Certificazioni per Modello 730 / Unico")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Spese Scolastiche Anno 2025")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Detraibile 19%")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Totale Tracciato: € 680,00 PagoPA")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Mensa (€ 420), Gite (€ 245), Assicurazione (€ 15)")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        Button(action: {}) {
                            Label("Scarica Prospetto Fiscale PDF", systemImage: "arrow.down.doc.fill")
                        }
                        .buttonStyle(.bordered)
                        .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Attestazioni 730")
        }
    }
}
