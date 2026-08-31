import SwiftUI

struct ParentDelegationsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Delegati al Ritiro Autorizzati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Sig.ra Anna Bianchi (Nonna)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Verificata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("CI N. CA92819XX • Scadenza: 2028")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Accolta dalla Segreteria per l'intero A.S.")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Deleghe & Ritiri")
        }
    }
}
