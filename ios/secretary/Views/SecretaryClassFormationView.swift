import SwiftUI

struct SecretaryClassFormationView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Formazione Sezioni Classi Prime")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Sezione 1ª A (Scientifico)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Ottimale 98%")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("24 Alunni (12 M / 12 F) • Media Uscita: 8.2")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Equilibrio provenienza: 4 Istituti • Studenti PDP: 2")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Formazione Classi")
        }
    }
}
