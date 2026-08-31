import SwiftUI

struct SecretaryMaterialsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Stato Storage Cloud Didattica")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Storage Utilizzato")
                                .fontWeight(.bold)
                            Spacer()
                            Text("42.8 GB / 250 GB")
                                .font(.caption2)
                                .foregroundColor(.blue)
                        }
                        ProgressView(value: 0.17)
                            .padding(.vertical, 2)
                        Text("1.482 File • Scansione Antivirus: OK")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Materiali Didattici")
        }
    }
}
