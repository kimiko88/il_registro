import SwiftUI

struct SecretarySubstitutionsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Coperture da Assegnare")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Classe 4ª B • 3ª Ora (10:00 - 11:00)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Docente Assente")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.red.opacity(0.2))
                                .foregroundColor(.red)
                                .cornerRadius(4)
                        }
                        Text("Docente Proposto: Prof.ssa Bianchi (Ora a disposizione)")
                            .font(.subheadline)
                            .foregroundColor(.blue)

                        Button("Conferma Assegnazione") {}
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Gestione Supplenze")
        }
    }
}
