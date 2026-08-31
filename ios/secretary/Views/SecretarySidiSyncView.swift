import SwiftUI

struct SecretarySidiSyncView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Esportazione Flussi SIDI (MIM)")) {
                    HStack {
                        VStack(alignment: .leading, spacing: 2) {
                            Text("1. Anagrafe Nazionale Studenti (ANS)")
                                .fontWeight(.bold)
                            Text("Dati anagrafici e iscrizioni")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button("Esporta XML") {}
                            .buttonStyle(.bordered)
                    }

                    HStack {
                        VStack(alignment: .leading, spacing: 2) {
                            Text("2. Esiti Scrutini Differiti")
                                .fontWeight(.bold)
                            Text("Sospensione del giudizio settembre")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button("Esporta XML") {}
                            .buttonStyle(.bordered)
                    }
                }
            }
            .navigationTitle("Flussi SIDI")
        }
    }
}
