import SwiftUI

struct SecretaryAlboPretorioView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Albo Pretorio Online (D.Lgs. 33/2013)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Determina Dirigenziale N. 104")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Affisso (15gg)")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("ALBO-2026-00342 • Fornitura Lab STEM PNRR")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Affisso: 25 Maggio • Scadenza: 09 Giugno • SHA-256")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Albo Pretorio")
        }
    }
}
