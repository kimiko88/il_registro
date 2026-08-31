import SwiftUI

struct ParentAnagraficaView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Dati Anagrafici & Contatti")) {
                    HStack {
                        Text("Tutore Principale")
                        Spacer()
                        Text("Giuseppe Rossi (Padre)")
                            .fontWeight(.bold)
                    }

                    HStack {
                        Text("Telefono Emergenza")
                        Spacer()
                        Text("+39 333 1234567")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Email Istituzionale")
                        Spacer()
                        Text("famiglia.rossi@email.it")
                            .foregroundColor(.secondary)
                    }

                    HStack {
                        Text("Stato Iscrizione")
                        Spacer()
                        Text("Confermata A.S. 25/26")
                            .fontWeight(.bold)
                            .foregroundColor(.green)
                    }
                }
            }
            .navigationTitle("Dati Famiglia")
        }
    }
}
