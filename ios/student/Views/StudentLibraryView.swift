import SwiftUI

struct StudentLibraryView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Prestiti Biblioteca Attivi")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Il Fu Mattia Pascal")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Prestito")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("L. Pirandello • Collocazione: NARR-PIR-04")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Scadenza: 10 Giugno • Inventario: BIB-8821")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Biblioteca & Libri")
        }
    }
}
