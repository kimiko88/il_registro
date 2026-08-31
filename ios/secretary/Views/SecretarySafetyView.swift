import SwiftUI

struct SecretarySafetyView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Prove di Evacuazione & Sicurezza")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Prova 2° Semestre (D.Lgs. 81/08)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Riuscita")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Tempo: 2 min 45 sec • 588 Persone evacuate")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Verbale RSPP N. 04 archiviato a norma")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Sicurezza & DVR")
        }
    }
}
