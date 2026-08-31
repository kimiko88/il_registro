import SwiftUI

struct StudentTextbooksView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Libri Adottati - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Matematica.blu 2.0 (Vol. 3)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Da Acquistare")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("M. Bergamini • Zanichelli")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text("ISBN: 978-8808520852 • € 34,50")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)

                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("La Divina Commedia (Inferno)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Possesso")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Dante Alighieri • Mondadori Scuola")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text("ISBN: 978-8824731201 • € 21,00")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Libri di Testo")
        }
    }
}
