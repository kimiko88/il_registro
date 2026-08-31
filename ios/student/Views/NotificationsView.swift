import SwiftUI

struct NotificationsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Notifiche Recenti")) {
                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Nuovo Voto Inserito: 8.5")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Oggi 11:30")
                                .font(.caption2)
                                .foregroundColor(.secondary)
                        }
                        Text("Matematica (Scritto) • Prof.ssa Bianchi")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)

                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Circolare N. 214")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Ieri")
                                .font(.caption2)
                                .foregroundColor(.secondary)
                        }
                        Text("Calendario Scrutini e Adempimenti Finali")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Centro Notifiche")
        }
    }
}
